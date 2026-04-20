package qpay

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCheckPayment_IsPaid(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token":
			tokenHandler(w)
		case "/v2/payment/check":
			var req checkPaymentRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			if req.ObjectType != "INVOICE" || req.ObjectID != "INV" {
				t.Fatalf("bad body: %+v", req)
			}
			_, _ = w.Write([]byte(`{"count":1,"paid_amount":500,"rows":[{"payment_id":"P1","payment_status":"PAID","payment_amount":500}]}`))
		}
	})
	defer srv.Close()

	pc, err := c.CheckPayment(context.Background(), "INV")
	if err != nil {
		t.Fatal(err)
	}
	paid, p := pc.IsPaid()
	if !paid || p.PaymentID != "P1" {
		t.Fatalf("expected paid row, got paid=%v p=%+v", paid, p)
	}
}

func TestIsPaid_nil(t *testing.T) {
	var pc *PaymentCheck
	paid, p := pc.IsPaid()
	if paid || p != nil {
		t.Fatalf("nil PaymentCheck must report unpaid")
	}
}

func TestCancelAndRefund(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token":
			tokenHandler(w)
		case "/v2/payment/cancel/P1", "/v2/payment/refund/P1":
			if r.Method != http.MethodDelete {
				t.Fatalf("want DELETE, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
		}
	})
	defer srv.Close()

	if err := c.CancelPayment(context.Background(), "P1", RefundRequest{Note: "test"}); err != nil {
		t.Fatal(err)
	}
	if err := c.RefundPayment(context.Background(), "P1", RefundRequest{Note: "test"}); err != nil {
		t.Fatal(err)
	}
}

func TestListPayments(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token":
			tokenHandler(w)
		case "/v2/payment/list":
			_, _ = w.Write([]byte(`{"count":2,"paid_amount":1000,"rows":[{"payment_id":"A"},{"payment_id":"B"}]}`))
		}
	})
	defer srv.Close()

	list, err := c.ListPayments(context.Background(), PaymentListRequest{MerchantID: "M1"})
	if err != nil {
		t.Fatal(err)
	}
	if list.Count != 2 || len(list.Rows) != 2 {
		t.Fatalf("bad list: %+v", list)
	}
}
