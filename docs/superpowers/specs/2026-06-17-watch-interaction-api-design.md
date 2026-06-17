# Watch Interaction API 设计

更新日期：2026-06-17

## 目标

把播放页评论和弹幕从前端 localStorage 迁移到 Go 后端和 PostgreSQL，为后续前端 API facade 接入提供稳定接口。

## 范围

本阶段实现：

- 评论列表、新增评论、评论点赞切换、删除自己的评论。
- 弹幕列表、新增弹幕。
- 按视频 id 隔离评论和弹幕。
- 使用 `X-User-ID` 识别当前开发态用户，未传时使用 `demo-user`。
- PostgreSQL migration、seed、service 测试、HTTP 路由测试。

本阶段不实现：

- 多级评论、举报、审核、分页。
- 真实登录鉴权和权限角色。
- 弹幕按播放时间轴滚动定位。
- 前端 store 替换。

## API

- `GET /videos/{videoId}/comments?sort=latest|hot`
  - 返回评论数组，字段对齐前端 `WatchCommentItem`。
  - `latest` 按创建时间倒序，`hot` 按点赞数倒序后按时间倒序。
- `POST /videos/{videoId}/comments`
  - 请求体：`{ "content": "评论内容" }`
  - 返回新评论。
- `POST /videos/{videoId}/comments/{commentId}/like`
  - 当前用户未点赞则点赞，已点赞则取消。
  - 返回更新后的评论。
- `DELETE /videos/{videoId}/comments/{commentId}`
  - 只允许删除当前用户创建的评论。
- `GET /videos/{videoId}/danmaku`
  - 返回弹幕数组，字段对齐前端 `WatchDanmakuItem`。
- `POST /videos/{videoId}/danmaku`
  - 请求体：`{ "content": "弹幕内容", "color": "white|green|yellow|pink" }`
  - 返回新弹幕，颜色为空或非法时使用 `white`。

## 数据库

新增表：

- `video_comments`
  - `id`、`video_id`、`user_id`、`author`、`content`、`created_at_ms`。
- `video_comment_likes`
  - `comment_id`、`user_id`，用于计算 `likes` 和 `likedByMe`。
- `video_danmaku`
  - `id`、`video_id`、`user_id`、`content`、`color`、`created_at_ms`。

## 约束

- 空白评论和弹幕返回 `400`。
- 删除不存在或不属于自己的评论返回 `404`。
- `commentId` 和弹幕 id 由后端生成，格式不暴露为稳定契约。
- service 层注入 clock 和 id 生成函数，保证测试可重复。

## 验证标准

- `go test ./internal/interaction` 通过。
- `go test ./internal/httpapi` 通过。
- `go test ./...` 通过。
- `go run ./cmd/migrate` 可重复执行。
- 服务启动后可通过 HTTP 写入和读取评论、弹幕，并可点赞和删除自己的评论。
