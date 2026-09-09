# Revirathya - Telbot

## Overview
Telegram bot webhook dispatcher service for conversations and command handling.

**Stack:**
- Go 1.25
- go-chi/chi
- gRPC

## Prerequisites

- Go 1.25
- Go Tools: `buf`
- Docker

## Quickstart

```sh
cp .env.example .env # After copy, manually adjust .env values

make proto-generate # Generate gRPC stubs
make run            # Start services
```

## Configuration

|Var|Description|
|:--|:--|
|APP_ENV|(Optional) default: dev. Environment of deployment|
|APP_PORT|(Optional) default: 8080. Exposes port of services|
|FINCORE_GRPC_HOST|Host of Finance core gRPC|
|TELEGRAM_BOT_TOKEN|API token of Telegram bot|
|TELEGRAM_WEBHOOK_URL|Exposed URL of service for Telegram webhook to hit (public)|

## Docker

```sh
make docker-build
make docker-run
```

## Development

```sh
make proto-generate # buf: gRPC stubs
```

## Changelog

See [CHANGELOG.md](CHANGELOG.md).