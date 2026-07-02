package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"next-video-golang/internal/account"
)

type AccountService interface {
	Profile(ctx context.Context, userID string) (account.UserProfile, error)
	Register(ctx context.Context, input account.RegisterInput) (account.UserProfile, error)
	Login(ctx context.Context, userID string, input account.LoginInput) (account.UserProfile, error)
	ActivateVIP(ctx context.Context, userID string, input account.VipInput) (account.UserProfile, error)
	Logout(ctx context.Context, userID string) (account.UserProfile, error)
	Favorites(ctx context.Context, userID string) ([]account.FavoriteItem, error)
	AddFavorite(ctx context.Context, userID string, input account.FavoriteInput) (account.FavoriteItem, error)
	DeleteFavorite(ctx context.Context, userID string, videoID string) error
	WatchHistory(ctx context.Context, userID string) ([]account.WatchHistoryItem, error)
	AddWatchHistory(ctx context.Context, userID string, input account.WatchHistoryInput) (account.WatchHistoryItem, error)
	DeleteWatchHistory(ctx context.Context, userID string, videoID string, episode *int) error
	ClearWatchHistory(ctx context.Context, userID string) error
}

// isAccountPath 判断请求是否属于账号状态路由；path 为请求路径。
func isAccountPath(path string) bool {
	return path == "/me" || strings.HasPrefix(path, "/me/")
}

// handleAccount 分发账号状态请求；response 为响应写入器，request 为当前请求。
func (server *Server) handleAccount(response http.ResponseWriter, request *http.Request) {
	if server.accountService == nil {
		writeError(response, http.StatusNotFound, "not found")
		return
	}

	switch {
	case request.URL.Path == "/me":
		server.handleProfile(response, request)
	case request.URL.Path == "/me/register":
		server.handleRegister(response, request)
	case request.URL.Path == "/me/login":
		server.handleLogin(response, request)
	case request.URL.Path == "/me/vip":
		server.handleActivateVIP(response, request)
	case request.URL.Path == "/me/logout":
		server.handleLogout(response, request)
	case request.URL.Path == "/me/favorites":
		server.handleFavorites(response, request)
	case strings.HasPrefix(request.URL.Path, "/me/favorites/"):
		server.handleFavoriteItem(response, request)
	case request.URL.Path == "/me/watch-history":
		server.handleWatchHistory(response, request)
	case strings.HasPrefix(request.URL.Path, "/me/watch-history/"):
		server.handleWatchHistoryItem(response, request)
	default:
		writeError(response, http.StatusNotFound, "not found")
	}
}

// handleProfile 返回当前用户资料；response 为响应写入器，request 为当前请求。
func (server *Server) handleProfile(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	profile, err := server.accountService.Profile(request.Context(), userIDFromRequest(request))
	if err != nil {
		writeError(response, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(response, http.StatusOK, profile)
}

// handleRegister 注册邮箱密码账号；response 为响应写入器，request 为当前请求。
func (server *Server) handleRegister(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input account.RegisterInput
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(response, http.StatusBadRequest, "invalid json body")
		return
	}

	profile, err := server.accountService.Register(request.Context(), input)
	if err != nil {
		writeError(response, accountErrorStatus(err), accountErrorMessage(err))
		return
	}

	writeJSON(response, http.StatusCreated, profile)
}

// handleLogin 写入开发态登录资料；response 为响应写入器，request 为当前请求。
func (server *Server) handleLogin(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input account.LoginInput
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(response, http.StatusBadRequest, "invalid json body")
		return
	}

	profile, err := server.accountService.Login(request.Context(), userIDFromRequest(request), input)
	if err != nil {
		writeError(response, accountErrorStatus(err), accountErrorMessage(err))
		return
	}

	writeJSON(response, http.StatusOK, profile)
}

