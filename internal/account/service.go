package account

import (
	"context"
	"errors"
	"strings"
	"time"
)

const DefaultUserID = "demo-user"

var ErrUserNotFound = errors.New("user not found")

type Repository interface {
	GetUser(ctx context.Context, userID string) (UserProfile, error)
	UpsertUser(ctx context.Context, profile UserProfile) error
	ListFavorites(ctx context.Context, userID string) ([]FavoriteItem, error)
	UpsertFavorite(ctx context.Context, userID string, item FavoriteItem) error
	DeleteFavorite(ctx context.Context, userID string, videoID string) error
	ListWatchHistory(ctx context.Context, userID string) ([]WatchHistoryItem, error)
	UpsertWatchHistory(ctx context.Context, userID string, item WatchHistoryItem) error
	DeleteWatchHistory(ctx context.Context, userID string, videoID string, episode *int) error
	ClearWatchHistory(ctx context.Context, userID string) error
}

type Service struct {
	repository Repository
	now        func() time.Time
}

// NewService 创建账号状态服务；repository 负责持久化，now 提供当前时间。
func NewService(repository Repository, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}

	return &Service{
		repository: repository,
		now:        now,
	}
}

// Profile 获取用户资料；ctx 为请求上下文，userID 为用户标识。
func (service *Service) Profile(ctx context.Context, userID string) (UserProfile, error) {
	normalizedUserID := normalizeUserID(userID)
	profile, err := service.repository.GetUser(ctx, normalizedUserID)
	if err == nil {
		return profile, nil
	}
	if !errors.Is(err, ErrUserNotFound) {
		return UserProfile{}, err
	}

	profile = defaultProfile(normalizedUserID)
	if err := service.repository.UpsertUser(ctx, profile); err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}

// Login 写入开发态登录资料；ctx 为请求上下文，userID 为用户标识，input 为登录输入。
func (service *Service) Login(ctx context.Context, userID string, input LoginInput) (UserProfile, error) {
	profile := defaultProfile(normalizeUserID(userID))
	contact := strings.TrimSpace(input.Contact)

	profile.AvatarURL = strings.TrimSpace(input.AvatarURL)
	profile.IsLoggedIn = true
	profile.Nickname = strings.TrimSpace(input.Nickname)
	if profile.Nickname == "" {
		profile.Nickname = defaultNickname
	}
	if strings.Contains(contact, "@") {
		profile.Email = contact
	} else {
		profile.Phone = contact
	}

	if err := service.repository.UpsertUser(ctx, profile); err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}

// ActivateVIP 开通当前用户 VIP；ctx 为请求上下文，userID 为用户标识，input 为 VIP 到期日。
func (service *Service) ActivateVIP(ctx context.Context, userID string, input VipInput) (UserProfile, error) {
	normalizedUserID := normalizeUserID(userID)
	profile, err := service.repository.GetUser(ctx, normalizedUserID)
	if errors.Is(err, ErrUserNotFound) {
		profile = defaultProfile(normalizedUserID)
	} else if err != nil {
		return UserProfile{}, err
	}

	profile.IsLoggedIn = true
	profile.IsVip = true
	profile.VipUntil = strings.TrimSpace(input.VipUntil)

	if err := service.repository.UpsertUser(ctx, profile); err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}

// Logout 清空开发态登录资料；ctx 为请求上下文，userID 为用户标识。
func (service *Service) Logout(ctx context.Context, userID string) (UserProfile, error) {
	profile := defaultProfile(normalizeUserID(userID))
	if err := service.repository.UpsertUser(ctx, profile); err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}

// Favorites 获取收藏列表；ctx 为请求上下文，userID 为用户标识。
func (service *Service) Favorites(ctx context.Context, userID string) ([]FavoriteItem, error) {
	items, err := service.repository.ListFavorites(ctx, normalizeUserID(userID))
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []FavoriteItem{}, nil
	}

	return items, nil
}

// AddFavorite 新增或更新收藏；ctx 为请求上下文，userID 为用户标识，input 为收藏摘要。
func (service *Service) AddFavorite(ctx context.Context, userID string, input FavoriteInput) (FavoriteItem, error) {
	item := FavoriteItem{
		ID:            strings.TrimSpace(input.ID),
		Title:         strings.TrimSpace(input.Title),
		Category:      strings.TrimSpace(input.Category),
		Progress:      strings.TrimSpace(input.Progress),
		CoverGradient: strings.TrimSpace(input.CoverGradient),
		Description:   strings.TrimSpace(input.Description),
		AddedAt:       service.now().UnixMilli(),
	}
	if err := service.repository.UpsertFavorite(ctx, normalizeUserID(userID), item); err != nil {
		return FavoriteItem{}, err
	}

	return item, nil
}

// DeleteFavorite 删除收藏；ctx 为请求上下文，userID 和 videoID 定位收藏。
func (service *Service) DeleteFavorite(ctx context.Context, userID string, videoID string) error {
	return service.repository.DeleteFavorite(ctx, normalizeUserID(userID), strings.TrimSpace(videoID))
}

// WatchHistory 获取观看历史；ctx 为请求上下文，userID 为用户标识。
func (service *Service) WatchHistory(ctx context.Context, userID string) ([]WatchHistoryItem, error) {
	items, err := service.repository.ListWatchHistory(ctx, normalizeUserID(userID))
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []WatchHistoryItem{}, nil
	}

	return items, nil
}

// AddWatchHistory 新增或更新观看历史；ctx 为请求上下文，userID 为用户标识，input 为历史摘要。
func (service *Service) AddWatchHistory(ctx context.Context, userID string, input WatchHistoryInput) (WatchHistoryItem, error) {
	item := WatchHistoryItem{
		ID:              strings.TrimSpace(input.ID),
		Title:           strings.TrimSpace(input.Title),
		Category:        strings.TrimSpace(input.Category),
		Progress:        strings.TrimSpace(input.Progress),
		CoverGradient:   strings.TrimSpace(input.CoverGradient),
		Episode:         input.Episode,
		WatchSeconds:    input.WatchSeconds,
		DurationSeconds: input.DurationSeconds,
		WatchedAt:       service.now().UnixMilli(),
	}
	if err := service.repository.UpsertWatchHistory(ctx, normalizeUserID(userID), item); err != nil {
		return WatchHistoryItem{}, err
	}

	return item, nil
}

// DeleteWatchHistory 删除观看历史；ctx 为请求上下文，userID、videoID 和 episode 定位历史。
func (service *Service) DeleteWatchHistory(ctx context.Context, userID string, videoID string, episode *int) error {
	return service.repository.DeleteWatchHistory(ctx, normalizeUserID(userID), strings.TrimSpace(videoID), episode)
}

// ClearWatchHistory 清空观看历史；ctx 为请求上下文，userID 为用户标识。
func (service *Service) ClearWatchHistory(ctx context.Context, userID string) error {
	return service.repository.ClearWatchHistory(ctx, normalizeUserID(userID))
}

const defaultNickname = "Next Video 用户"

// defaultProfile 生成默认用户资料；userID 为用户标识。
func defaultProfile(userID string) UserProfile {
	return UserProfile{
		ID:       userID,
		Nickname: defaultNickname,
	}
}

// normalizeUserID 规范化用户标识；userID 为空时返回开发态默认用户。
func normalizeUserID(userID string) string {
	normalizedUserID := strings.TrimSpace(userID)
	if normalizedUserID == "" {
		return DefaultUserID
	}

	return normalizedUserID
}
