package account_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"next-video-golang/internal/account"
)

type fakeRepository struct {
	users        map[string]account.UserProfile
	favorites    map[string][]account.FavoriteItem
	watchHistory map[string][]account.WatchHistoryItem
}

// GetUser 返回测试用户；ctx 为调用上下文，userID 为用户标识。
func (repo *fakeRepository) GetUser(ctx context.Context, userID string) (account.UserProfile, error) {
	if profile, ok := repo.users[userID]; ok {
		return profile, nil
	}

	return account.UserProfile{}, account.ErrUserNotFound
}

// UpsertUser 写入测试用户；ctx 为调用上下文，profile 为用户资料。
func (repo *fakeRepository) UpsertUser(ctx context.Context, profile account.UserProfile) error {
	repo.users[profile.ID] = profile
	return nil
}

// ListFavorites 返回测试收藏；ctx 为调用上下文，userID 为用户标识。
func (repo *fakeRepository) ListFavorites(ctx context.Context, userID string) ([]account.FavoriteItem, error) {
	return append([]account.FavoriteItem(nil), repo.favorites[userID]...), nil
}

// UpsertFavorite 写入测试收藏；ctx 为调用上下文，userID 为用户标识，item 为收藏项。
func (repo *fakeRepository) UpsertFavorite(ctx context.Context, userID string, item account.FavoriteItem) error {
	items := repo.favorites[userID]
	for index, existing := range items {
		if existing.ID == item.ID {
			items[index] = item
			repo.favorites[userID] = items
			return nil
		}
	}
	repo.favorites[userID] = append([]account.FavoriteItem{item}, items...)
	return nil
}

// DeleteFavorite 删除测试收藏；ctx 为调用上下文，userID 和 videoID 定位收藏。
func (repo *fakeRepository) DeleteFavorite(ctx context.Context, userID string, videoID string) error {
	items := repo.favorites[userID]
	nextItems := []account.FavoriteItem{}
	for _, item := range items {
		if item.ID != videoID {
			nextItems = append(nextItems, item)
		}
	}
	repo.favorites[userID] = nextItems
	return nil
}

// ListWatchHistory 返回测试历史；ctx 为调用上下文，userID 为用户标识。
func (repo *fakeRepository) ListWatchHistory(ctx context.Context, userID string) ([]account.WatchHistoryItem, error) {
	return append([]account.WatchHistoryItem(nil), repo.watchHistory[userID]...), nil
}

// UpsertWatchHistory 写入测试历史；ctx 为调用上下文，userID 为用户标识，item 为历史项。
func (repo *fakeRepository) UpsertWatchHistory(ctx context.Context, userID string, item account.WatchHistoryItem) error {
	items := repo.watchHistory[userID]
	nextItems := []account.WatchHistoryItem{item}
	for _, existing := range items {
		if existing.ID == item.ID && sameEpisode(existing.Episode, item.Episode) {
			continue
		}
		nextItems = append(nextItems, existing)
	}
	repo.watchHistory[userID] = nextItems
	return nil
}

// DeleteWatchHistory 删除测试历史；ctx 为调用上下文，userID、videoID 和 episode 定位历史。
func (repo *fakeRepository) DeleteWatchHistory(ctx context.Context, userID string, videoID string, episode *int) error {
	items := repo.watchHistory[userID]
	nextItems := []account.WatchHistoryItem{}
	for _, item := range items {
		if item.ID != videoID {
			nextItems = append(nextItems, item)
			continue
		}
		if episode != nil && !sameEpisode(item.Episode, episode) {
			nextItems = append(nextItems, item)
		}
	}
	repo.watchHistory[userID] = nextItems
	return nil
}

// ClearWatchHistory 清空测试历史；ctx 为调用上下文，userID 为用户标识。
func (repo *fakeRepository) ClearWatchHistory(ctx context.Context, userID string) error {
	repo.watchHistory[userID] = []account.WatchHistoryItem{}
	return nil
}

func TestServiceProfileCreatesDevelopmentUserWhenMissing(t *testing.T) {
	svc := account.NewService(newFakeRepository(), fixedClock())

	got, err := svc.Profile(context.Background(), "demo-user")
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != "demo-user" {
		t.Fatalf("ID = %q, want demo-user", got.ID)
	}
	if got.Nickname != "Next Video 用户" {
		t.Fatalf("Nickname = %q", got.Nickname)
	}
	if got.IsLoggedIn {
		t.Fatal("new development user should start logged out")
	}
}

