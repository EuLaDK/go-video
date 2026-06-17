# Next Video Go 后端开发记录

更新日期：2026-06-17

## 当前完成

- 已初始化 Go 模块 `next-video-golang`。
- 已确定第一版范围：视频只读 API、PostgreSQL schema、seed 数据、前端 API facade 对齐。
- 已添加设计文档：`docs/superpowers/specs/2026-06-16-video-read-api-design.md`。
- 已添加实施计划：`docs/superpowers/plans/2026-06-16-video-read-api.md`。
- 已实现配置读取：
  - `PORT` 默认 `8080`。
  - `DATABASE_URL` 默认连接本地 `nextvideo`。
- 已实现视频领域服务：
  - 首页聚合。
  - 排行榜排序和频道筛选。
  - 频道页筛选。
  - 搜索页筛选。
  - 播放详情和相关推荐。
  - 视频 id 列表。
- 已实现 HTTP API：
  - `GET /health`
  - `GET /videos/home`
  - `GET /videos/rank`
  - `GET /videos/channel/{slug}`
  - `GET /videos/search`
  - `GET /videos/{id}`
  - `GET /videos/ids`
- 已实现 PostgreSQL 连接和 repository。
- 已添加数据库 migration 和 seed SQL。
- 已添加 `cmd/migrate`，可直接执行 migration 和 seed。
- 已实现开发态用户、收藏和观看历史 API：
  - `GET /me`
  - `POST /me/login`
  - `POST /me/logout`
  - `GET /me/favorites`
  - `POST /me/favorites`
  - `DELETE /me/favorites/{videoId}`
  - `GET /me/watch-history`
  - `POST /me/watch-history`
  - `DELETE /me/watch-history/{videoId}?episode=1`
  - `DELETE /me/watch-history`
- 已新增账号状态相关表：
  - `users`
  - `user_favorites`
  - `user_watch_history`
- 已添加开发态默认用户 `demo-user`。
- 前端 `next-video` 已开始接入 `/me` 系列接口：
  - 新增 `src/lib/account-api.ts`。
  - `use-user-store.ts`、`use-favorite-store.ts`、`use-watch-history-store.ts` 已优先同步 Go API，并保留本地 fallback。
- 已实现视频互动 API 第一版：
  - `GET /videos/{id}/comments?sort=latest|hot`
  - `POST /videos/{id}/comments`
  - `POST /videos/{id}/comments/{commentId}/like`
  - `DELETE /videos/{id}/comments/{commentId}`
  - `GET /videos/{id}/danmaku`
  - `POST /videos/{id}/danmaku`
- 已新增互动数据相关表：
  - `video_comments`
  - `video_comment_likes`
  - `video_danmaku`

## 后续计划

1. 联调前端，把 `NEXT_PUBLIC_API_BASE_URL` 指向 Go 服务。
2. 将前端 localStorage 状态逐步替换为服务端接口：
   - 用户资料。
   - 收藏。
   - 观看历史。
3. 增加真实用户系统：
   - 登录注册。
   - 用户资料。
   - VIP 状态持久化。
4. 增加互动数据服务端化：
   - 评论。
   - 弹幕。
5. 增加真实播放器能力：
   - HLS/DASH 播放源。
   - 清晰度切换。
   - 试看和会员鉴权。
   - 断点续播同步。
6. 增加后台内容管理：
   - 视频新增和编辑。
   - 上下架。
   - 频道配置。
   - 推荐位配置。

## 下一步准备做

### 下次继续开发入口

今天先暂停在这里。下次建议从 **前端联调已完成的 Go API** 开始，不急着继续扩后端新能力。

优先顺序：

1. 在 `next-video` 前端配置：

```text
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

2. 前端 localStorage 状态替换已完成第一轮联调：
   - 用户资料、追剧收藏、观看历史已接入 API facade。
   - 首页、播放页、收藏页、历史页已通过本地 Next dev server 请求验证。

3. 评论和弹幕服务端化第一版已完成，前端 `use-watch-interaction-store` 已接入 API facade。下一步建议做真实播放器接口能力。

当前后端服务上次验证已启动在：

```text
http://localhost:8080
PID: 21260
```

如果下次端口被占用，先检查或停止旧服务：

```powershell
netstat -ano | findstr :8080
Stop-Process -Id 21260
```

### 常规启动步骤

1. 执行：

```powershell
go run ./cmd/migrate
```

2. 启动后端：

```powershell
go run ./cmd/api
```

3. 在前端配置：

```text
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

4. 打开首页、排行榜、频道页、搜索页和播放页，确认数据来自 Go 后端。

5. 调用 `/me`、`/me/favorites` 和 `/me/watch-history`，确认用户状态接口可用。

## 验证记录

- `go test ./internal/config`：通过。
- `go test ./internal/video`：通过。
- `go test ./internal/httpapi`：通过。
- `go test ./internal/account`：通过。
- `go test ./cmd/migrate`：通过。
- `go test ./...`：通过。
- `go run ./cmd/migrate`：通过，已执行 migration 和 seed。
- `GET http://127.0.0.1:8080/health`：返回 `{"status":"ok"}`。
- `GET http://127.0.0.1:8080/videos/ids`：返回 14 个视频 id。
- `GET http://127.0.0.1:8080/me`：返回 `demo-user` 用户资料。
- `POST http://127.0.0.1:8080/me/login`：UTF-8 JSON 请求体可正确写入中文昵称。
- `POST http://127.0.0.1:8080/me/favorites`：可写入 `xinghe` 收藏。
- `GET http://127.0.0.1:8080/me/favorites`：可读取 `xinghe` 收藏。
- `POST http://127.0.0.1:8080/me/watch-history`：可写入 `xinghe` 第 2 集观看历史。
- `GET http://127.0.0.1:8080/me/watch-history`：可读取 `xinghe` 第 2 集观看历史。
- 当前后端服务已启动在 `http://localhost:8080`，进程名 `api`，PID `21260`。
- `node --test src\lib\account-api.test.mjs src\stores\account-sync-stores.test.mjs`：通过，6 个账号同步测试全部通过。
- `npm.cmd run lint`：通过。
- `GET http://127.0.0.1:3000/`、`/profile/favorites`、`/profile/history`、`/watch/xinghe`：均返回 200。
- 使用 `X-User-ID: codex-smoke` 验证 `/me/login`、`/me/favorites`、`/me/watch-history` 写入读取链路：通过，测试收藏和历史已清理。
- `go test ./internal/interaction`：通过。
- `go test ./internal/httpapi`：通过，已覆盖评论和弹幕路由。
- `go test ./...`：通过，已包含视频互动模块。
- `go run ./cmd/migrate`：通过，已执行 `003_create_watch_interaction_tables.sql` 和 `003_seed_watch_interactions.sql`。
- 使用 8083 临时服务验证视频互动接口：`/health`、评论新增、评论点赞、评论列表、删除自己的评论、弹幕新增、弹幕列表均通过。
- `node --test src\lib\interaction-api.test.mjs src\stores\watch-interaction-sync-store.test.mjs`：通过，6 个评论/弹幕前端同步测试全部通过。
