# Next Video Go 后端第一版设计

更新日期：2026-06-16

## 目标

在 `next-video-golang` 中初始化 Go 后端项目，使用本地 PostgreSQL 数据库 `nextvideo` 提供当前前端 API facade 需要的视频只读接口。第一版聚焦视频播放和列表数据，不包含登录、收藏、观看历史、评论、弹幕写入等用户互动能力。

## 范围

本阶段实现：

- Go 项目骨架、配置读取、HTTP 服务启动。
- PostgreSQL 连接池。
- 视频、频道、选集、更新时间线、相关推荐的表结构。
- 可重复执行的 seed SQL，把前端 mock 数据落入数据库。
- 与前端 `src/lib/video-api.ts` 对齐的 JSON 接口。
- `DEVELOPMENT.md` 开发记录，记录当前完成、后续计划和下一步。

本阶段不实现：

- 用户注册登录、JWT、权限鉴权。
- 收藏、缓存、观看历史、评论、弹幕的服务端持久化。
- 视频文件上传、转码、HLS/DASH 切片。
- 后台内容管理页面。

## API 契约

服务默认监听 `http://localhost:8080`，前端可设置：

```text
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

接口列表：

- `GET /health`
  - 返回服务和数据库连接状态。
- `GET /videos/home`
  - 返回 `{ featuredVideo, hotVideos, rankVideos, recommendationVideos }`。
- `GET /videos/rank?sort=&channel=`
  - 返回视频数组。
  - `sort` 支持 `hot`、`score`、`new`、`rising`、`reputation`、`vip`，默认 `hot`。
  - `channel` 为空、`all` 或 `featured` 时不过滤。
- `GET /videos/channel/{slug}?type=&year=&sort=`
  - 返回 `{ channel, heroVideo, videos }`。
  - 未命中频道时回退到 `featured`。
- `GET /videos/search?q=&channel=&quality=&type=&year=&sort=`
  - 返回 `{ hotSearchKeywords, recommendationVideos, videos }`。
  - 空关键词返回空 `videos`。
- `GET /videos/{id}`
  - 返回 `{ video, relatedVideos }`。
  - 未命中视频时回退到首页主推视频。
- `GET /videos/ids`
  - 返回视频 id 字符串数组。

JSON 字段采用 camelCase，对齐前端 `VideoItem`、`ChannelItem`、`VideoEpisode` 和 `VideoReleaseCalendarItem` 类型。

## 数据库设计

使用 PostgreSQL 默认端口 `5432`，数据库名 `nextvideo`。本地连接默认值：

```text
DATABASE_URL=postgres://postgres:dengke258567@localhost:5432/nextvideo?sslmode=disable
```

表结构：

- `channels`
  - `slug` 主键。
  - `label`、`description`、`keywords`、`accent`、`display_order`。
- `videos`
  - `id` 主键。
  - 前端 `VideoItem` 的主体字段。
  - `search_text` 由 seed 写入，用于简单搜索。
  - `display_order` 用于保留 mock 数据默认顺序。
- `video_episodes`
  - `video_id` 外键。
  - `episode`、`title`、`duration`、`status`。
- `video_release_calendar`
  - `video_id` 外键。
  - `item_order`、`time_text`、`detail`、`active`。
- `video_related`
  - `video_id`、`related_video_id` 外键。
  - `display_order`。

本阶段使用 SQL migration 和 seed 文件，不引入 ORM。视频的 `tags`、`cast_names`、频道 `keywords` 使用 `text[]`，方便 Go 映射和查询。

## 代码结构

```text
cmd/api/main.go
internal/config
internal/database
internal/httpapi
internal/video
migrations
seeds
docs/superpowers
DEVELOPMENT.md
README.md
```

职责划分：

- `internal/config`：读取端口、数据库连接串、请求超时等配置。
- `internal/database`：创建 PostgreSQL 连接池，执行健康检查。
- `internal/video`：定义领域模型、查询参数、Repository 和 Service。
- `internal/httpapi`：路由、JSON 响应、错误响应、CORS。
- `cmd/api`：组装配置、数据库、服务和 HTTP server。

## 查询行为

后端第一版复刻前端 mock 查询逻辑：

- 频道匹配：根据频道 `keywords` 匹配视频标题、副标题、分类、徽标、进度和标签。
- `featured`、`all`、空频道：不过滤。
- `type`：模糊匹配视频标题、副标题、分类、徽标、进度和标签。
- `year`、`quality`：精确匹配。
- `search`：关键词匹配标题、副标题、描述、分类、徽标、进度、年份、地区、标签、演员。
- 排序：
  - `hot`：热度数字降序。
  - `score`：评分降序。
  - `new` / `rising`：年份降序，再按热度降序。
  - `reputation`：评分降序，再按热度降序。
  - `vip`：仅返回包含会员、独播或 VIP 语义的视频，再按热度降序。

## 错误处理

- 数据库不可用时 `/health` 返回 `503`。
- API 查询失败返回 `500` 和 `{ "error": "internal server error" }`。
- 未命中频道或视频不返回 `404`，按前端现有 fallback 行为回退到精选频道或主推视频。
- 所有响应设置 `Content-Type: application/json; charset=utf-8`。
- 开发期允许跨域访问，便于 Next.js 前端联调。

## 测试策略

按 TDD 实现：

- `internal/config` 单元测试覆盖默认配置和环境变量覆盖。
- `internal/video` 服务层测试使用内存假 Repository，覆盖首页聚合、排行榜排序、频道筛选、搜索、详情 fallback。
- `internal/httpapi` handler 测试使用 `httptest`，覆盖主要路由 JSON 结构、CORS、错误响应。
- 数据库 migration/seed 通过本地 `go test` 或启动服务后的健康检查验证。

验收标准：

- `go test ./...` 通过。
- 服务能启动并连接本地 `nextvideo`。
- 当前前端配置 `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080` 后，可从真实接口拿到播放页和列表页数据。
- `DEVELOPMENT.md` 记录当前完成、后续计划和下一步。

## 自检

- 没有留下待定占位内容。
- 范围只覆盖视频只读 API，未混入用户系统和后台管理。
- API 路径与前端 `src/lib/video-api.ts` 保持一致。
- 表结构支持当前前端字段和第一版查询行为。
