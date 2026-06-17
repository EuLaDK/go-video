package httpapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"next-video-golang/internal/httpapi"
	"next-video-golang/internal/interaction"
)

type fakeInteractionService struct {
	userID   string
	videoID  string
	sort     interaction.CommentSort
	comments []interaction.CommentItem
	danmaku  []interaction.DanmakuItem
}

// Comments 模拟评论列表；ctx 为请求上下文，userID/videoID/sort 用于定位和排序。
func (service *fakeInteractionService) Comments(ctx context.Context, userID string, videoID string, sort interaction.CommentSort) ([]interaction.CommentItem, error) {
	service.userID = userID
	service.videoID = videoID
	service.sort = sort
	return service.comments, nil
}

// AddComment 模拟新增评论；ctx 为请求上下文，userID/videoID/input 用于创建评论。
func (service *fakeInteractionService) AddComment(ctx context.Context, userID string, videoID string, input interaction.CommentInput) (interaction.CommentItem, error) {
	service.userID = userID
	service.videoID = videoID
	if input.Content == "" {
		return interaction.CommentItem{}, interaction.ErrInvalidContent
	}
	item := interaction.CommentItem{ID: "comment-1", VideoID: videoID, Content: input.Content, Author: "我", CreatedAt: 123}
	service.comments = []interaction.CommentItem{item}
	return item, nil
}

// ToggleCommentLike 模拟评论点赞；ctx 为请求上下文，userID/videoID/commentID 定位评论。
func (service *fakeInteractionService) ToggleCommentLike(ctx context.Context, userID string, videoID string, commentID string) (interaction.CommentItem, error) {
	service.userID = userID
	service.videoID = videoID
	item := interaction.CommentItem{ID: commentID, VideoID: videoID, Content: "值得二刷", Author: "我", LikedByMe: true, Likes: 1, CreatedAt: 123}
	service.comments = []interaction.CommentItem{item}
	return item, nil
}

// DeleteComment 模拟删除评论；ctx 为请求上下文，userID/videoID/commentID 定位评论。
func (service *fakeInteractionService) DeleteComment(ctx context.Context, userID string, videoID string, commentID string) error {
	service.userID = userID
	service.videoID = videoID
	service.comments = []interaction.CommentItem{}
	return nil
}

// Danmaku 模拟弹幕列表；ctx 为请求上下文，videoID 定位视频。
func (service *fakeInteractionService) Danmaku(ctx context.Context, videoID string) ([]interaction.DanmakuItem, error) {
	service.videoID = videoID
	return service.danmaku, nil
}

// AddDanmaku 模拟新增弹幕；ctx 为请求上下文，userID/videoID/input 用于创建弹幕。
func (service *fakeInteractionService) AddDanmaku(ctx context.Context, userID string, videoID string, input interaction.DanmakuInput) (interaction.DanmakuItem, error) {
	service.userID = userID
	service.videoID = videoID
	if input.Content == "" {
		return interaction.DanmakuItem{}, interaction.ErrInvalidContent
	}
	item := interaction.DanmakuItem{ID: "danmaku-1", VideoID: videoID, Content: input.Content, Color: input.Color, CreatedAt: 456}
	service.danmaku = []interaction.DanmakuItem{item}
	return item, nil
}

func TestInteractionRoutesComments(t *testing.T) {
	interactionService := &fakeInteractionService{}
	srv := httpapi.NewServerWithServices(&fakeVideoService{}, nil, interactionService)

	rec := requestJSON(srv, http.MethodPost, "/videos/xinghe/comments", `{"content":"值得二刷"}`, "custom-user")
	if rec.Code != http.StatusOK {
		t.Fatalf("add comment status = %d, want 200", rec.Code)
	}
	var comment interaction.CommentItem
	if err := json.Unmarshal(rec.Body.Bytes(), &comment); err != nil {
		t.Fatal(err)
	}
	if comment.ID != "comment-1" || comment.VideoID != "xinghe" || interactionService.userID != "custom-user" {
		t.Fatalf("comment = %#v userID = %q", comment, interactionService.userID)
	}

	rec = requestJSON(srv, http.MethodGet, "/videos/xinghe/comments?sort=hot", "", "custom-user")
	var comments []interaction.CommentItem
	if err := json.Unmarshal(rec.Body.Bytes(), &comments); err != nil {
		t.Fatal(err)
	}
	if len(comments) != 1 || interactionService.sort != interaction.CommentSortHot {
		t.Fatalf("comments = %#v sort = %q", comments, interactionService.sort)
	}

	rec = requestJSON(srv, http.MethodPost, "/videos/xinghe/comments/comment-1/like", "", "custom-user")
	if rec.Code != http.StatusOK {
		t.Fatalf("like status = %d, want 200", rec.Code)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &comment); err != nil {
		t.Fatal(err)
	}
	if !comment.LikedByMe || comment.Likes != 1 {
		t.Fatalf("liked comment = %#v", comment)
	}

	rec = requestJSON(srv, http.MethodDelete, "/videos/xinghe/comments/comment-1", "", "custom-user")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete comment status = %d, want 204", rec.Code)
	}
	if !reflect.DeepEqual(interactionService.comments, []interaction.CommentItem{}) {
		t.Fatalf("comments after delete = %#v", interactionService.comments)
	}
}

func TestInteractionRoutesDanmaku(t *testing.T) {
	interactionService := &fakeInteractionService{}
	srv := httpapi.NewServerWithServices(&fakeVideoService{}, nil, interactionService)

	rec := requestJSON(srv, http.MethodPost, "/videos/xinghe/danmaku", `{"content":"前方高能","color":"green"}`, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("add danmaku status = %d, want 200", rec.Code)
	}
	var danmaku interaction.DanmakuItem
	if err := json.Unmarshal(rec.Body.Bytes(), &danmaku); err != nil {
		t.Fatal(err)
	}
	if danmaku.ID != "danmaku-1" || danmaku.VideoID != "xinghe" || interactionService.userID != "demo-user" {
		t.Fatalf("danmaku = %#v userID = %q", danmaku, interactionService.userID)
	}

	rec = requestJSON(srv, http.MethodGet, "/videos/xinghe/danmaku", "", "")
	var items []interaction.DanmakuItem
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Content != "前方高能" {
		t.Fatalf("items = %#v", items)
	}
}

func TestInteractionRoutesValidationErrors(t *testing.T) {
	interactionService := &fakeInteractionService{}
	srv := httpapi.NewServerWithServices(&fakeVideoService{}, nil, interactionService)

	rec := requestJSON(srv, http.MethodPost, "/videos/xinghe/comments", `{"content":""}`, "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("empty comment status = %d, want 400", rec.Code)
	}

	rec = requestJSON(srv, http.MethodPost, "/videos/xinghe/danmaku", `{"content":""}`, "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("empty danmaku status = %d, want 400", rec.Code)
	}
}
