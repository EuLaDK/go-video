# Next Video Go 后端开发记录

更新日期：2026-07-01

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
- 已实现播放配置接口第一版：
  - `GET /videos/{id}` 的响应新增 `playback` 字段。
  - `playback.sources` 已支持从独立播放源表读取；没有配置时继续从现有 `source_url` 和 `quality` 派生。
  - 已返回 `defaultQuality`、`requiresVip`、`canPlay`、`trialSeconds` 和 `message`，用于前端播放器控制播放源、清晰度和试看提示。
  - `canPlay` 已根据当前用户 VIP 状态计算；非 VIP 会员内容保留 360 秒试看，VIP 用户完整放行。
- 已新增多播放源持久化：
  - 新增 `video_playback_sources` 表。
  - 每条播放源保存 `quality`、`label`、`source_url`、`mime_type` 和 `display_order`。
  - `GET /videos/{id}` 优先返回表内多源配置，表内无记录时回退到 `videos.source_url`。
  - seed 已给 `xinghe` 配置 `4K HDR / 1080P / 720P` 三档示例播放源。
- 已接入真实 VIP 鉴权第一版：
  - 播放详情路由读取 `X-User-ID` 对应账号资料，按 `is_vip` 和 `vip_until` 生成播放鉴权状态。
  - 默认 `demo-user` 为非 VIP；新增 `demo-vip` 开发态用户用于本地验证 VIP 放行。
  - CORS 已允许 `X-User-ID` 请求头。
- 已实现前端 VIP 状态同步到后端第一版：
  - 新增 `POST /me/vip`，请求体为 `{"vipUntil":"YYYY-MM-DD"}`。
  - 接口会把当前 `X-User-ID` 对应用户写为 `is_logged_in=true`、`is_vip=true`，并更新 `vip_until`。
  - 前端 VIP 页面套餐购买已通过 `activateAccountVip` 写回后端，同时保留本地 fallback。
- 已实现断点续播服务端策略第一版：
  - `GET /videos/{id}` 的 `playback` 新增 `resume` 字段。
  - 播放详情路由会读取当前 `X-User-ID` 用户的 `/me/watch-history`，为当前视频提取最近有效观看秒数。
  - 前端播放页无 URL `episode/t` 参数时，优先使用后端 `playback.resume` 作为初始集数和播放秒数。
- 前端 `next-video` 播放页已读取播放配置：
  - `src/lib/video-api.ts` 的 mock fallback 返回同形 `playback` 配置。
  - `PlayerShell` 使用 `playback.sources` 作为 `<video>` 播放源，并提供清晰度选择。
- 已完成前端 HLS/DASH 播放适配第一版：
  - `next-video` 新增 `src/lib/player-source.ts`，按 `mimeType`/URL 将播放源识别为 `native`、`hls` 或 `dash`。
  - `PlayerShell` 对 MP4 和浏览器原生 HLS 继续直连 `video.src`，对非原生 HLS 动态加载 `hls.js`，对 DASH 动态加载 `dashjs`。
  - 切换清晰度或卸载播放器时会清理 HLS/DASH 实例和旧 `video.src`，避免旧播放源继续占用资源。

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
5. 继续增强真实播放器能力：
   - 前端 HLS/DASH 播放引擎适配已完成第一版。
   - 继续完善播放页冒烟测试和组件级测试。
6. 增加后台内容管理：
   - 视频新增和编辑。
   - 上下架。
   - 频道配置。
   - 推荐位配置。

## 下一步准备做

### 下次继续开发入口

本次已完成 **真实 VIP 鉴权 v1**、**前端 VIP 状态同步 v1**、**断点续播服务端策略 v1** 和 **前端 HLS/DASH 播放适配 v1**。下次建议从播放页冒烟测试和组件级测试继续，不急着做后台内容管理。

优先顺序：

1. 播放页冒烟和组件测试：
   - 覆盖播放详情读取 `playback.sources`、`playback.resume`、清晰度切换和初始恢复秒数。
   - 验证 MP4、本地 mock fallback，以及未来 HLS/DASH 示例源不会破坏播放页加载。

2. 继续补真实用户体系：
   - 登录注册。
   - 用户资料持久化。
   - VIP 状态续期和过期处理。

3. 本地联调时，在 `next-video` 前端配置：

