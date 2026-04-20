package qpay

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// fetchToken requests a fresh access/refresh token using Basic auth.
func (c *Client) fetchToken(ctx context.Context) error {
	var body io.Reader
	if c.terminalID != "" {
		b, _ := json.Marshal(map[string]string{"terminal_id": c.terminalID})
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v2/auth/token", body)
	if err != nil {
		return fmt.Errorf("qpay: build token request: %w", err)
	}

	basic := base64.StdEncoding.EncodeToString([]byte(c.username + ":" + c.password))
	req.Header.Set("Authorization", "Basic "+basic)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("qpay: execute token request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("qpay: read token response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrInvalidCredentials
	}
	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, respBody)
	}

	var tr TokenResponse
	if err := json.Unmarshal(respBody, &tr); err != nil {
		return fmt.Errorf("qpay: parse token response: %w", err)
	}

	c.storeToken(tr)
	return nil
}

// refresh swaps the refresh token for a new access token. Falls back to fetchToken on any failure.
func (c *Client) refresh(ctx context.Context) error {
	c.mu.RLock()
	refresh := c.refreshToken
	c.mu.RUnlock()

	if refresh == "" {
		return c.fetchToken(ctx)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v2/auth/refresh", nil)
	if err != nil {
		return fmt.Errorf("qpay: build refresh request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+refresh)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("qpay: execute refresh request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("qpay: read refresh response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return c.fetchToken(ctx)
	}

	var tr TokenResponse
	if err := json.Unmarshal(respBody, &tr); err != nil {
		return fmt.Errorf("qpay: parse refresh response: %w", err)
	}

	c.storeToken(tr)
	return nil
}

// ensureValidToken fetches a new token if none exists or the cached one is near expiry.
func (c *Client) ensureValidToken(ctx context.Context) error {
	c.mu.RLock()
	token := c.accessToken
	expiry := c.tokenExpiry
	c.mu.RUnlock()

	if token == "" || time.Now().Add(60*time.Second).After(expiry) {
		return c.refresh(ctx)
	}
	return nil
}

func (c *Client) storeToken(tr TokenResponse) {
	c.mu.Lock()
	c.accessToken = tr.AccessToken
	c.refreshToken = tr.RefreshToken
	c.tokenExpiry = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	c.mu.Unlock()
}
