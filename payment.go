package qpay

import (
	"context"
	"net/http"
)

// CheckPayment returns all payments recorded against an invoice.
func (c *Client) CheckPayment(ctx context.Context, invoiceID string) (*PaymentCheck, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/v2/payment/check", checkPaymentRequest{
		ObjectType: "INVOICE",
		ObjectID:   invoiceID,
	})
	if err != nil {
		return nil, err
	}
	var pc PaymentCheck
	if err := decodeJSON(body, &pc); err != nil {
		return nil, err
	}
	return &pc, nil
}

// GetPayment fetches a single payment by ID.
func (c *Client) GetPayment(ctx context.Context, paymentID string) (*Payment, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/v2/payment/"+paymentID, nil)
	if err != nil {
		return nil, err
	}
	var p Payment
	if err := decodeJSON(body, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// CancelPayment cancels a payment with an optional note.
func (c *Client) CancelPayment(ctx context.Context, paymentID string, req RefundRequest) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/v2/payment/cancel/"+paymentID, req)
	return err
}

// RefundPayment refunds a payment with an optional note.
func (c *Client) RefundPayment(ctx context.Context, paymentID string, req RefundRequest) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/v2/payment/refund/"+paymentID, req)
	return err
}

// ListPayments searches payments by merchant, status, and date range.
func (c *Client) ListPayments(ctx context.Context, req PaymentListRequest) (*PaymentList, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/v2/payment/list", req)
	if err != nil {
		return nil, err
	}
	var pl PaymentList
	if err := decodeJSON(body, &pl); err != nil {
		return nil, err
	}
	return &pl, nil
}
