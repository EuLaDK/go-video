# Account State API 设计

更新日期：2026-06-17

## 目标

在现有视频只读 API 基础上，增加开发态用户、收藏和观看历史接口，为前端后续把 localStorage 状态迁移到服务端提供稳定落点。

## 范围

本阶段实现：

- 开发态用户资料读取、登录和退出。
- 当前用户收藏列表的读取、新增和删除。
- 当前用户观看历史的读取、新增、删除单条和清空。
- PostgreSQL migration 和 seed。
- HTTP API、service 和 repository 测试。

本阶段不实现：

- 密码、短信验证码、JWT、刷新 token。
- 多租户权限模型。
- 评论、弹幕、点赞、缓存状态。
- 前端 store 替换。

## 用户识别

开发期优先使用请求头：

```text
X-User-ID: demo-user
```

未传时默认使用 `demo-user`。这样前端可以先无登录成本联调，后续真实鉴权上线时再替换用户解析逻辑。

## API

- `GET /me`
  - 返回当前用户资料，字段对齐前端 `UserProfileState`。
- `POST /me/login`
  - 请求体：`{ "nickname": "用户", "contact": "phone-or-email", "avatarUrl": "" }`
  - 返回登录后的用户资料。
- `POST /me/logout`
  - 返回退出后的用户资料。
- `GET /me/favorites`
  - 返回收藏数组，字段对齐前端 `FavoriteItem`。
- `POST /me/favorites`
  - 请求体为不含 `addedAt` 的收藏摘要。
  - 返回写入后的收藏项。
- `DELETE /me/favorites/{videoId}`
  - 删除指定视频收藏。
- `GET /me/watch-history`
  - 返回观看历史数组，字段对齐前端 `WatchHistoryItem`。
- `POST /me/watch-history`
  - 请求体为不含 `watchedAt` 的观看历史摘要。
  - 返回写入后的历史项。
- `DELETE /me/watch-history/{videoId}?episode=1`
  - 删除指定视频历史；episode 可选，不传则删除该视频全部历史。
- `DELETE /me/watch-history`
  - 清空当前用户观看历史。

## 数据库

新增表：

- `users`
  - `id`、`avatar_url`、`email`、`is_logged_in`、`is_vip`、`nickname`、`phone`、`vip_until`。
- `user_favorites`
  - `user_id`、`video_id`、收藏摘要字段、`added_at`。
- `user_watch_history`
  - `user_id`、`video_id`、`episode`、历史摘要字段、播放秒数、总秒数、`watched_at`。

## 验证标准

- `go test ./...` 通过。
- `go run ./cmd/migrate` 可重复执行。
- 服务启动后：
  - `GET /me` 返回默认开发用户。
  - `POST /me/favorites` 后 `GET /me/favorites` 能看到写入项。
  - `POST /me/watch-history` 后 `GET /me/watch-history` 能看到写入项。

## 自检

- 范围只覆盖用户资料、收藏、观看历史。
- API 字段对齐当前前端 store 类型。
- 用户识别是开发态实现，后续可替换为真实鉴权。
