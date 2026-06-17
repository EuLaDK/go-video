package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"next-video-golang/internal/account"
	"next-video-golang/internal/interaction"
	"next-video-golang/internal/video"
)

type VideoService interface {
	HealthCheck(ctx context.Context) error
	HomePage(ctx context.Context) (video.HomePageData, error)
	RankVideos(ctx context.Context, query video.RankQuery) ([]video.Video, error)
	ChannelPage(ctx context.Context, query video.ChannelQuery) (video.ChannelPageData, error)
	SearchPage(ctx context.Context, query video.SearchQuery) (video.SearchPageData, error)
	WatchPage(ctx context.Context, videoID string) (video.WatchPageData, error)
	VideoIDs(ctx context.Context) ([]string, error)
}

type Server struct {
	videoService       VideoService
	accountService     AccountService
	interactionService InteractionService
}

// NewServer 创建 HTTP 服务；videoService 提供视频查询能力，accountServices 可选提供用户状态能力。
func NewServer(videoService VideoService, accountServices ...AccountService) *Server {
	var accountService AccountService
	if len(accountServices) > 0 {
		accountService = accountServices[0]
	}

	return NewServerWithServices(videoService, accountService, nil)
}

// NewServerWithServices 创建 HTTP 服务；videoService 提供视频查询，accountService 和 interactionService 提供可选状态能力。
func NewServerWithServices(videoService VideoService, accountService AccountService, interactionService InteractionService) *Server {
	return &Server{
		videoService:       videoService,
		accountService:     accountService,
		interactionService: interactionService,
	}
}

// ServeHTTP 分发 HTTP 请求；response 为响应写入器，request 为当前请求。
func (server *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	addCORSHeaders(response)

	if request.Method == http.MethodOptions {
		response.WriteHeader(http.StatusNoContent)
		return
	}

	if isAccountPath(request.URL.Path) {
		server.handleAccount(response, request)
		return
	}

	if isInteractionPath(request.URL.Path) {
		server.handleInteraction(response, request)
		return
	}

	if request.Method != http.MethodGet {
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	switch {
	case request.URL.Path == "/health":
		server.handleHealth(response, request)
	case request.URL.Path == "/videos/home":
		server.handleHome(response, request)
	case request.URL.Path == "/videos/rank":
		server.handleRank(response, request)
	case request.URL.Path == "/videos/search":
		server.handleSearch(response, request)
	case request.URL.Path == "/videos/ids":
		server.handleVideoIDs(response, request)
	case strings.HasPrefix(request.URL.Path, "/videos/channel/"):
		server.handleChannel(response, request)
	case strings.HasPrefix(request.URL.Path, "/videos/"):
		server.handleWatch(response, request)
	default:
		writeError(response, http.StatusNotFound, "not found")
	}
}

// handleHealth 返回健康检查结果；response 为响应写入器，request 为当前请求。
func (server *Server) handleHealth(response http.ResponseWriter, request *http.Request) {
	if err := server.videoService.HealthCheck(request.Context()); err != nil {
		writeJSON(response, http.StatusServiceUnavailable, map[string]string{
			"status": "unhealthy",
		})
		return
	}

	writeJSON(response, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// handleHome 返回首页数据；response 为响应写入器，request 为当前请求。
func (server *Server) handleHome(response http.ResponseWriter, request *http.Request) {
	data, err := server.videoService.HomePage(request.Context())
	if err != nil {
		writeError(response, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(response, http.StatusOK, data)
}

// handleRank 返回排行榜数据；response 为响应写入器，request 为当前请求。
func (server *Server) handleRank(response http.ResponseWriter, request *http.Request) {
	query := video.RankQuery{
		Channel: request.URL.Query().Get("channel"),
		Sort:    request.URL.Query().Get("sort"),
	}
	data, err := server.videoService.RankVideos(request.Context(), query)
	if err != nil {
		writeError(response, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(response, http.StatusOK, data)
}

// handleChannel 返回频道页数据；response 为响应写入器，request 为当前请求。
func (server *Server) handleChannel(response http.ResponseWriter, request *http.Request) {
	slug := strings.TrimPrefix(request.URL.Path, "/videos/channel/")
	slug = pathValue(slug)
	query := video.ChannelQuery{
		Slug: slug,
		Type: request.URL.Query().Get("type"),
		Year: request.URL.Query().Get("year"),
		Sort: request.URL.Query().Get("sort"),
	}
	data, err := server.videoService.ChannelPage(request.Context(), query)
	if err != nil {
		writeError(response, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(response, http.StatusOK, data)
}

// handleSearch 返回搜索页数据；response 为响应写入器，request 为当前请求。
func (server *Server) handleSearch(response http.ResponseWriter, request *http.Request) {
	query := video.SearchQuery{
		Q:       request.URL.Query().Get("q"),
		Channel: request.URL.Query().Get("channel"),
		Quality: request.URL.Query().Get("quality"),
		Type:    request.URL.Query().Get("type"),
		Year:    request.URL.Query().Get("year"),
		Sort:    request.URL.Query().Get("sort"),
	}
	data, err := server.videoService.SearchPage(request.Context(), query)
	if err != nil {
		writeError(response, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(response, http.StatusOK, data)
}

// handleWatch 返回播放页数据；response 为响应写入器，request 为当前请求。
func (server *Server) handleWatch(response http.ResponseWriter, request *http.Request) {
	videoID := strings.TrimPrefix(request.URL.Path, "/videos/")
	videoID = pathValue(videoID)
	data, err := server.videoService.WatchPage(request.Context(), videoID)
	if err != nil {
		writeError(response, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(response, http.StatusOK, data)
}

// handleVideoIDs 返回视频 id 列表；response 为响应写入器，request 为当前请求。
func (server *Server) handleVideoIDs(response http.ResponseWriter, request *http.Request) {
	data, err := server.videoService.VideoIDs(request.Context())
	if err != nil {
		writeError(response, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(response, http.StatusOK, data)
}

// addCORSHeaders 添加开发期跨域响应头；response 为响应写入器。
func addCORSHeaders(response http.ResponseWriter) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// writeJSON 写入 JSON 响应；response 为响应写入器，status 为状态码，data 为响应体。
func writeJSON(response http.ResponseWriter, status int, data any) {
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	response.WriteHeader(status)
	if err := json.NewEncoder(response).Encode(data); err != nil {
		return
	}
}

// writeError 写入错误 JSON；response 为响应写入器，status 为状态码，message 为错误消息。
func writeError(response http.ResponseWriter, status int, message string) {
	writeJSON(response, status, map[string]string{
		"error": message,
	})
}

// pathValue 解码路径参数；value 为 URL 路径中的单段参数。
func pathValue(value string) string {
	decoded, err := url.PathUnescape(value)
	if err != nil {
		return value
	}

	return decoded
}

// userIDFromRequest 读取开发态用户标识；request 为当前请求，未传时返回默认用户。
func userIDFromRequest(request *http.Request) string {
	userID := strings.TrimSpace(request.Header.Get("X-User-ID"))
	if userID == "" {
		return account.DefaultUserID
	}

	return userID
}

// commentSortFromRequest 读取评论排序参数；request 为当前请求。
func commentSortFromRequest(request *http.Request) interaction.CommentSort {
	if request.URL.Query().Get("sort") == string(interaction.CommentSortHot) {
		return interaction.CommentSortHot
	}

	return interaction.CommentSortLatest
}
