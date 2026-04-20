package qpay

import (
	"log/slog"
	"net/http"
	"time"
)

const (
	sandboxBaseURL    = "https://merchant-sandbox.qpay.mn"
	productionBaseURL = "https://merchant.qpay.mn"
	defaultTimeout    = 30 * time.Second
)

// Option configures a Client.
type Option func(*Client)

func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

func WithSandbox() Option {
	return func(c *Client) { c.baseURL = sandboxBaseURL }
}

func WithProduction() Option {
	return func(c *Client) { c.baseURL = productionBaseURL }
}

func WithCredentials(username, password string) Option {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

func WithTerminalID(id string) Option {
	return func(c *Client) { c.terminalID = id }
}

func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.httpClient = h
		}
	}
}

func WithLogger(l *slog.Logger) Option {
	return func(c *Client) {
		if l != nil {
			c.log = l
		}
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.httpClient.Timeout = d
		}
	}
}
