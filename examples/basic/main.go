// Example: authenticate, create an invoice, poll for payment.
//
//	QPAY_USERNAME=... QPAY_PASSWORD=... QPAY_MERCHANT_ID=... \
//	  go run ./examples/basic
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	qpay "github.com/codify-mn/qpay-sdk"
)

func main() {
	client, err := qpay.New(
		qpay.WithSandbox(),
		qpay.WithCredentials(os.Getenv("QPAY_USERNAME"), os.Getenv("QPAY_PASSWORD")),
	)
	if err != nil {
		log.Fatalf("client: %v", err)
	}

	ctx := context.Background()

	if err := client.Ping(ctx); err != nil {
		log.Fatalf("ping: %v", err)
	}
	fmt.Println("authenticated with QPay")

	inv, err := client.CreateInvoice(ctx, qpay.CreateInvoiceRequest{
		MerchantID:   os.Getenv("QPAY_MERCHANT_ID"),
		Amount:       100,
		Currency:     "MNT",
		CustomerName: "Example Co",
		CallbackURL:  "https://example.com/webhooks/qpay",
		Description:  "qpay-sdk example invoice",
	})
	if err != nil {
		log.Fatalf("create invoice: %v", err)
	}
	fmt.Printf("invoice %s\n  QR: %s\n  short URL: %s\n", inv.InvoiceID, inv.QRImage, inv.ShortURL)

	pc, err := client.CheckPayment(ctx, inv.InvoiceID)
	if err != nil {
		log.Fatalf("check: %v", err)
	}
	if paid, p := pc.IsPaid(); paid {
		fmt.Printf("paid: %s amount=%v\n", p.PaymentID, p.PaymentAmount)
	} else {
		fmt.Println("not yet paid — scan the QR with a Mongolian bank app")
	}
}
