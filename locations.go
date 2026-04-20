package qpay

import (
	"context"
	"net/http"
)

// GetCities lists QPay cities (aimags).
func (c *Client) GetCities(ctx context.Context) ([]Location, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/v2/aimaghot", nil)
	if err != nil {
		return nil, err
	}
	var out []Location
	if err := decodeJSON(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDistricts lists districts (sums) for a given city code.
func (c *Client) GetDistricts(ctx context.Context, cityCode string) ([]Location, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/v2/sumduureg/"+cityCode, nil)
	if err != nil {
		return nil, err
	}
	var out []Location
	if err := decodeJSON(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetBanks lists supported banks for merchant onboarding.
func (c *Client) GetBanks(ctx context.Context) ([]Bank, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/v2/bank/list", nil)
	if err != nil {
		return nil, err
	}
	var out []Bank
	if err := decodeJSON(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}
