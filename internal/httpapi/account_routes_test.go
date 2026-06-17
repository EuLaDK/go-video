package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"next-video-golang/internal/account"
	"next-video-golang/internal/httpapi"
)

type fakeAccountService struct {
	userID       string
	profile      account.UserProfile
	favorites    []account.FavoriteItem
	watchHistory []account.WatchHistoryItem
}

// Profile 模拟用户资料；ctx 为请求上下文，userID 为用户标识。
func (service *fakeAccountService) Profile(ctx context.Context, userID string) (account.UserProfile, error) {
	service.userID = userID
	if service.profile.ID == "" {
		service.profile = account.UserProfile{ID: userID, Nickname: "Next Video 用户"}
	}
	return service.profile, nil
}

// Login 模拟登录；ctx 为请求上下文，userID 为用户标识，input 为登录输入。
func (service *fakeAccountService) Login(ctx context.Context, userID string, input account.LoginInput) (account.UserProfile, error) {
	service.userID = userID
	service.profile = account.UserProfile{ID: userID, IsLoggedIn: true, Nickname: input.Nickname, Email: input.Contact}
	return service.profile, nil
}

// Logout 模拟退出；ctx 为请求上下文，userID 为用户标识。
func (service *fakeAccountService) Logout(ctx context.Context, userID string) (account.UserProfile, error) {
	service.userID = userID
	service.profile = account.UserProfile{ID: userID, Nickname: "Next Video 用户"}
	return service.profile, nil
}

// Favorites 模拟收藏列表；ctx 为请求上下文，userID 为用户标识。
func (service *fakeAccountService) Favorites(ctx context.Context, userID string) ([]account.FavoriteItem, error) {
	service.userID = userID
	return service.favorites, nil
}

// AddFavorite 模拟新增收藏；ctx 为请求上下文，userID 为用户标识，input 为收藏输入。
func (service *fakeAccountService) AddFavorite(ctx context.Context, userID string, input account.FavoriteInput) (account.FavoriteItem, error) {
	service.userID = userID
	item := account.FavoriteItem{
		ID:            input.ID,
		Title:         input.Title,
		Category:      input.Category,
		Progress:      input.Progress,
		CoverGradient: input.CoverGradient,
		Description:   input.Description,
		AddedAt:       123,
	}
	service.favorites = []account.FavoriteItem{item}
	return item, nil
}

// DeleteFavorite 模拟删除收藏；ctx 为请求上下文，userID 和 videoID 定位收藏。
func (service *fakeAccountService) DeleteFavorite(ctx context.Context, userID string, videoID string) error {
	service.userID = userID
	service.favorites = []account.FavoriteItem{}
	return nil
}

// WatchHistory 模拟观看历史；ctx 为请求上下文，userID 为用户标识。
func (service *fakeAccountService) WatchHistory(ctx context.Context, userID string) ([]account.WatchHistoryItem, error) {
	service.userID = userID
	return service.watchHistory, nil
}

// AddWatchHistory 模拟新增历史；ctx 为请求上下文，userID 为用户标识，input 为历史输入。
func (service *fakeAccountService) AddWatchHistory(ctx context.Context, userID string, input account.WatchHistoryInput) (account.WatchHistoryItem, error) {
	service.userID = userID
	item := account.WatchHistoryItem{
		ID:              input.ID,
		Title:           input.Title,
		Category:        input.Category,
		Progress:        input.Progress,
		CoverGradient:   input.CoverGradient,
		Episode:         input.Episode,
		WatchSeconds:    input.WatchSeconds,
		DurationSeconds: input.DurationSeconds,
		WatchedAt:       456,
	}
	service.watchHistory = []account.WatchHistoryItem{item}
	return item, nil
}

// DeleteWatchHistory 模拟删除历史；ctx 为请求上下文，userID、videoID 和 episode 定位历史。
func (service *fakeAccountService) DeleteWatchHistory(ctx context.Context, userID string, videoID string, episode *int) error {
	service.userID = userID
	service.watchHistory = []account.WatchHistoryItem{}
	return nil
}

