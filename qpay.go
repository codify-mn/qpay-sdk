// Package qpay is a Go SDK for the QPay v2 Mongolian payment API.
//
// Create a client with functional options and call domain methods directly:
//
//	client, _ := qpay.New(
//	    qpay.WithSandbox(),
//	    qpay.WithCredentials("user", "pass"),
//	)
//	inv, err := client.CreateInvoice(ctx, qpay.CreateInvoiceRequest{...})
package qpay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// Client is the main entry point for QPay API calls.
type Client struct {
	httpClient *http.Client
	log        *slog.Logger

	baseURL    string
	username   string
	password   string
	terminalID string

	mu           sync.RWMutex
	accessToken  string
	refreshToken string
	tokenExpiry  time.Time
}

// New constructs a Client. WithCredentials is required.
func New(opts ...Option) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    productionBaseURL,
		log:        slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.username == "" || c.password == "" {
		return nil, ErrMissingCredentials
	}
	return c, nil
}

// NewFromConfig is a shortcut for consumers who have a Config struct loaded from YAML.
func NewFromConfig(cfg Config, opts ...Option) (*Client, error) {
	all := []Option{
		WithCredentials(cfg.Username, cfg.Password),
		WithTerminalID(cfg.TerminalID),
	}
	if cfg.BaseURL != "" {
		all = append(all, WithBaseURL(cfg.BaseURL))
	}
	all = append(all, opts...)
	return New(all...)
}

// Ping verifies credentials by requesting a token. Use as a startup health check.
func (c *Client) Ping(ctx context.Context) error {
	return c.fetchToken(ctx)
}

// doRequest performs an authenticated request, refreshing the token once on 401.
func (c *Client) doRequest(ctx context.Context, method, path string, body any) ([]byte, error) {
	if err := c.ensureValidToken(ctx); err != nil {
		return nil, err
	}
	return c.do(ctx, method, path, body, true)
}

func (c *Client) do(ctx context.Context, method, path string, body any, allowRetry bool) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("qpay: marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("qpay: build request: %w", err)
	}

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("qpay: execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("qpay: read response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized && allowRetry {
		c.log.Debug("qpay: 401, refreshing token", "path", path)
		if err := c.refresh(ctx); err != nil {
			return nil, err
		}
		return c.do(ctx, method, path, body, false)
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, respBody)
	}

	return respBody, nil
}

// decodeJSON unmarshals a JSON response into dst.
func decodeJSON(body []byte, dst any) error {
	if err := json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("qpay: decode response: %w", err)
	}
	return nil
}
