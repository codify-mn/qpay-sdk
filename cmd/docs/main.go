// Command docs serves the qpay-go OpenAPI spec rendered with Scalar.
//
//	go run ./cmd/docs              # listens on :8080
//	go run ./cmd/docs -addr :9000  # custom port
//
// Then open the printed URL in a browser.
package main

import (
	"flag"
	"io/fs"
	"log"
	"net/http"

	qpay "github.com/codify-mn/qpay-go"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	sub, err := fs.Sub(qpay.DocsFS, "docs")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(sub)))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })

	log.Printf("qpay-go docs → http://localhost%s/scalar.html", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
