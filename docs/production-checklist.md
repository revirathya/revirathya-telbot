# Telegram Bot Service: Production Checklist

Checklist for a Go service that receives Telegram webhooks, handles conversations, sends/edits messages, and communicates with other gRPC services.

> **Top priorities:** scan git history for secrets before going public, verify the webhook secret, deduplicate by `update_id`, and set timeouts on all gRPC calls.

---

## ✅ Ready now (no deployment target needed)

### Minimal production setup

#### Webhook handling
- [ ] Verify the `X-Telegram-Bot-Api-Secret-Token` header on every incoming request and reject mismatches, so nobody can forge updates to your endpoint.
- [ ] Return `200 OK` quickly and process the update asynchronously (worker pool or queue). Telegram retries slow or failed responses, which causes duplicates.
- [ ] Deduplicate by `update_id` so retried updates aren't processed twice. This matters most for actions with side effects, such as gRPC writes.
- [ ] Recover from panics in handlers so one bad update can't crash the service.

#### Conversation state
- [ ] Store conversation state in a persistent store (Redis or a DB) rather than in memory, or it's lost on restart and you can't run more than one replica.
- [ ] Add expiry/TTL for abandoned conversations.

#### Telegram API calls
- [ ] Respect rate limits (roughly 30 messages/sec overall, about 1/sec per chat, about 20/min per group), and handle `429` responses by honoring `retry_after`.
- [ ] Handle common errors gracefully: user blocked the bot (`403`), "message is not modified" on edits, and message too old to edit.

#### gRPC clients
- [ ] Set a timeout on every call with `context.WithTimeout`. A hung downstream service shouldn't hang webhook processing.
- [ ] Make the insecure/TLS flag default to secure.
- [ ] Decide on retry behavior for transient errors (`Unavailable`, `DeadlineExceeded`), and only retry idempotent calls.

#### Service basics
- [ ] Load all config from environment variables, and fail fast at startup if required ones (bot token, API keys) are missing.
- [ ] Use structured logging (`log/slog`). Never log the bot token, API keys, or full user messages if they may contain personal data.
- [ ] Add a `/healthz` endpoint (liveness) and a `/readyz` endpoint (readiness, e.g., checks the state store).
- [ ] Implement graceful shutdown: on `SIGTERM`, stop accepting requests, finish in-flight work, then close gRPC connections.
- [ ] Write a multi-stage Dockerfile that produces a small final image and runs as a non-root user.

#### Code quality
- [ ] Get `golangci-lint` passing.
- [ ] Run tests with `go test -race ./...`.
- [ ] Run `govulncheck ./...` to catch vulnerable dependencies.

### Documents
- [ ] **README:** what the bot does, how to run it locally, and the required env vars.
- [ ] **Configuration reference:** every env var, whether it's required, its default, and an example.
- [ ] **Architecture overview:** Telegram → webhook → processing → gRPC services, ideally with a simple diagram.
- [ ] **Conversation flows:** states and transitions for each multi-step conversation.
- [ ] **External dependencies:** which gRPC services are called, which methods, and what happens when each one is down.
- [ ] **Local development guide:** how to receive webhooks locally (e.g., via a tunnel such as ngrok or cloudflared) and how to run with insecure gRPC.

### GitHub public repo
- [ ] **Scan the full git history for secrets** with gitleaks or trufflehog. If a bot token or API key was *ever* committed, rotate it. Deleting it in a later commit doesn't remove it from history.
- [ ] Add `.env`, credential files, and local configs to `.gitignore`, and commit a `.env.example` with placeholder values instead.
- [ ] Remove hardcoded internal hostnames, IPs, and ports from code and configs.
- [ ] Decide whether the `.proto` files of internal services should be public. They reveal your internal API surface.
- [ ] Add a `LICENSE`. Without one, the code is public but legally not reusable.
- [ ] Add a `SECURITY.md` explaining how to report vulnerabilities privately.
- [ ] Turn on GitHub secret scanning and push protection, plus Dependabot for Go modules and GitHub Actions.
- [ ] Set up branch protection on `main`: require PRs and passing CI.
- [ ] Add a CI workflow (GitHub Actions) that runs lint, tests, `govulncheck`, and a Docker build.

---

## ⏳ Waiting for deployment targets

### Minimal production setup
- [ ] **Public HTTPS endpoint:** a domain and valid TLS certificate. Telegram only delivers webhooks over HTTPS on ports 443, 80, 88, or 8443.
- [ ] **Register the webhook** with `setWebhook`, including `secret_token`, `allowed_updates` (only the update types you handle), and a suitable `max_connections`. Confirm with `getWebhookInfo`.
- [ ] **Secrets management:** put the bot token and API keys in the platform's secret store, not in plain env files on the server.
- [ ] **TLS for gRPC** between the bot and backend services, with the insecure flag off.
- [ ] **Network access:** make sure the bot can reach the gRPC services, and that those services aren't exposed publicly.
- [ ] **Health probes:** wire `/healthz` and `/readyz` into the platform (Kubernetes probes, load balancer checks, etc.).
- [ ] **Resource limits:** set CPU/memory limits and a restart policy.
- [ ] **Hosted state store:** a managed Redis or DB instance, with backups if conversation state matters.
- [ ] **Observability:** ship logs somewhere searchable, and add basic metrics (updates received, processing errors, gRPC latency and errors, Telegram `429`s).
- [ ] **Alerts:** at minimum for the service being down, a rising error rate, and webhook delivery errors (`getWebhookInfo` reports `last_error_message`).
- [ ] **Deployment pipeline:** build and push the image, deploy, and have a documented rollback path.

### Documents
- [ ] **Deployment guide:** how to deploy, where secrets live, and how to register or update the webhook.
- [ ] **Environment matrix:** dev, staging, and prod URLs, which bot token each uses, and which gRPC endpoints each points to. Use a separate test bot for non-prod.
- [ ] **Runbook:** what to do when Telegram stops delivering, a gRPC dependency is down, or a token needs rotating.
- [ ] **Monitoring reference:** where the logs and dashboards are, and what each alert means.

### GitHub public repo
- [ ] Add deployment secrets as GitHub Actions secrets or environment secrets, never in workflow files.
- [ ] If CI deploys, use GitHub Environments with required reviewers for production.
- [ ] Make sure workflows triggered by pull requests from forks can't access deployment secrets. This is important for public repos.
