package qpay

import (
	"context"
	"net/http"
	"testing"
)

func TestCreateEbarimt(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token":
			tokenHandler(w)
		case "/v2/ebarimt/create":
			if r.Method != http.MethodPost {
				t.Fatalf("want POST")
			}
			_, _ = w.Write([]byte(`{"id":"E1","payment_id":"P1","ebarimt_type":"CITIZEN","ebarimt_receipt":"..."}`))
		}
	})
	defer srv.Close()

	e, err := c.CreateEbarimt(context.Background(), CreateEbarimtRequest{PaymentID: "P1", EbarimtReceiver: "CITIZEN"})
	if err != nil {
		t.Fatal(err)
	}
	if e.EbarimtID != "E1" || e.PaymentID != "P1" {
		t.Fatalf("bad decode: %+v", e)
	}
}