// accountErrorStatus 将账号服务错误映射为 HTTP 状态码；err 为服务层错误。
func accountErrorStatus(err error) int {
	switch {
	case errors.Is(err, account.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, account.ErrEmailAlreadyRegistered):
		return http.StatusConflict
	case errors.Is(err, account.ErrInvalidAuthInput):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// accountErrorMessage 将账号服务错误映射为安全响应消息；err 为服务层错误。
func accountErrorMessage(err error) string {
	switch {
	case errors.Is(err, account.ErrInvalidCredentials):
		return "invalid credentials"
	case errors.Is(err, account.ErrEmailAlreadyRegistered):
		return "email already registered"
	case errors.Is(err, account.ErrInvalidAuthInput):
		return "invalid auth input"
	default:
		return "internal server error"
	}
}

// handleActivateVIP 写入当前用户 VIP 状态；response 为响应写入器，request 为当前请求。
func (server *Server) handleActivateVIP(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input account.VipInput
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(response, http.StatusBadRequest, "invalid json body")
		return
	}

	profile, err := server.accountService.ActivateVIP(request.Context(), userIDFromRequest(request), input)
	if err != nil {
		writeError(response, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(response, http.StatusOK, profile)
}

// handleLogout 清空开发态登录资料；response 为响应写入器，request 为当前请求。
func (server *Server) handleLogout(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	profile, err := server.accountService.Logout(request.Context(), userIDFromRequest(request))
	if err != nil {
		writeError(response, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(response, http.StatusOK, profile)
}

// handleFavorites 读取或新增收藏；response 为响应写入器，request 为当前请求。
func (server *Server) handleFavorites(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		items, err := server.accountService.Favorites(request.Context(), userIDFromRequest(request))
		if err != nil {
			writeError(response, http.StatusInternalServerError, "internal server error")
			return
		}
		writeJSON(response, http.StatusOK, items)
	case http.MethodPost:
		var input account.FavoriteInput
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			writeError(response, http.StatusBadRequest, "invalid json body")
			return
		}
		item, err := server.accountService.AddFavorite(request.Context(), userIDFromRequest(request), input)
		if err != nil {
			writeError(response, http.StatusInternalServerError, "internal server error")
			return
		}
		writeJSON(response, http.StatusOK, item)
	default:
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleFavoriteItem 删除单条收藏；response 为响应写入器，request 为当前请求。
func (server *Server) handleFavoriteItem(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	videoID := pathValue(strings.TrimPrefix(request.URL.Path, "/me/favorites/"))
	if err := server.accountService.DeleteFavorite(request.Context(), userIDFromRequest(request), videoID); err != nil {
		writeError(response, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

// handleWatchHistory 读取、新增或清空观看历史；response 为响应写入器，request 为当前请求。
func (server *Server) handleWatchHistory(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		items, err := server.accountService.WatchHistory(request.Context(), userIDFromRequest(request))
		if err != nil {
			writeError(response, http.StatusInternalServerError, "internal server error")
			return
		}
		writeJSON(response, http.StatusOK, items)
	case http.MethodPost:
		var input account.WatchHistoryInput
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			writeError(response, http.StatusBadRequest, "invalid json body")
			return
		}
		item, err := server.accountService.AddWatchHistory(request.Context(), userIDFromRequest(request), input)
		if err != nil {
			writeError(response, http.StatusInternalServerError, "internal server error")
			return
		}
		writeJSON(response, http.StatusOK, item)
	case http.MethodDelete:
		if err := server.accountService.ClearWatchHistory(request.Context(), userIDFromRequest(request)); err != nil {
			writeError(response, http.StatusInternalServerError, "internal server error")
			return
		}
		response.WriteHeader(http.StatusNoContent)
	default:
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleWatchHistoryItem 删除单条观看历史；response 为响应写入器，request 为当前请求。
func (server *Server) handleWatchHistoryItem(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	episode, ok := episodeFromRequest(response, request)
	if !ok {
		return
	}

	videoID := pathValue(strings.TrimPrefix(request.URL.Path, "/me/watch-history/"))
	if err := server.accountService.DeleteWatchHistory(request.Context(), userIDFromRequest(request), videoID, episode); err != nil {
		writeError(response, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

// episodeFromRequest 读取可选集数参数；response 为错误写入器，request 为当前请求。
func episodeFromRequest(response http.ResponseWriter, request *http.Request) (*int, bool) {
	rawEpisode := strings.TrimSpace(request.URL.Query().Get("episode"))
	if rawEpisode == "" {
		return nil, true
	}

	episode, err := strconv.Atoi(rawEpisode)
	if err != nil || episode <= 0 {
		writeError(response, http.StatusBadRequest, "invalid episode")
		return nil, false
	}

	return &episode, true
}
