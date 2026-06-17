package interaction

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

const (
	defaultAuthor = "我"
	defaultColor  = "white"
	defaultUserID = "demo-user"
)

var (
	ErrInvalidContent  = errors.New("invalid content")
	ErrCommentNotFound = errors.New("comment not found")
	defaultIDCounter   atomic.Int64
)

type Repository interface {
	ListComments(ctx context.Context, videoID string, userID string, sort CommentSort) ([]CommentItem, error)
	InsertComment(ctx context.Context, userID string, item CommentItem) error
	ToggleCommentLike(ctx context.Context, videoID string, commentID string, userID string) (CommentItem, error)
	DeleteComment(ctx context.Context, videoID string, commentID string, userID string) error
	ListDanmaku(ctx context.Context, videoID string) ([]DanmakuItem, error)
	InsertDanmaku(ctx context.Context, userID string, item DanmakuItem) error
}

type Service struct {
	repository Repository
	now        func() time.Time
	newID      func(prefix string) string
}

// NewService 创建视频互动服务；repository 负责持久化，now/newID 用于生成可测试的时间和 id。
func NewService(repository Repository, now func() time.Time, newID func(prefix string) string) *Service {
	if now == nil {
		now = time.Now
	}
	if newID == nil {
		newID = nextID
	}

	return &Service{
		repository: repository,
		now:        now,
		newID:      newID,
	}
}

// Comments 获取视频评论；ctx 为请求上下文，userID/videoID/sort 分别用于当前用户、视频和排序。
func (service *Service) Comments(ctx context.Context, userID string, videoID string, sort CommentSort) ([]CommentItem, error) {
	items, err := service.repository.ListComments(ctx, normalizeVideoID(videoID), normalizeUserID(userID), normalizeSort(sort))
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []CommentItem{}, nil
	}

	return items, nil
}

// AddComment 新增评论；ctx 为请求上下文，userID/videoID 定位用户和视频，input 为评论内容。
func (service *Service) AddComment(ctx context.Context, userID string, videoID string, input CommentInput) (CommentItem, error) {
	content := strings.TrimSpace(input.Content)
	if content == "" {
		return CommentItem{}, ErrInvalidContent
	}

	item := CommentItem{
		ID:        service.newID("comment"),
		VideoID:   normalizeVideoID(videoID),
		Content:   content,
		Author:    defaultAuthor,
		LikedByMe: false,
		Likes:     0,
		CreatedAt: service.now().UnixMilli(),
	}
	if err := service.repository.InsertComment(ctx, normalizeUserID(userID), item); err != nil {
		return CommentItem{}, err
	}

	return item, nil
}

// ToggleCommentLike 切换评论点赞；ctx 为请求上下文，userID/videoID/commentID 定位点赞目标。
func (service *Service) ToggleCommentLike(ctx context.Context, userID string, videoID string, commentID string) (CommentItem, error) {
	return service.repository.ToggleCommentLike(ctx, normalizeVideoID(videoID), strings.TrimSpace(commentID), normalizeUserID(userID))
}

// DeleteComment 删除当前用户自己的评论；ctx 为请求上下文，userID/videoID/commentID 定位评论。
func (service *Service) DeleteComment(ctx context.Context, userID string, videoID string, commentID string) error {
	return service.repository.DeleteComment(ctx, normalizeVideoID(videoID), strings.TrimSpace(commentID), normalizeUserID(userID))
}

// Danmaku 获取视频弹幕；ctx 为请求上下文，videoID 定位视频。
func (service *Service) Danmaku(ctx context.Context, videoID string) ([]DanmakuItem, error) {
	items, err := service.repository.ListDanmaku(ctx, normalizeVideoID(videoID))
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []DanmakuItem{}, nil
	}

	return items, nil
}

// AddDanmaku 新增弹幕；ctx 为请求上下文，userID/videoID 定位用户和视频，input 为弹幕内容和颜色。
func (service *Service) AddDanmaku(ctx context.Context, userID string, videoID string, input DanmakuInput) (DanmakuItem, error) {
	content := strings.TrimSpace(input.Content)
	if content == "" {
		return DanmakuItem{}, ErrInvalidContent
	}

	item := DanmakuItem{
		ID:        service.newID("danmaku"),
		VideoID:   normalizeVideoID(videoID),
		Content:   content,
		Color:     normalizeColor(input.Color),
		CreatedAt: service.now().UnixMilli(),
	}
	if err := service.repository.InsertDanmaku(ctx, normalizeUserID(userID), item); err != nil {
		return DanmakuItem{}, err
	}

	return item, nil
}

// SortComments 按排序规则返回新评论列表；items 为原始评论，sortMode 为 latest 或 hot。
func SortComments(items []CommentItem, sortMode CommentSort) []CommentItem {
	nextItems := append([]CommentItem(nil), items...)
	sort.SliceStable(nextItems, func(firstIndex int, secondIndex int) bool {
		first := nextItems[firstIndex]
		second := nextItems[secondIndex]
		if normalizeSort(sortMode) == CommentSortHot && first.Likes != second.Likes {
			return first.Likes > second.Likes
		}

		return first.CreatedAt > second.CreatedAt
	})

	return nextItems
}

// nextID 生成默认互动 id；prefix 为评论或弹幕前缀。
func nextID(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), defaultIDCounter.Add(1))
}

// normalizeSort 规范化评论排序；sortMode 为空或未知时使用 latest。
func normalizeSort(sortMode CommentSort) CommentSort {
	if sortMode == CommentSortHot {
		return CommentSortHot
	}

	return CommentSortLatest
}

// normalizeUserID 规范化用户标识；userID 为空时使用开发态默认用户。
func normalizeUserID(userID string) string {
	normalizedUserID := strings.TrimSpace(userID)
	if normalizedUserID == "" {
		return defaultUserID
	}

	return normalizedUserID
}

// normalizeVideoID 规范化视频标识；videoID 为路径中的视频 id。
func normalizeVideoID(videoID string) string {
	return strings.TrimSpace(videoID)
}

// normalizeColor 规范化弹幕颜色；color 不在前端支持范围内时返回 white。
func normalizeColor(color string) string {
	switch strings.TrimSpace(color) {
	case "green", "yellow", "pink", "white":
		return strings.TrimSpace(color)
	default:
		return defaultColor
	}
}
