// Example: authenticate, create an invoice, poll for payment.
//
//	QPAY_USERNAME=... QPAY_PASSWORD=... QPAY_INVOICE_CODE=... \
//	  go run ./examples/basic
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	fmt.Println("authenticated with QPay sandbox")

	inv, err := client.CreateInvoice(ctx, qpay.CreateInvoiceRequest{
		InvoiceCode:         os.Getenv("QPAY_INVOICE_CODE"),
		SenderInvoiceNo:     fmt.Sprintf("EX-%d", time.Now().Unix()),
		InvoiceReceiverCode: "terminal",
		InvoiceDescription:  "qpay-sdk example invoice",
		Amount:              100,
		CallbackURL:         "https://example.com/webhooks/qpay",
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
