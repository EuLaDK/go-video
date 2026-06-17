package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"next-video-golang/internal/interaction"
)

type InteractionService interface {
	Comments(ctx context.Context, userID string, videoID string, sort interaction.CommentSort) ([]interaction.CommentItem, error)
	AddComment(ctx context.Context, userID string, videoID string, input interaction.CommentInput) (interaction.CommentItem, error)
	ToggleCommentLike(ctx context.Context, userID string, videoID string, commentID string) (interaction.CommentItem, error)
	DeleteComment(ctx context.Context, userID string, videoID string, commentID string) error
	Danmaku(ctx context.Context, videoID string) ([]interaction.DanmakuItem, error)
	AddDanmaku(ctx context.Context, userID string, videoID string, input interaction.DanmakuInput) (interaction.DanmakuItem, error)
}

type interactionPath struct {
	videoID  string
	resource string
	itemID   string
	action   string
}

// isInteractionPath 判断请求是否属于视频互动路由；path 为请求路径。
func isInteractionPath(path string) bool {
	_, ok := parseInteractionPath(path)
	return ok
}

// handleInteraction 分发视频互动请求；response 为响应写入器，request 为当前请求。
func (server *Server) handleInteraction(response http.ResponseWriter, request *http.Request) {
	if server.interactionService == nil {
		writeError(response, http.StatusNotFound, "not found")
		return
	}

	parsedPath, ok := parseInteractionPath(request.URL.Path)
	if !ok {
		writeError(response, http.StatusNotFound, "not found")
		return
	}

	switch parsedPath.resource {
	case "comments":
		server.handleCommentRoutes(response, request, parsedPath)
	case "danmaku":
		server.handleDanmakuRoutes(response, request, parsedPath)
	default:
		writeError(response, http.StatusNotFound, "not found")
	}
}

// handleCommentRoutes 分发评论请求；response 为响应写入器，request 为当前请求，path 为解析后的路由参数。
func (server *Server) handleCommentRoutes(response http.ResponseWriter, request *http.Request, path interactionPath) {
	switch {
	case path.itemID == "" && path.action == "" && request.Method == http.MethodGet:
		items, err := server.interactionService.Comments(request.Context(), userIDFromRequest(request), path.videoID, commentSortFromRequest(request))
		if err != nil {
			writeError(response, http.StatusInternalServerError, "internal server error")
			return
		}
		writeJSON(response, http.StatusOK, items)
	case path.itemID == "" && path.action == "" && request.Method == http.MethodPost:
		var input interaction.CommentInput
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			writeError(response, http.StatusBadRequest, "invalid json body")
			return
		}
		item, err := server.interactionService.AddComment(request.Context(), userIDFromRequest(request), path.videoID, input)
		writeInteractionItem(response, item, err)
	case path.itemID != "" && path.action == "like" && request.Method == http.MethodPost:
		item, err := server.interactionService.ToggleCommentLike(request.Context(), userIDFromRequest(request), path.videoID, path.itemID)
		writeInteractionItem(response, item, err)
	case path.itemID != "" && path.action == "" && request.Method == http.MethodDelete:
		err := server.interactionService.DeleteComment(request.Context(), userIDFromRequest(request), path.videoID, path.itemID)
		if writeInteractionError(response, err) {
			return
		}
		response.WriteHeader(http.StatusNoContent)
	default:
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleDanmakuRoutes 分发弹幕请求；response 为响应写入器，request 为当前请求，path 为解析后的路由参数。
func (server *Server) handleDanmakuRoutes(response http.ResponseWriter, request *http.Request, path interactionPath) {
	if path.itemID != "" || path.action != "" {
		writeError(response, http.StatusNotFound, "not found")
		return
	}

	switch request.Method {
	case http.MethodGet:
		items, err := server.interactionService.Danmaku(request.Context(), path.videoID)
		if err != nil {
			writeError(response, http.StatusInternalServerError, "internal server error")
			return
		}
		writeJSON(response, http.StatusOK, items)
	case http.MethodPost:
		var input interaction.DanmakuInput
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			writeError(response, http.StatusBadRequest, "invalid json body")
			return
		}
		item, err := server.interactionService.AddDanmaku(request.Context(), userIDFromRequest(request), path.videoID, input)
		if writeInteractionError(response, err) {
			return
		}
		writeJSON(response, http.StatusOK, item)
	default:
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// writeInteractionItem 写入评论结果；response 为响应写入器，item 为评论数据，err 为服务错误。
func writeInteractionItem(response http.ResponseWriter, item interaction.CommentItem, err error) {
	if writeInteractionError(response, err) {
		return
	}

	writeJSON(response, http.StatusOK, item)
}

// writeInteractionError 写入互动错误；response 为响应写入器，err 为服务错误，返回值表示是否已写响应。
func writeInteractionError(response http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, interaction.ErrInvalidContent) {
		writeError(response, http.StatusBadRequest, "invalid content")
		return true
	}
	if errors.Is(err, interaction.ErrCommentNotFound) {
		writeError(response, http.StatusNotFound, "not found")
		return true
	}

	writeError(response, http.StatusInternalServerError, "internal server error")
	return true
}

// parseInteractionPath 解析视频互动路径；path 为请求路径。
func parseInteractionPath(path string) (interactionPath, bool) {
	rest := strings.TrimPrefix(path, "/videos/")
	if rest == path {
		return interactionPath{}, false
	}

	parts := strings.Split(rest, "/")
	if len(parts) < 2 || len(parts) > 4 {
		return interactionPath{}, false
	}
	if parts[0] == "" || parts[1] == "" {
		return interactionPath{}, false
	}
	if parts[1] != "comments" && parts[1] != "danmaku" {
		return interactionPath{}, false
	}
	if len(parts) == 3 && parts[2] == "" {
		return interactionPath{}, false
	}
	if len(parts) == 4 && (parts[2] == "" || parts[3] == "") {
		return interactionPath{}, false
	}

	parsedPath := interactionPath{
		videoID:  pathValue(parts[0]),
		resource: parts[1],
	}
	if len(parts) >= 3 {
		parsedPath.itemID = pathValue(parts[2])
	}
	if len(parts) == 4 {
		parsedPath.action = parts[3]
	}

	return parsedPath, true
}
