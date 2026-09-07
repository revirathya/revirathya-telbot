# ----- Builder -----
FROM golang:1.25-bookworm AS builder

# Install gPRC Tool
RUN go install github.com/bufbuild/buf/cmd/buf@v1.72.0 && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.12 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.2
ENV PATH="/root/go/bin:${PATH}"

WORKDIR /src

# Cache module downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy sources
COPY . .

# Generate proto stubs
RUN buf generate

# Static binary
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/telbot ./cmd/api

# -------------------

# ----- Runtime -----
FROM gcr.io/distroless/static-debian12:nonroot AS runtime
WORKDIR /app
COPY --from=builder /out/telbot ./telbot

ENTRYPOINT ["/app/telbot"]