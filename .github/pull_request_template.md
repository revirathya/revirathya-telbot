# Pull Request

## Overview

<!-- What changes, and why. Link the issue or task if there is one. -->

## Checklist

- [ ] The PR title follows Conventional Commits: `type(scope): summary`, e.g. `feat(expense): add request_id`
      (types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `ci`, `build`, `perf`, `revert`; `!` marks a
      breaking change)
- [ ] `make check` pass
- [ ] If there is new image build, `make docker-scan` pass
- [ ] Docs updated: README (API behaviour), CHANGELOG, `docs/` if a decision changed (new ADR in `docs/adr/`)
- [ ] No secrets, tokens or DSNs in code, tests, logs or this description
