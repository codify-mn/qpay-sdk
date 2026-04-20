package qpay

import "embed"

// DocsFS exposes the OpenAPI spec and Scalar HTML bundled with the SDK.
// Mount it as a static file server to host the API reference locally.
//
//	mux.Handle("/", http.FileServer(http.FS(qpay.DocsFS)))
//
// Files: docs/openapi.yaml, docs/scalar.html.
//
//go:embed docs
var DocsFS embed.FS
