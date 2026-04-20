package qpay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestNew_requiresCredentials(t *testing.T) {
	if _, err := New(WithSandbox()); err == nil {
		t.Fatal("expected error without credentials")
	}
}

func TestPing_successCachesToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/auth/token" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Basic ") {
			t.Fatalf("missing basic auth header: %q", r.Header.Get("Authorization"))
		}
		_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "A", RefreshToken: "R", ExpiresIn: 3600})
	}))
	defer srv.Close()

	c, err := New(WithBaseURL(srv.URL), WithCredentials("u", "p"))
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.accessToken != "A" || c.refreshToken != "R" {
		t.Fatalf("tokens not stored: access=%q refresh=%q", c.accessToken, c.refreshToken)
	}
}

func TestPing_invalidCredentials(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer srv.Close()

	c, _ := New(WithBaseURL(srv.URL), WithCredentials("bad", "bad"))
	if err := c.Ping(context.Background()); err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestDoRequest_retriesOn401(t *testing.T) {
	var tokenCalls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token", "/v2/auth/refresh":
			atomic.AddInt32(&tokenCalls, 1)
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "A", RefreshToken: "R", ExpiresIn: 3600})
		case "/v2/invoice/test":
			// First call returns 401 (stale token). Second must succeed.
			if r.Header.Get("Authorization") == "Bearer A" && atomic.LoadInt32(&tokenCalls) == 1 {
				http.Error(w, "expired", http.StatusUnauthorized)
				return
			}
			_, _ = w.Write([]byte(`{"invoice_id":"test","invoice_code":"I","amount":1,"invoice_status":"OPEN","invoice_description":"","created_date":""}`))
		}
	}))
	defer srv.Close()

	c, _ := New(WithBaseURL(srv.URL), WithCredentials("u", "p"))
	if err := c.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
	atomic.StoreInt32(&tokenCalls, 1) // reset counter after Ping
	_, err := c.GetInvoice(context.Background(), "test")
	if err != nil {
		t.Fatalf("expected retry to succeed, got %v", err)
	}
}
