-include .env
export

TRIVY_IMAGE := aquasec/trivy@sha256:be1190afcb28352bfddc4ddeb71470835d16462af68d310f9f4bca710961a41e # 0.70.0

.PHONY: run build test test-integration vet leaks check tidy lint \
	docker-build docker-run

run:
	go run -ldflags "-X main.version=$(VERSION)" ./cmd

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/telbot ./cmd

test:
	go test -count=1 ./test/...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

leaks:
	gitleaks git -v --log-opts="--all"

check: vet lint leaks

tidy:
	go mod tidy

http:
	ngrok http 8080 --url $(NGROK_HOST)

hooks:
	git config core.hooksPath .githooks

## --- local infra ---
docker-build:
	docker build \
		-f infra/docker/Dockerfile \
		--build-arg VERSION=$(VERSION) \
		-t revirathya/revirathya-telbot:$(VERSION) .

docker-up:
	docker run \
	--rm -p 8080:8080 \
	--env-file .env \
	revirathya/revirathya-telbot:$(VERSION)

docker-scan:
	docker run --rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v financecore-trivy-cache:/root/.cache \
		-v $(CURDIR)/.trivyignore:/work/.trivyignore:ro \
		$(TRIVY_IMAGE) image --severity HIGH,CRITICAL --ignore-unfixed --exit-code 1 \
		--ignorefile /work/.trivyignore \
		revirathya/revirathya-telbot:$(VERSION)
