# qpay-go

Go SDK for [QPay](https://qpay.mn) v2 — Mongolian payment gateway.

## Install

```bash
go get github.com/codify-mn/qpay-go
```

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"

	qpay "github.com/codify-mn/qpay-go"
)

func main() {
	client, err := qpay.New(
		qpay.WithSandbox(),
		qpay.WithCredentials("YOUR_USERNAME", "YOUR_PASSWORD"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	inv, err := client.CreateInvoice(ctx, qpay.CreateInvoiceRequest{
		MerchantID:   "YOUR_MERCHANT_ID",
		InvoiceCode:  "INV-001",
		Description:  "Order #1",
		Amount:       10000,
		Currency:     "MNT",
		CallbackURL:  "https://yourapp.com/webhook",
		CustomerName: "Customer",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Invoice ID: %s\nQR: %s\n", inv.InvoiceID, inv.QRImage)
}
```

## Features

- Full QPay v2 API coverage: auth, merchant, invoice, payment, ebarimt, locations.
- Automatic token refresh (access + refresh tokens, mutex-protected).
- Retries 401 transparently by refreshing token.
- Typed `APIError` for all non-2xx responses.
- Context-aware (`context.Context` on every method).
- Zero external dependencies — stdlib only.
- Optional `*slog.Logger` integration.

## Documentation

Browse the full API spec locally:

```bash
make docs   # opens http://localhost:8080/docs/scalar.html
```

Powered by [Scalar](https://scalar.com) over the `docs/openapi.yaml` spec bundled in the module.

## Configuration

| Option | Purpose |
|---|---|
| `WithSandbox()` | Use `https://merchant-sandbox.qpay.mn`. |
| `WithProduction()` | Use `https://merchant.qpay.mn`. |
| `WithBaseURL(url)` | Override base URL (useful for tests). |
| `WithCredentials(user, pass)` | **Required.** Basic-auth creds for `/v2/auth/token`. |
| `WithTerminalID(id)` | Default terminal for invoice creation. |
| `WithHTTPClient(*http.Client)` | Inject custom transport. |
| `WithLogger(*slog.Logger)` | Optional structured logging. |
| `WithTimeout(d)` | HTTP client timeout (default 30s). |

## Webhooks

```go
http.HandleFunc("/webhooks/qpay", func(w http.ResponseWriter, r *http.Request) {
	payload, err := qpay.ParseWebhook(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Invoice %s status=%s", payload.ObjectID, payload.Status)
	w.WriteHeader(http.StatusOK)
})
```

## Testing

```bash
make test   # go test -race -cover ./...
```

All tests use `httptest` — no live QPay calls.

## License

MIT
