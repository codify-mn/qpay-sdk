package qpay

import (
	"context"
	"net/http"
)


// CreateCompanyMerchant registers a company as a QPay merchant.
func (c *Client) CreateCompanyMerchant(ctx context.Context, req CreateCompanyMerchantRequest) (*Merchant, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/v2/merchant/company", req)
	if err != nil {
		return nil, err
	}
	var m Merchant
	if err := decodeJSON(body, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// CreatePersonMerchant registers an individual as a QPay merchant.
func (c *Client) CreatePersonMerchant(ctx context.Context, req CreatePersonMerchantRequest) (*Merchant, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/v2/merchant/person", req)
	if err != nil {
		return nil, err
	}
	var m Merchant
	if err := decodeJSON(body, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMerchant fetches a merchant by ID.
func (c *Client) GetMerchant(ctx context.Context, merchantID string) (*Merchant, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/v2/merchant/"+merchantID, nil)
	if err != nil {
		return nil, err
	}
	var m Merchant
	if err := decodeJSON(body, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// ListMerchants returns a paginated slice of merchants.
func (c *Client) ListMerchants(ctx context.Context, opts ListOptions) (*MerchantList, error) {
	if opts.Limit <= 0 {
		opts.Limit = 25
	}
	body, err := c.doRequest(ctx, http.MethodPost, "/v2/merchant/list", opts)
	if err != nil {
		return nil, err
	}
	var list MerchantList
	if err := decodeJSON(body, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// UpdateMerchant patches a merchant record.
func (c *Client) UpdateMerchant(ctx context.Context, merchantID string, req UpdateMerchantRequest) (*Merchant, error) {
	body, err := c.doRequest(ctx, http.MethodPut, "/v2/merchant/"+merchantID, req)
	if err != nil {
		return nil, err
	}
	var m Merchant
	if err := decodeJSON(body, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// DeleteMerchant removes a merchant by ID.
func (c *Client) DeleteMerchant(ctx context.Context, merchantID string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/v2/merchant/"+merchantID, nil)
	return err
}
