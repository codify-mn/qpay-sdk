package qpay

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestCreateCompanyMerchant(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token":
			tokenHandler(w)
		case "/v2/merchant/company":
			_, _ = w.Write([]byte(`{"id":"M1","type":"COMPANY","name":"Acme"}`))
		}
	})
	defer srv.Close()

	m, err := c.CreateCompanyMerchant(context.Background(), CreateCompanyMerchantRequest{Name: "Acme"})
	if err != nil {
		t.Fatal(err)
	}
	if m.ID != "M1" || m.Type != "COMPANY" {
		t.Fatalf("bad decode: %+v", m)
	}
}

func TestListMerchants_paginationBody(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v2/auth/token":
			tokenHandler(w)
		case strings.HasPrefix(r.URL.Path, "/v2/merchant/list"):
			var body ListOptions
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body.PageNumber != 2 || body.PageLimit != 50 {
				t.Fatalf("pagination not forwarded: %+v", body)
			}
			_, _ = w.Write([]byte(`{"count":0,"rows":[]}`))
		}
	})
	defer srv.Close()

	if _, err := c.ListMerchants(context.Background(), ListOptions{PageNumber: 2, PageLimit: 50}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteMerchant(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token":
			tokenHandler(w)
		case "/v2/merchant/M1":
			if r.Method != http.MethodDelete {
				t.Fatalf("want DELETE")
			}
			w.WriteHeader(http.StatusOK)
		}
	})
	defer srv.Close()

	if err := c.DeleteMerchant(context.Background(), "M1"); err != nil {
		t.Fatal(err)
	}
}
