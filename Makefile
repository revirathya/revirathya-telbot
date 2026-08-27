run:
	go run ./cmd/api

build:
	go build ./...

vet:
	go vet ./...

## --- codegen ---

proto-generate:
	buf generate

proto-lint:
	buf lint