// ClearWatchHistory 模拟清空历史；ctx 为请求上下文，userID 为用户标识。
func (service *fakeAccountService) ClearWatchHistory(ctx context.Context, userID string) error {
	service.userID = userID
	service.watchHistory = []account.WatchHistoryItem{}
	return nil
}

func TestAccountRoutesUseDefaultUser(t *testing.T) {
	accountService := &fakeAccountService{}
	srv := httpapi.NewServer(&fakeVideoService{}, accountService)
	rec := requestJSON(srv, http.MethodGet, "/me", "", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if accountService.userID != account.DefaultUserID {
		t.Fatalf("userID = %q, want %q", accountService.userID, account.DefaultUserID)
	}
}

func TestAccountRoutesLoginAndLogout(t *testing.T) {
	accountService := &fakeAccountService{}
	srv := httpapi.NewServer(&fakeVideoService{}, accountService)

	rec := requestJSON(srv, http.MethodPost, "/me/login", `{"nickname":"小夏","contact":"xia@example.com"}`, "custom-user")
	if rec.Code != http.StatusOK {
		t.Fatalf("login status = %d, want 200", rec.Code)
	}

	var profile account.UserProfile
	if err := json.Unmarshal(rec.Body.Bytes(), &profile); err != nil {
		t.Fatal(err)
	}
	if !profile.IsLoggedIn || profile.Nickname != "小夏" || accountService.userID != "custom-user" {
		t.Fatalf("profile = %#v userID = %q", profile, accountService.userID)
	}

	rec = requestJSON(srv, http.MethodPost, "/me/logout", "", "custom-user")
	if rec.Code != http.StatusOK {
		t.Fatalf("logout status = %d, want 200", rec.Code)
	}
}

func TestAccountRoutesFavorites(t *testing.T) {
	accountService := &fakeAccountService{}
	srv := httpapi.NewServer(&fakeVideoService{}, accountService)

	rec := requestJSON(srv, http.MethodPost, "/me/favorites", `{"id":"xinghe","title":"星河回响","category":"科幻 / 悬疑","progress":"会员抢先看","coverGradient":"gradient","description":"深空信号"}`, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("add favorite status = %d, want 200", rec.Code)
	}

	rec = requestJSON(srv, http.MethodGet, "/me/favorites", "", "")
	var items []account.FavoriteItem
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != "xinghe" {
		t.Fatalf("items = %#v", items)
	}

	rec = requestJSON(srv, http.MethodDelete, "/me/favorites/xinghe", "", "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete favorite status = %d, want 204", rec.Code)
	}
}

func TestAccountRoutesWatchHistory(t *testing.T) {
	accountService := &fakeAccountService{}
	srv := httpapi.NewServer(&fakeVideoService{}, accountService)

	rec := requestJSON(srv, http.MethodPost, "/me/watch-history", `{"id":"xinghe","title":"星河回响","category":"科幻 / 悬疑","progress":"会员抢先看","coverGradient":"gradient","episode":2,"watchSeconds":90,"durationSeconds":2700}`, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("add history status = %d, want 200", rec.Code)
	}

	rec = requestJSON(srv, http.MethodGet, "/me/watch-history", "", "")
	var items []account.WatchHistoryItem
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != "xinghe" || items[0].Episode == nil || *items[0].Episode != 2 {
		t.Fatalf("items = %#v", items)
	}

	rec = requestJSON(srv, http.MethodDelete, "/me/watch-history/xinghe?episode=2", "", "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete history status = %d, want 204", rec.Code)
	}

	accountService.watchHistory = []account.WatchHistoryItem{{ID: "anye"}}
	rec = requestJSON(srv, http.MethodDelete, "/me/watch-history", "", "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("clear history status = %d, want 204", rec.Code)
	}
	if !reflect.DeepEqual(accountService.watchHistory, []account.WatchHistoryItem{}) {
		t.Fatalf("watchHistory = %#v", accountService.watchHistory)
	}
}

// requestJSON 发送测试 JSON 请求；handler 为被测服务，method/path/body/userID 描述请求。
func requestJSON(handler http.Handler, method string, path string, body string, userID string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if userID != "" {
		req.Header.Set("X-User-ID", userID)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}
