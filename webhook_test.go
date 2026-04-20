package qpay

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseWebhook(t *testing.T) {
	body := `{"type":"payment","object_id":"INV1","payment_id":"P1","status":"PAID"}`
	r := httptest.NewRequest(http.MethodPost, "/webhooks/qpay", strings.NewReader(body))
	p, err := ParseWebhook(r)
	if err != nil {
		t.Fatal(err)
	}
	if p.ObjectID != "INV1" || p.PaymentID != "P1" || p.Status != "PAID" {
		t.Fatalf("bad decode: %+v", p)
	}
}

func TestParseWebhook_nilRequest(t *testing.T) {
	if _, err := ParseWebhook(nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseWebhook_badJSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/webhooks/qpay", strings.NewReader("{not json"))
	if _, err := ParseWebhook(r); err == nil {
		t.Fatal("expected error")
	}
}
