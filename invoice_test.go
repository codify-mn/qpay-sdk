package qpay

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, h http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(h)
	c, err := New(WithBaseURL(srv.URL), WithCredentials("u", "p"))
	if err != nil {
		t.Fatal(err)
	}
	return c, srv
}

func tokenHandler(w http.ResponseWriter) {
	_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "A", RefreshToken: "R", ExpiresIn: 3600})
}

func TestCreateInvoice(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token":
			tokenHandler(w)
		case "/v2/invoice":
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}
			var req CreateInvoiceRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatal(err)
			}
			if req.Amount != 12345 {
				t.Fatalf("amount not forwarded: got %v", req.Amount)
			}
			_, _ = w.Write([]byte(`{"id":"INV1","qr_code":"qr","qr_image":"img","urls":[{"name":"Khan","link":"khan://"}]}`))
		}
	})
	defer srv.Close()

	inv, err := c.CreateInvoice(context.Background(), CreateInvoiceRequest{
		MerchantID: "M1", Amount: 12345, Description: "Test", Currency: "MNT",
	})
	if err != nil {
		t.Fatal(err)
	}
	if inv.InvoiceID != "INV1" || len(inv.URLs) != 1 {
		t.Fatalf("bad decode: %+v", inv)
	}
}

func TestGetInvoice_apiError(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token":
			tokenHandler(w)
		case "/v2/invoice/missing":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"code":"INVOICE_NOT_FOUND","message":"not found"}`))
		}
	})
	defer srv.Close()

	_, err := c.GetInvoice(context.Background(), "missing")
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %v", err)
	}
	if apiErr.StatusCode != 404 || apiErr.Code != "INVOICE_NOT_FOUND" {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}
}

func TestCancelInvoice(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token":
			tokenHandler(w)
		case "/v2/invoice/abc":
			if r.Method != http.MethodDelete {
				t.Fatalf("expected DELETE")
			}
			w.WriteHeader(http.StatusOK)
		}
	})
	defer srv.Close()

	if err := c.CancelInvoice(context.Background(), "abc"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateInvoice_appliesClientTerminalID(t *testing.T) {
	var captured CreateInvoiceRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token":
			tokenHandler(w)
		case "/v2/invoice":
			_ = json.NewDecoder(r.Body).Decode(&captured)
			_, _ = w.Write([]byte(`{"id":"X"}`))
		}
	}))
	defer srv.Close()

	c, _ := New(WithBaseURL(srv.URL), WithCredentials("u", "p"), WithTerminalID("TERM123"))
	if _, err := c.CreateInvoice(context.Background(), CreateInvoiceRequest{Amount: 1}); err != nil {
		t.Fatal(err)
	}
	if captured.TerminalID != "TERM123" {
		t.Fatalf("terminal_id not applied: got %q", captured.TerminalID)
	}
}
