# Next Video Go Backend

Next Video 的 Go 后端第一版，提供前端当前 API facade 需要的视频只读接口。

## 环境

- Go 1.26+
- PostgreSQL
- 本地数据库：`nextvideo`
- 默认连接串：

```text
postgres://postgres:dengke258567@localhost:5432/nextvideo?sslmode=disable
```

可用环境变量覆盖：

```powershell
$env:PORT="8080"
$env:DATABASE_URL="postgres://postgres:dengke258567@localhost:5432/nextvideo?sslmode=disable"
```

## 初始化数据库

项目内置迁移和 seed 工具，不依赖本机安装 `psql`：

```powershell
go run ./cmd/migrate
```

该命令会依次执行：

- `migrations/*.sql`
- `seeds/*.sql`

seed 可以重复执行，会更新当前演示片库、频道、选集、更新时间线和相关推荐。

## 启动服务

```powershell
go run ./cmd/api
```

默认地址：

```text
http://localhost:8080
```

健康检查：

```text
GET http://localhost:8080/health
```

## 前端接入

在 `next-video` 前端项目中配置：

```text
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

当前支持接口：

- `GET /health`
- `GET /videos/home`
- `GET /videos/rank?sort=&channel=`
- `GET /videos/channel/{slug}?type=&year=&sort=`
- `GET /videos/search?q=&channel=&quality=&type=&year=&sort=`
- `GET /videos/{id}`
- `GET /videos/ids`
- `GET /me`
- `POST /me/register`
- `POST /me/login`
- `POST /me/logout`
- `GET /me/favorites`
- `POST /me/favorites`
- `DELETE /me/favorites/{videoId}`
- `GET /me/watch-history`
- `POST /me/watch-history`
- `DELETE /me/watch-history/{videoId}?episode=1`
- `DELETE /me/watch-history`
- `GET /videos/{id}/comments?sort=latest|hot`
- `POST /videos/{id}/comments`
- `POST /videos/{id}/comments/{commentId}/like`
- `DELETE /videos/{id}/comments/{commentId}`
- `GET /videos/{id}/danmaku`
- `POST /videos/{id}/danmaku`

开发期用户识别：

```text
X-User-ID=demo-user
```

未传 `X-User-ID` 时默认使用 `demo-user`。

真实登录注册 v1：

- `POST /me/register` 请求体：`{"email":"xia@example.com","password":"password123","nickname":"小夏"}`
- `POST /me/login` 请求体：`{"email":"xia@example.com","password":"password123"}`
- 登录/注册成功后，前端会把响应里的 `id` 作为后续 `X-User-ID`，用于收藏、观看历史、VIP、评论和弹幕的数据隔离。

## 验证

```powershell
go test ./...
```
