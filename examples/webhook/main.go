// Example: minimal HTTP handler consuming QPay webhooks.
//
//	go run ./examples/webhook
//	curl -X POST -H 'Content-Type: application/json' \
//	  -d '{"type":"payment","object_id":"INV1","payment_id":"P1","status":"PAID"}' \
//	  http://localhost:8081/webhooks/qpay
package main

import (
	"fmt"
	"log"
	"net/http"

	qpay "github.com/codify-mn/qpay-sdk"
)

func main() {
	http.HandleFunc("/webhooks/qpay", func(w http.ResponseWriter, r *http.Request) {
		payload, err := qpay.ParseWebhook(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("webhook: invoice=%s payment=%s status=%s", payload.ObjectID, payload.PaymentID, payload.Status)
		fmt.Fprintln(w, "ok")
	})

	log.Println("listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
