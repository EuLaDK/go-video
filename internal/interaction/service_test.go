package interaction_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"next-video-golang/internal/interaction"
)

type fakeRepository struct {
	commentOwners map[string]string
	comments      map[string][]interaction.CommentItem
	danmaku       map[string][]interaction.DanmakuItem
	likedComments map[string]map[string]bool
}

// ListComments 返回测试评论列表；ctx 为调用上下文，videoID/userID/sort 用于定位和排序。
func (repo *fakeRepository) ListComments(ctx context.Context, videoID string, userID string, sort interaction.CommentSort) ([]interaction.CommentItem, error) {
	items := append([]interaction.CommentItem(nil), repo.comments[videoID]...)
	for index := range items {
		items[index].LikedByMe = repo.likedComments[items[index].ID][userID]
	}
	return interaction.SortComments(items, sort), nil
}

// InsertComment 写入测试评论；ctx 为调用上下文，userID 为作者，item 为评论数据。
func (repo *fakeRepository) InsertComment(ctx context.Context, userID string, item interaction.CommentItem) error {
	repo.commentOwners[item.ID] = userID
	repo.comments[item.VideoID] = append([]interaction.CommentItem{item}, repo.comments[item.VideoID]...)
	return nil
}

// ToggleCommentLike 切换测试评论点赞；ctx 为调用上下文，videoID/commentID/userID 定位点赞。
func (repo *fakeRepository) ToggleCommentLike(ctx context.Context, videoID string, commentID string, userID string) (interaction.CommentItem, error) {
	for index, item := range repo.comments[videoID] {
		if item.ID != commentID {
			continue
		}
		if repo.likedComments[commentID] == nil {
			repo.likedComments[commentID] = map[string]bool{}
		}
		if repo.likedComments[commentID][userID] {
			repo.likedComments[commentID][userID] = false
			item.LikedByMe = false
			item.Likes--
		} else {
			repo.likedComments[commentID][userID] = true
			item.LikedByMe = true
			item.Likes++
		}
		repo.comments[videoID][index] = item
		return item, nil
	}
	return interaction.CommentItem{}, interaction.ErrCommentNotFound
}

// DeleteComment 删除当前用户自己的测试评论；ctx 为调用上下文，videoID/commentID/userID 定位评论。
func (repo *fakeRepository) DeleteComment(ctx context.Context, videoID string, commentID string, userID string) error {
	if repo.commentOwners[commentID] != userID {
		return interaction.ErrCommentNotFound
	}
	nextItems := []interaction.CommentItem{}
	for _, item := range repo.comments[videoID] {
		if item.ID != commentID {
			nextItems = append(nextItems, item)
		}
	}
	repo.comments[videoID] = nextItems
	return nil
}

// ListDanmaku 返回测试弹幕列表；ctx 为调用上下文，videoID 定位视频。
func (repo *fakeRepository) ListDanmaku(ctx context.Context, videoID string) ([]interaction.DanmakuItem, error) {
	return append([]interaction.DanmakuItem(nil), repo.danmaku[videoID]...), nil
}

// InsertDanmaku 写入测试弹幕；ctx 为调用上下文，userID 为发送用户，item 为弹幕数据。
func (repo *fakeRepository) InsertDanmaku(ctx context.Context, userID string, item interaction.DanmakuItem) error {
	repo.danmaku[item.VideoID] = append([]interaction.DanmakuItem{item}, repo.danmaku[item.VideoID]...)
	return nil
}

func TestServiceAddsCommentsAndSortsByLatestOrHot(t *testing.T) {
	svc := interaction.NewService(newFakeRepository(), fixedClock(), fixedIDGenerator())

	first, err := svc.AddComment(context.Background(), " demo-user ", " xinghe ", interaction.CommentInput{Content: " 第一条 "})
	if err != nil {
		t.Fatal(err)
	}
	second, err := svc.AddComment(context.Background(), "demo-user", "xinghe", interaction.CommentInput{Content: "第二条"})
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != "comment-1" || first.Author != "我" || first.Content != "第一条" || first.CreatedAt != fixedUnixMillis {
		t.Fatalf("first = %#v", first)
	}

	_, err = svc.ToggleCommentLike(context.Background(), "other-user", "xinghe", first.ID)
	if err != nil {
		t.Fatal(err)
	}
	latest, err := svc.Comments(context.Background(), "demo-user", "xinghe", interaction.CommentSortLatest)
	if err != nil {
		t.Fatal(err)
	}
	hot, err := svc.Comments(context.Background(), "demo-user", "xinghe", interaction.CommentSortHot)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(commentIDs(latest), []string{second.ID, first.ID}) {
		t.Fatalf("latest ids = %#v", commentIDs(latest))
	}
	if !reflect.DeepEqual(commentIDs(hot), []string{first.ID, second.ID}) {
		t.Fatalf("hot ids = %#v", commentIDs(hot))
	}
}

