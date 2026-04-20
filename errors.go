package qpay

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrMissingCredentials = errors.New("qpay: credentials required (use WithCredentials)")
	ErrInvalidCredentials = errors.New("qpay: invalid credentials")
	ErrTokenExpired       = errors.New("qpay: token expired")
)

// APIError represents a non-2xx response from the QPay API.
type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
	Raw        []byte `json:"-"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("qpay: %d %s: %s", e.StatusCode, e.Code, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("qpay: %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("qpay: %d: %s", e.StatusCode, string(e.Raw))
}

func parseAPIError(status int, body []byte) *APIError {
	e := &APIError{StatusCode: status, Raw: body}
	_ = json.Unmarshal(body, e)
	return e
}
