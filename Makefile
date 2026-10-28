.PHONY: run build test test-integration vet check tidy lint \
	docker-build docker-run

run:
	go run -ldflags "-X main.version=$(VERSION)" ./cmd

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/financecore ./cmd

test:
	go test -count=1 ./test/...

test-integration:
	go test -race -count=1 -tags=integration ./test/...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: vet lint test proto-lint

tidy:
	go mod tidy

http:
	ngrok http 8080 --url https://semicommercial-quantummechanical-eva.ngrok-free.dev


## --- docker ---
docker-build:
	docker build -t revirathya/revirathya-telbot:latest .

docker-run:
	docker run --rm -p 8080:8080 --env-file .env revirathya/revirathya-telbot:latest

