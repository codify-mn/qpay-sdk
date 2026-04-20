package qpay

import (
	"encoding/json"
	"errors"
	"net/http"
)

// ParseWebhook decodes the JSON body of a QPay callback POST request.
// The request body is closed by the caller (http.Server's responsibility).
func ParseWebhook(r *http.Request) (*WebhookPayload, error) {
	if r == nil || r.Body == nil {
		return nil, errors.New("qpay: nil request")
	}
	var p WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}
