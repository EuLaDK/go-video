# Video Read API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the first Go backend for Next Video, backed by PostgreSQL, serving the frontend video read API.

**Architecture:** A small layered Go service using `net/http`, `pgxpool`, and focused internal packages. HTTP handlers delegate to a video service, the service owns query behavior, and the PostgreSQL repository owns persistence.

**Tech Stack:** Go 1.26, PostgreSQL, `github.com/jackc/pgx/v5`, standard `net/http`, standard `httptest`.

---

## File Structure

- Create: `go.mod` and `go.sum`
- Create: `cmd/api/main.go`
- Create: `internal/config/config.go`
- Test: `internal/config/config_test.go`
- Create: `internal/database/database.go`
- Create: `internal/video/models.go`
- Create: `internal/video/service.go`
- Test: `internal/video/service_test.go`
- Create: `internal/video/postgres_repository.go`
- Create: `internal/httpapi/server.go`
- Test: `internal/httpapi/server_test.go`
- Create: `migrations/001_create_video_read_tables.sql`
- Create: `seeds/001_seed_video_read_data.sql`
- Create: `README.md`
- Create: `DEVELOPMENT.md`

### Task 1: Go Module And Dependencies

**Files:**
- Create: `go.mod`
- Create: `go.sum`

- [ ] **Step 1: Initialize module**

Run:

```powershell
go mod init next-video-golang
```

Expected: `go.mod` exists with module path `next-video-golang`.

- [ ] **Step 2: Add PostgreSQL driver**

Run:

```powershell
go get github.com/jackc/pgx/v5/pgxpool
```

Expected: `go.mod` contains `github.com/jackc/pgx/v5`.

- [ ] **Step 3: Verify empty module**

Run:

```powershell
go test ./...
```

Expected: packages either report no test files or pass.

### Task 2: Configuration

**Files:**
- Create: `internal/config/config_test.go`
- Create: `internal/config/config.go`

- [ ] **Step 1: Write failing config tests**

Test default values and environment overrides:

```go
func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")
	cfg := config.Load()
	if cfg.Port != "8080" {
		t.Fatalf("Port = %q, want 8080", cfg.Port)
	}
	if cfg.DatabaseURL == "" {
		t.Fatal("DatabaseURL should have a development default")
	}
}
```

Run:

```powershell
go test ./internal/config
```

Expected: fail because package does not exist.

- [ ] **Step 2: Implement config**

Create `Config`, `Load`, and helper environment readers. Defaults:

```text
PORT=8080
DATABASE_URL=postgres://postgres:dengke258567@localhost:5432/nextvideo?sslmode=disable
```

- [ ] **Step 3: Verify config**

Run:

```powershell
go test ./internal/config
```

Expected: pass.

### Task 3: Video Service Models And Query Behavior

**Files:**
- Create: `internal/video/models.go`
- Create: `internal/video/service_test.go`
- Create: `internal/video/service.go`

- [ ] **Step 1: Write failing service tests**

Cover the service contract:

```go
func TestServiceRankVideosSortsByScore(t *testing.T) {
	svc := video.NewService(fakeRepositoryWithSampleData())
	got, err := svc.RankVideos(context.Background(), video.RankQuery{Sort: "score", Channel: "movie"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("expected ranked videos")
	}
	if got[0].Score < got[len(got)-1].Score {
		t.Fatalf("videos are not sorted by score desc")
	}
}
```

Run:

```powershell
go test ./internal/video
```

Expected: fail because package does not exist.

- [ ] **Step 2: Implement models and service**

Define these exported types with JSON tags matching the frontend:

```go
type Video struct {
	ID              string                `json:"id"`
	Title           string                `json:"title"`
	Subtitle        string                `json:"subtitle"`
	Description     string                `json:"description"`
	Score           string                `json:"score"`
	Heat            string                `json:"heat"`
	Update          string                `json:"update"`
	Category        string                `json:"category"`
	Year            string                `json:"year"`
	Region          string                `json:"region"`
	TotalEpisodes   int                   `json:"totalEpisodes"`
	Quality         string                `json:"quality"`
	Badge           string                `json:"badge"`
	Progress        string                `json:"progress"`
	Duration        string                `json:"duration"`
	SourceURL       string                `json:"sourceUrl"`
	CoverGradient   string                `json:"coverGradient"`
	Tags            []string              `json:"tags"`
	CastNames       []string              `json:"castNames"`
	ReleaseCalendar []ReleaseCalendarItem `json:"releaseCalendar"`
	Episodes        []Episode             `json:"episodes"`
	RelatedVideoIDs []string              `json:"relatedVideoIds"`
}
```

Add repository interface methods for loading channels, videos, and related videos. Keep query logic in service so it can be tested without a database.

- [ ] **Step 3: Verify service**

Run:

```powershell
go test ./internal/video
```

Expected: pass.

### Task 4: PostgreSQL Schema And Seed Data

**Files:**
- Create: `migrations/001_create_video_read_tables.sql`
- Create: `seeds/001_seed_video_read_data.sql`
- Create: `internal/database/database.go`
- Create: `internal/video/postgres_repository.go`

- [ ] **Step 1: Create migration**

Create tables `channels`, `videos`, `video_episodes`, `video_release_calendar`, and `video_related` with foreign keys and stable ordering columns.

- [ ] **Step 2: Create seed SQL**

Insert the current frontend channel list and video library. Seed SQL uses `ON CONFLICT DO UPDATE` so it can be rerun during development.

- [ ] **Step 3: Implement database connection and repository**

Use `pgxpool.New`, `Ping`, and typed row scanning. Repository returns complete `Video` values with nested episodes, calendar items, and related ids.

- [ ] **Step 4: Verify package build**

Run:

```powershell
go test ./internal/database ./internal/video
```

Expected: pass.

### Task 5: HTTP API

**Files:**
- Create: `internal/httpapi/server_test.go`
- Create: `internal/httpapi/server.go`
- Create: `cmd/api/main.go`

- [ ] **Step 1: Write failing handler tests**

Use `httptest` and a fake video service to verify:

```go
func TestServerHomeRouteReturnsHomeData(t *testing.T) {
	srv := httpapi.NewServer(fakeVideoService{})
	req := httptest.NewRequest(http.MethodGet, "/videos/home", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}
```

Run:

```powershell
go test ./internal/httpapi
```

Expected: fail because package does not exist.

- [ ] **Step 2: Implement routes**

Implement:

```text
GET /health
GET /videos/home
GET /videos/rank
GET /videos/channel/{slug}
GET /videos/search
GET /videos/{id}
GET /videos/ids
```

Add JSON helpers, CORS headers, and method checks.

- [ ] **Step 3: Implement main**

Wire config, database, repository, service, and HTTP server.

- [ ] **Step 4: Verify HTTP API tests**

Run:

```powershell
go test ./internal/httpapi ./...
```

Expected: pass.

### Task 6: Docs And Development Record

**Files:**
- Create: `README.md`
- Create: `DEVELOPMENT.md`

- [ ] **Step 1: Write README**

Include local run commands, database connection defaults, migration/seed execution guidance, and frontend integration:

```text
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

- [ ] **Step 2: Write development record**

Record:

- 当前完成：项目初始化、设计、接口范围、数据库方案。
- 后续计划：用户系统、互动数据、HLS/DASH、后台管理。
- 下一步：跑 migration/seed，启动后端，联调前端。

- [ ] **Step 3: Final verification**

Run:

```powershell
gofmt -w cmd internal
go test ./...
```

Expected: all packages pass.
