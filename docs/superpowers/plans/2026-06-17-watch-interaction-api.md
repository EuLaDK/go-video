# Watch Interaction API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add comment and danmaku APIs to the Go backend.

**Architecture:** Add a focused `internal/interaction` package with models, service, and PostgreSQL repository. Extend `internal/httpapi` with `/videos/{id}/comments` and `/videos/{id}/danmaku` routes while keeping existing video read routes stable.

**Tech Stack:** Go 1.26, PostgreSQL, `pgxpool`, standard `net/http`, `httptest`.

---

## File Structure

- Create: `internal/interaction/models.go`
- Create: `internal/interaction/service.go`
- Test: `internal/interaction/service_test.go`
- Create: `internal/interaction/postgres_repository.go`
- Create: `internal/httpapi/interaction_routes.go`
- Test: `internal/httpapi/interaction_routes_test.go`
- Modify: `internal/httpapi/server.go`
- Modify: `cmd/api/main.go`
- Create: `migrations/003_create_watch_interaction_tables.sql`
- Create: `seeds/003_seed_watch_interactions.sql`
- Modify: `README.md`
- Modify: `DEVELOPMENT.md`

## Tasks

### Task 1: Interaction Service

- [ ] Write failing tests for comment creation, sorting, like toggling, own-comment deletion, danmaku creation, and color fallback.
- [ ] Implement `internal/interaction/models.go` and `service.go` with injected clock and id generator.
- [ ] Run `go test ./internal/interaction`.

### Task 2: HTTP Routes

- [ ] Write failing route tests for list/add/like/delete comments and list/add danmaku.
- [ ] Add optional interaction service to `httpapi.Server`.
- [ ] Add method-aware interaction route dispatch before the existing `/videos/{id}` watch route.
- [ ] Run `go test ./internal/httpapi`.

### Task 3: PostgreSQL

- [ ] Add migration for `video_comments`, `video_comment_likes`, and `video_danmaku`.
- [ ] Add seed data for `xinghe`.
- [ ] Implement PostgreSQL repository.
- [ ] Run `go test ./...`.

### Task 4: Wiring And Docs

- [ ] Wire the interaction repository/service in `cmd/api/main.go`.
- [ ] Update README and DEVELOPMENT.md.
- [ ] Run `gofmt -w cmd internal`.
- [ ] Run `go test ./...`.
- [ ] Run `go run ./cmd/migrate`.
- [ ] Start service and verify comments and danmaku by HTTP.
