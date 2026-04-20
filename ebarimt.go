package qpay

import (
	"context"
	"net/http"
)

// CreateEbarimt creates an eBarimt receipt for a completed payment.
func (c *Client) CreateEbarimt(ctx context.Context, req CreateEbarimtRequest) (*Ebarimt, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/v2/ebarimt/create", req)
	if err != nil {
		return nil, err
	}
	var e Ebarimt
	if err := decodeJSON(body, &e); err != nil {
		return nil, err
	}
	return &e, nil
}
