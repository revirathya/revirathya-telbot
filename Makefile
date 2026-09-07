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

## --- docker ---
docker-build:
	docker build -t revirathya/revirathya-telbot:latest .

docker-run:
	docker run --rm -p 8080:8080 --env-file .env revirathya/revirathya-telbot:latest
