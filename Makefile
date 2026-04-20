.PHONY: test docs lint tidy vet build

test:
	go test -race -cover ./...

vet:
	go vet ./...

tidy:
	go mod tidy

build:
	go build ./...

docs:
	go run ./cmd/docs

lint: vet
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || echo "golangci-lint not installed, skipping"