```text
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

4. 前端 localStorage 状态替换已完成第一轮联调：
   - 用户资料、追剧收藏、观看历史已接入 API facade。
   - 首页、播放页、收藏页、历史页已通过本地 Next dev server 请求验证。

5. 评论和弹幕服务端化第一版已完成，前端 `use-watch-interaction-store` 已接入 API facade。
6. VIP 套餐购买状态同步已完成第一版：
   - 后端新增 `POST /me/vip`。
   - 前端 `use-user-store.activateVip` 已优先写本地状态，再同步 Go API。
7. 断点续播服务端策略已完成第一版：
   - 后端 `GET /videos/{id}` 已返回 `playback.resume`。
   - 前端播放页已在无 URL 秒数时使用后端恢复点。

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

### 2026-06-25 播放配置 v1

- `go test ./internal/video`：通过，已覆盖 `WatchPage` 返回 `playback.sources`、默认清晰度、会员试看和可播放状态。
- `go test ./internal/httpapi`：通过，播放详情响应新增字段未破坏现有 HTTP 路由。
- `go test ./...`：通过。
- `node --test src\lib\video-api.test.mjs`：通过，播放页 mock fallback 已返回同形 `playback` 配置。
- `node --test src\lib\*.test.mjs`：通过，89 个前端 lib 测试全部通过。
- `npm.cmd run lint`：通过。

### 2026-06-25 多播放源持久化 v1

- `go test ./internal/video`：通过，已覆盖 repository 多播放源优先于旧 `source_url` 兜底。
- `go test ./cmd/api`：通过，`PostgresRepository` 已满足视频服务接口。
- `go test ./...`：通过。
- `go run ./cmd/migrate`：通过，已执行 `004_create_video_playback_sources.sql` 和 `004_seed_video_playback_sources.sql`。
- 使用 8091 临时服务验证 `GET /videos/xinghe`：返回 `4K HDR`、`1080P`、`720P` 三个 `playback.sources`，`defaultQuality` 为 `4K HDR`，`trialSeconds` 为 `360`。

### 2026-06-25 真实 VIP 鉴权 v1

- `go test ./internal/video`：通过，已覆盖非 VIP 用户观看会员内容时 `canPlay=false`，VIP 用户 `canPlay=true` 且 `trialSeconds=0`。
- `go test ./internal/httpapi`：通过，已覆盖播放详情路由读取 `X-User-ID` 对应账号资料并注入播放鉴权状态。
- `go test ./...`：通过。
- `go run ./cmd/migrate`：通过，已写入 `demo-vip` 开发态用户。
- 使用 8091 临时服务验证默认 `demo-user` 请求 `GET /videos/xinghe`：返回 `requiresVip=true`、`canPlay=false`、`trialSeconds=360`。
- 使用 8091 临时服务验证 `X-User-ID: demo-vip` 请求 `GET /videos/xinghe`：返回 `requiresVip=true`、`canPlay=true`、`trialSeconds=0`。

### 2026-06-25 前端 VIP 状态同步 v1

- `go test ./internal/account`：通过，已覆盖 VIP 开通保留现有用户资料、缺失用户创建默认资料。
- `go test ./internal/httpapi`：通过，已覆盖 `POST /me/vip` 读取 `X-User-ID` 并返回更新后的用户资料。
- `go test -count=1 ./...`：通过，VIP 状态写入未破坏现有视频、账号、互动接口测试。
- `node --test src\lib\account-api.test.mjs`：通过，已覆盖 `activateAccountVip` 向 `/me/vip` 发送 `vipUntil`。
- `node --test src\stores\account-sync-stores.test.mjs`：通过，已覆盖 `activateVip` 乐观更新并同步后端响应。
- `node --test src\lib\account-api.test.mjs src\stores\account-sync-stores.test.mjs`：通过，8 个账号同步相关用例全部通过。
- `npm.cmd run lint`：通过。

### 2026-06-25 断点续播服务端策略 v1

- `go test ./internal/video`：通过，已覆盖 `playback.resume` 从播放上下文进入播放详情响应。
- `go test ./internal/httpapi`：通过，已覆盖播放详情路由读取用户观看历史并注入当前视频恢复点。
- `go test -count=1 ./...`：通过，断点续播恢复点未破坏现有后端接口测试。
- `node --test src\lib\video-api.test.mjs`：通过，已覆盖 mock fallback 返回同形 `playback.resume`。
- `node --test src\lib\playback-resume.test.mjs`：通过，已覆盖 URL 参数优先、无参数时使用后端恢复点。
- `node --test src\lib\*.test.mjs`：通过，92 个前端 lib 用例全部通过。
- `npm.cmd run lint`：通过。

### 2026-07-01 前端 HLS/DASH 播放适配 v1

- `node --test src\lib\player-source.test.mjs`：通过，已覆盖 MP4、HLS、DASH 和未知类型的播放引擎识别。
- `node --test src\lib\player-source.test.mjs src\lib\player-controls.test.mjs src\lib\video-api.test.mjs`：通过，16 个播放器和视频 API 相关用例全部通过。
- `node --test src\lib\*.test.mjs`：通过，97 个前端 lib 用例全部通过。
- `npm.cmd run lint`：通过。

### 2026-07-01 播放器控制修复 v1

- 前端 `PlayerShell` 已把清晰度、倍速和弹幕速度从原生 `select` 改为现有 shadcn/Radix `DropdownMenuRadioGroup`，修复深色播放器里 option 文字不可读的问题。
- 修复播放源切换 key：同一集同一 `sourceUrl` 不再因为清晰度标签变化重建 video 节点；切换下一集时会重新绑定当前 video 源。
- `play()` 被切源或暂停打断时的 `AbortError` 已被安全处理，不再形成未处理 Promise 错误。
- 新增 `test:e2e` 和播放页 Playwright 冒烟入口；该入口默认连接手动启动的 `http://127.0.0.1:3000`，不再自动启动 Next 服务。
- `node --test src\lib\player-source.test.mjs src\lib\player-controls.test.mjs src\lib\video-api.test.mjs`：通过，18 个播放器和视频 API 相关用例全部通过。
- `npx.cmd playwright test --list`：通过，可识别播放页冒烟用例。
- `npm.cmd run lint`：通过。
