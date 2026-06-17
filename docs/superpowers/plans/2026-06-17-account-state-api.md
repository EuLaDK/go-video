# Account State API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add development user, favorites, and watch history APIs to the Go backend.

**Architecture:** Add a focused `internal/account` package with models, service, and PostgreSQL repository. Extend `internal/httpapi` with `/me` routes while keeping existing video routes stable.

**Tech Stack:** Go 1.26, PostgreSQL, `pgxpool`, standard `net/http`, `httptest`.

---

## File Structure

- Create: `internal/account/models.go`
- Create: `internal/account/service.go`
- Test: `internal/account/service_test.go`
- Create: `internal/account/postgres_repository.go`
- Modify: `internal/httpapi/server.go`
- Create: `internal/httpapi/account_routes.go`
- Test: `internal/httpapi/account_routes_test.go`
- Modify: `cmd/api/main.go`
- Create: `migrations/002_create_account_state_tables.sql`
- Create: `seeds/002_seed_demo_user.sql`
- Modify: `README.md`
- Modify: `DEVELOPMENT.md`

## Tasks

### Task 1: Account Service

- [ ] Write failing tests for login/logout/profile/favorites/history behavior.
- [ ] Implement account models and service with an injectable clock.
- [ ] Run `go test ./internal/account`.

### Task 2: Database

- [ ] Add migration for `users`, `user_favorites`, `user_watch_history`.
- [ ] Add seed for `demo-user`.
- [ ] Implement PostgreSQL repository.
- [ ] Run `go test ./internal/account`.

### Task 3: HTTP Routes

- [ ] Write failing `httptest` coverage for `/me`, `/me/favorites`, `/me/watch-history`.
- [ ] Extend server with optional account service and method-aware routing.
- [ ] Run `go test ./internal/httpapi`.

### Task 4: Wiring And Docs

- [ ] Wire account repository/service in `cmd/api/main.go`.
- [ ] Update README and DEVELOPMENT.md.
- [ ] Run `gofmt -w cmd internal`.
- [ ] Run `go test ./...`.
- [ ] Run `go run ./cmd/migrate`.
- [ ] Start service and verify `/me`, `/me/favorites`, `/me/watch-history`.