func TestServiceRejectsEmptyCommentAndTogglesLike(t *testing.T) {
	svc := interaction.NewService(newFakeRepository(), fixedClock(), fixedIDGenerator())
	if _, err := svc.AddComment(context.Background(), "demo-user", "xinghe", interaction.CommentInput{Content: "   "}); err != interaction.ErrInvalidContent {
		t.Fatalf("err = %v, want ErrInvalidContent", err)
	}

	comment, err := svc.AddComment(context.Background(), "demo-user", "xinghe", interaction.CommentInput{Content: "值得二刷"})
	if err != nil {
		t.Fatal(err)
	}
	liked, err := svc.ToggleCommentLike(context.Background(), "demo-user", "xinghe", comment.ID)
	if err != nil {
		t.Fatal(err)
	}
	unliked, err := svc.ToggleCommentLike(context.Background(), "demo-user", "xinghe", comment.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !liked.LikedByMe || liked.Likes != 1 {
		t.Fatalf("liked = %#v", liked)
	}
	if unliked.LikedByMe || unliked.Likes != 0 {
		t.Fatalf("unliked = %#v", unliked)
	}
}

func TestServiceDeletesOnlyOwnComment(t *testing.T) {
	svc := interaction.NewService(newFakeRepository(), fixedClock(), fixedIDGenerator())
	comment, err := svc.AddComment(context.Background(), "demo-user", "xinghe", interaction.CommentInput{Content: "我的评论"})
	if err != nil {
		t.Fatal(err)
	}

	if err := svc.DeleteComment(context.Background(), "other-user", "xinghe", comment.ID); err != interaction.ErrCommentNotFound {
		t.Fatalf("delete other err = %v, want ErrCommentNotFound", err)
	}
	if err := svc.DeleteComment(context.Background(), "demo-user", "xinghe", comment.ID); err != nil {
		t.Fatal(err)
	}
	items, err := svc.Comments(context.Background(), "demo-user", "xinghe", interaction.CommentSortLatest)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("items after delete = %#v", items)
	}
}

func TestServiceAddsDanmakuWithColorFallback(t *testing.T) {
	svc := interaction.NewService(newFakeRepository(), fixedClock(), fixedIDGenerator())

	green, err := svc.AddDanmaku(context.Background(), "demo-user", "xinghe", interaction.DanmakuInput{Content: " 前方高能 ", Color: "green"})
	if err != nil {
		t.Fatal(err)
	}
	white, err := svc.AddDanmaku(context.Background(), "demo-user", "xinghe", interaction.DanmakuInput{Content: "默认颜色", Color: "blue"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AddDanmaku(context.Background(), "demo-user", "xinghe", interaction.DanmakuInput{Content: " "}); err != interaction.ErrInvalidContent {
		t.Fatalf("err = %v, want ErrInvalidContent", err)
	}

	items, err := svc.Danmaku(context.Background(), "xinghe")
	if err != nil {
		t.Fatal(err)
	}
	if green.ID != "danmaku-1" || green.Color != "green" || green.Content != "前方高能" {
		t.Fatalf("green = %#v", green)
	}
	if white.ID != "danmaku-2" || white.Color != "white" {
		t.Fatalf("white = %#v", white)
	}
	if !reflect.DeepEqual(danmakuIDs(items), []string{white.ID, green.ID}) {
		t.Fatalf("danmaku ids = %#v", danmakuIDs(items))
	}
}

const fixedUnixMillis = int64(1_766_000_000_000)

// fixedClock 返回固定时间；用于让服务层时间戳测试可重复。
func fixedClock() func() time.Time {
	return func() time.Time {
		return time.UnixMilli(fixedUnixMillis)
	}
}

// fixedIDGenerator 返回可重复 id 生成器；prefix 为评论或弹幕 id 前缀。
func fixedIDGenerator() func(prefix string) string {
	counters := map[string]int{}
	return func(prefix string) string {
		counters[prefix]++
		return prefix + "-" + string(rune('0'+counters[prefix]))
	}
}

// newFakeRepository 创建测试仓库；返回值带有空的评论、弹幕和点赞映射。
func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		commentOwners: map[string]string{},
		comments:      map[string][]interaction.CommentItem{},
		danmaku:       map[string][]interaction.DanmakuItem{},
		likedComments: map[string]map[string]bool{},
	}
}

// commentIDs 提取评论 id；items 为评论列表。
func commentIDs(items []interaction.CommentItem) []string {
	ids := []string{}
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

// danmakuIDs 提取弹幕 id；items 为弹幕列表。
func danmakuIDs(items []interaction.DanmakuItem) []string {
	ids := []string{}
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}
