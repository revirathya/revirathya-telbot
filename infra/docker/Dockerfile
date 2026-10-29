# ----- Builder -----
FROM golang:1.26.10-bookworm AS builder

WORKDIR /src

# Cache module downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy sources
COPY . .

# Static binary, stamped with the image version
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /out/telbot ./cmd/api

# -------------------

# ----- Runtime -----
FROM gcr.io/distroless/static-debian12:nonroot AS runtime
WORKDIR /app
COPY --from=builder /out/telbot ./telbot

ENTRYPOINT ["/app/telbot"]