func TestServiceLoginAndLogout(t *testing.T) {
	svc := account.NewService(newFakeRepository(), fixedClock())

	loggedIn, err := svc.Login(context.Background(), "demo-user", account.LoginInput{
		Nickname:  "  小夏  ",
		Contact:   "xia@example.com",
		AvatarURL: "  /avatar.png  ",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !loggedIn.IsLoggedIn || loggedIn.Nickname != "小夏" || loggedIn.Email != "xia@example.com" || loggedIn.AvatarURL != "/avatar.png" {
		t.Fatalf("loggedIn = %#v", loggedIn)
	}

	loggedOut, err := svc.Logout(context.Background(), "demo-user")
	if err != nil {
		t.Fatal(err)
	}

	if loggedOut.IsLoggedIn || loggedOut.IsVip {
		t.Fatalf("loggedOut = %#v", loggedOut)
	}
}

func TestServiceFavorites(t *testing.T) {
	svc := account.NewService(newFakeRepository(), fixedClock())

	item, err := svc.AddFavorite(context.Background(), "demo-user", account.FavoriteInput{
		ID:            "xinghe",
		Title:         "星河回响",
		Category:      "科幻 / 悬疑",
		Progress:      "会员抢先看",
		CoverGradient: "linear-gradient(...)",
		Description:   "深空信号",
	})
	if err != nil {
		t.Fatal(err)
	}

	if item.AddedAt != fixedUnixMillis {
		t.Fatalf("AddedAt = %d, want %d", item.AddedAt, fixedUnixMillis)
	}

	items, err := svc.Favorites(context.Background(), "demo-user")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != "xinghe" {
		t.Fatalf("items = %#v", items)
	}

	if err := svc.DeleteFavorite(context.Background(), "demo-user", "xinghe"); err != nil {
		t.Fatal(err)
	}
	items, _ = svc.Favorites(context.Background(), "demo-user")
	if len(items) != 0 {
		t.Fatalf("items after delete = %#v", items)
	}
}

func TestServiceWatchHistory(t *testing.T) {
	svc := account.NewService(newFakeRepository(), fixedClock())
	episode := 2
	watchSeconds := 90
	durationSeconds := 2700

	item, err := svc.AddWatchHistory(context.Background(), "demo-user", account.WatchHistoryInput{
		ID:              "xinghe",
		Title:           "星河回响",
		Category:        "科幻 / 悬疑",
		Progress:        "会员抢先看",
		CoverGradient:   "linear-gradient(...)",
		Episode:         &episode,
		WatchSeconds:    &watchSeconds,
		DurationSeconds: &durationSeconds,
	})
	if err != nil {
		t.Fatal(err)
	}

	if item.WatchedAt != fixedUnixMillis {
		t.Fatalf("WatchedAt = %d, want %d", item.WatchedAt, fixedUnixMillis)
	}

	items, err := svc.WatchHistory(context.Background(), "demo-user")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != "xinghe" || *items[0].Episode != 2 {
		t.Fatalf("items = %#v", items)
	}

	if err := svc.DeleteWatchHistory(context.Background(), "demo-user", "xinghe", &episode); err != nil {
		t.Fatal(err)
	}
	items, _ = svc.WatchHistory(context.Background(), "demo-user")
	if len(items) != 0 {
		t.Fatalf("items after delete = %#v", items)
	}

	_, _ = svc.AddWatchHistory(context.Background(), "demo-user", account.WatchHistoryInput{ID: "anye", Title: "暗夜追光"})
	if err := svc.ClearWatchHistory(context.Background(), "demo-user"); err != nil {
		t.Fatal(err)
	}
	items, _ = svc.WatchHistory(context.Background(), "demo-user")
	if !reflect.DeepEqual(items, []account.WatchHistoryItem{}) {
		t.Fatalf("items after clear = %#v", items)
	}
}

const fixedUnixMillis = int64(1_766_000_000_000)

// fixedClock 返回固定时间；用于让服务层时间戳测试可重复。
func fixedClock() func() time.Time {
	return func() time.Time {
		return time.UnixMilli(fixedUnixMillis)
	}
}

// newFakeRepository 创建测试仓库；返回值带有空的用户、收藏和历史映射。
func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		users:        map[string]account.UserProfile{},
		favorites:    map[string][]account.FavoriteItem{},
		watchHistory: map[string][]account.WatchHistoryItem{},
	}
}

// sameEpisode 判断两个可选集数是否相同；first 和 second 为待比较集数。
func sameEpisode(first *int, second *int) bool {
	if first == nil || second == nil {
		return first == second
	}

	return *first == *second
}
