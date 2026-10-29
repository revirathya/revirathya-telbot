---
name: repo-security-review
description: Security review for the revirathya-telbot public repo. Use when asked to do a security review, a public-repo readiness check, or before opening/merging a PR, and when touching config, Makefile, infra/, .github/, or gRPC/Telegram client code. Covers hardcoded internal hostnames, IPs, and ports.
---

# Repo security review (revirathya-telbot)

This repo is public. Anything committed, including history, is visible to everyone.

## Rules

- **Never edit anything in this repo.** This is a read-only review: report findings, don't apply fixes.
- **Never give sample code** (snippets, diffs, patched config) unless the user explicitly asks for it. Describe
  the fix in words instead.
- The commands below are for running the review, not for showing to the user as fixes.

## Scope

- `cmd/`, `internal/` (especially `internal/config`, `internal/platform/grpc`, `internal/platform/telegram`, `internal/repository/grpc`)
- `infra/` (Dockerfile, compose files), `Makefile`, `.env.example`
- `.github/` (workflows, templates), `docs/`, `README.md`, `test/`
- Skip `.env`. It's gitignored and holds real local values. Never print its contents.

## Checks

### 1. Hardcoded internal hostnames, IPs, and ports

Internal endpoints reveal the network layout and must come from env vars (see `.env.example`), not code or committed config.

Search the tracked files only:

```bash
# IPv4 addresses (private ranges are the main concern)
git grep -nE '\b([0-9]{1,3}\.){3}[0-9]{1,3}\b' -- . ':!go.sum'

# host:port pairs and URLs with explicit ports
git grep -nE '[A-Za-z0-9.-]+:[0-9]{2,5}\b' -- . ':!go.sum'

# internal-looking hostnames
git grep -niE '\b[a-z0-9-]+\.(internal|local|lan|corp|svc|cluster\.local|home)\b'
git grep -niE '(revirathya|finance)[a-z0-9-]*\.[a-z]{2,}'

# ngrok / tunnel hosts
git grep -niE '[a-z0-9-]+\.(ngrok(-free)?\.(app|io|dev)|trycloudflare\.com)'
```

Allowed (don't report):
- `localhost`, `127.0.0.1`, `0.0.0.0`, and `:8080`-style listen addresses used as **defaults** for the local listen port.
- Public, well-known endpoints: `api.telegram.org`, `github.com`, `proxy.golang.org`, image registries.
- Placeholders such as `example.com`, `<host>`, `your-host:443`.
- Version strings that look like IPs (e.g. `1.2.3.4` in go.sum or tags).

Report:
- Any private IP (`10.x`, `172.16–31.x`, `192.168.x`), VPN/Tailscale IP (`100.64–127.x`), or real internal hostname.
- A real gRPC target (`FINANCE_HOST`-like value) set as a Go default, in the Makefile, Dockerfile, compose file, or workflow.
- A real ngrok/tunnel host or webhook URL committed anywhere.
- Ports of internal services hardcoded in code instead of being part of the host env var.

### 2. Secrets

```bash
gitleaks git --redact -v .   # full history
gitleaks dir --redact -v .   # working tree; findings in .env are expected, it's gitignored
```

- Any finding in history means the secret must be **rotated**, not just removed.
- `.env.example` values must be empty or obvious placeholders.
- Workflows must read secrets from `${{ secrets.* }}`, never inline values.

### 3. Ignored files

- `.gitignore` and `.dockerignore` both exclude `.env`. The Docker build context must not contain credentials.
- `git ls-files | grep -iE '\.env$|\.pem$|\.key$|credential|secret'` returns nothing except `.env.example`.

### 4. Logging

- The bot token, `FINANCE_API_KEY`, and full user messages must not be logged. The Telegram bot token is part of
  Telegram API URLs, so watch for logged request URLs and wrapped HTTP errors.

### 5. GitHub Actions (public repo)

- Workflows live in `.github/workflows/` (otherwise they never run).
- Top-level `permissions: contents: read`; widen per job only.
- No `pull_request_target` combined with checking out PR code.
- Third-party actions pinned to a full commit SHA (with a `# vX.Y.Z` comment). Check the owner name is spelled right to avoid typosquatted actions.

## Report format

For each finding: `file:line`, what was found (redact secrets and internal hosts, e.g. `10.x.x.x`, `fin****.internal`),
why it matters, and the fix described in words, with no code (usually: move it to an env var and document it in
`.env.example`). End with a short
summary: clean, or the number of findings by check.
