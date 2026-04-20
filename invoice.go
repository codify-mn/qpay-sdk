package qpay

import (
	"context"
	"net/http"
)

// CreateInvoice creates a new invoice and returns QR + bank deep links.
func (c *Client) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*Invoice, error) {
	if req.TerminalID == "" && c.terminalID != "" {
		req.TerminalID = c.terminalID
	}
	body, err := c.doRequest(ctx, http.MethodPost, "/v2/invoice", req)
	if err != nil {
		return nil, err
	}
	var inv Invoice
	if err := decodeJSON(body, &inv); err != nil {
		return nil, err
	}
	return &inv, nil
}

// GetInvoice fetches details for a known invoice ID.
func (c *Client) GetInvoice(ctx context.Context, invoiceID string) (*InvoiceDetail, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/v2/invoice/"+invoiceID, nil)
	if err != nil {
		return nil, err
	}
	var inv InvoiceDetail
	if err := decodeJSON(body, &inv); err != nil {
		return nil, err
	}
	return &inv, nil
}

// CancelInvoice cancels an open invoice.
func (c *Client) CancelInvoice(ctx context.Context, invoiceID string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/v2/invoice/"+invoiceID, nil)
	return err
}
