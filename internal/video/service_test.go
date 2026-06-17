package video_test

import (
	"context"
	"reflect"
	"testing"

	"next-video-golang/internal/video"
)

type fakeRepository struct {
	channels []video.Channel
	videos   []video.Video
}

// ListChannels 返回测试用频道列表；ctx 为调用上下文。
func (repo fakeRepository) ListChannels(ctx context.Context) ([]video.Channel, error) {
	return repo.channels, nil
}

// ListVideos 返回测试用视频列表；ctx 为调用上下文。
func (repo fakeRepository) ListVideos(ctx context.Context) ([]video.Video, error) {
	return repo.videos, nil
}

func TestServiceHomePageUsesDefaultSlices(t *testing.T) {
	svc := video.NewService(sampleRepository())

	got, err := svc.HomePage(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if got.FeaturedVideo.ID != "xinghe" {
		t.Fatalf("FeaturedVideo.ID = %q, want xinghe", got.FeaturedVideo.ID)
	}

	if idsOf(got.RecommendationVideos) != "movie-a,drama-a,vip-a,doc-a" {
		t.Fatalf("recommendation ids = %s", idsOf(got.RecommendationVideos))
	}

	if len(got.RankVideos) != 3 {
		t.Fatalf("len(RankVideos) = %d, want 3", len(got.RankVideos))
	}
}

func TestServiceRankVideosSortsAndFilters(t *testing.T) {
	svc := video.NewService(sampleRepository())

	got, err := svc.RankVideos(context.Background(), video.RankQuery{
		Channel: "movie",
		Sort:    "score",
	})
	if err != nil {
		t.Fatal(err)
	}

	if idsOf(got) != "vip-a,movie-a" {
		t.Fatalf("rank ids = %s, want vip-a,movie-a", idsOf(got))
	}
}

func TestServiceChannelPageFiltersAndFallsBack(t *testing.T) {
	svc := video.NewService(sampleRepository())

	got, err := svc.ChannelPage(context.Background(), video.ChannelQuery{
		Slug: "missing",
		Sort: "hot",
	})
	if err != nil {
		t.Fatal(err)
	}

	if got.Channel.Slug != "featured" {
		t.Fatalf("Channel.Slug = %q, want featured", got.Channel.Slug)
	}

	if got.HeroVideo.ID != got.Videos[0].ID {
		t.Fatalf("HeroVideo.ID = %q, want first video %q", got.HeroVideo.ID, got.Videos[0].ID)
	}
}

func TestServiceSearchPageFiltersAndSorts(t *testing.T) {
	svc := video.NewService(sampleRepository())

	got, err := svc.SearchPage(context.Background(), video.SearchQuery{
		Q:       "电影",
		Quality: "4K",
		Sort:    "hot",
	})
	if err != nil {
		t.Fatal(err)
	}

	if idsOf(got.Videos) != "vip-a,movie-a" {
		t.Fatalf("search ids = %s, want vip-a,movie-a", idsOf(got.Videos))
	}

	if !reflect.DeepEqual(got.HotSearchKeywords, []string{"星河回响", "电影", "悬疑", "纪录片", "会员", "科幻"}) {
		t.Fatalf("HotSearchKeywords = %#v", got.HotSearchKeywords)
	}
}

func TestServiceWatchPageAndVideoIDs(t *testing.T) {
	svc := video.NewService(sampleRepository())

	page, err := svc.WatchPage(context.Background(), "vip-a")
	if err != nil {
		t.Fatal(err)
	}

	if page.Video.ID != "vip-a" {
		t.Fatalf("Video.ID = %q, want vip-a", page.Video.ID)
	}

	if idsOf(page.RelatedVideos) != "xinghe,movie-a,drama-a,doc-a" {
		t.Fatalf("related ids = %s", idsOf(page.RelatedVideos))
	}

	ids, err := svc.VideoIDs(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(ids, []string{"xinghe", "movie-a", "drama-a", "vip-a", "doc-a"}) {
		t.Fatalf("ids = %#v", ids)
	}
}

func TestServiceWatchPageFallsBackToFeaturedVideo(t *testing.T) {
	svc := video.NewService(sampleRepository())

	got, err := svc.WatchPage(context.Background(), "unknown")
	if err != nil {
		t.Fatal(err)
	}

	if got.Video.ID != "xinghe" {
		t.Fatalf("Video.ID = %q, want xinghe", got.Video.ID)
	}
}

// sampleRepository 构造测试用仓库；数据顺序模拟前端 mock 数据的默认顺序。
func sampleRepository() fakeRepository {
	return fakeRepository{
		channels: []video.Channel{
			{Slug: "featured", Label: "精选", Description: "精选内容", Keywords: []string{}},
			{Slug: "movie", Label: "电影", Description: "电影频道", Keywords: []string{"电影"}},
			{Slug: "tv", Label: "电视剧", Description: "剧集频道", Keywords: []string{"剧集", "都市"}},
			{Slug: "vip", Label: "VIP", Description: "会员频道", Keywords: []string{"会员", "独播", "会员抢先看"}},
		},
		videos: []video.Video{
			{
				ID:              "xinghe",
				Title:           "星河回响",
				Subtitle:        "科幻悬疑",
				Description:     "深空信号",
				Score:           "9.3",
				Heat:            "热度 10026",
				Update:          "更新至 18 集",
				Category:        "科幻 / 悬疑",
				Year:            "2026",
				Region:          "中国大陆",
				TotalEpisodes:   24,
				Quality:         "4K HDR",
				Badge:           "独播",
				Progress:        "会员抢先看",
				Duration:        "45:00",
				SourceURL:       "/assets/video/staticTest.mp4",
				CoverGradient:   "linear-gradient(135deg,#0f766e,#111827,#be123c)",
				Tags:            []string{"科幻", "悬疑", "会员抢先看"},
				CastNames:       []string{"林舟"},
				RelatedVideoIDs: []string{"movie-a", "drama-a", "vip-a"},
			},
			{
				ID:              "movie-a",
				Title:           "归途列车",
				Subtitle:        "电影冒险",
				Description:     "暴雪列车",
				Score:           "8.1",
				Heat:            "热度 6920",
				Update:          "本周新片",
				Category:        "电影 / 冒险",
				Year:            "2026",
				Region:          "中国大陆",
				TotalEpisodes:   1,
				Quality:         "4K",
				Badge:           "本周新片",
				Progress:        "本周新片",
				Duration:        "01:48:20",
				SourceURL:       "/assets/video/staticTest.mp4",
				CoverGradient:   "linear-gradient(135deg,#be123c,#312e81,#111827)",
				Tags:            []string{"电影", "冒险"},
				CastNames:       []string{"梁舟"},
				RelatedVideoIDs: []string{"xinghe", "vip-a"},
			},
			{
				ID:              "drama-a",
				Title:           "春日事务所",
				Subtitle:        "剧集都市",
				Description:     "旧街事务所",
				Score:           "8.2",
				Heat:            "热度 6638",
				Update:          "更新至 12 集",
				Category:        "剧集 / 都市",
				Year:            "2026",
				Region:          "中国大陆",
				TotalEpisodes:   24,
				Quality:         "1080P",
				Badge:           "轻喜热播",
				Progress:        "更新至 12 集",
				Duration:        "40:10",
				SourceURL:       "/assets/video/staticTest.mp4",
				CoverGradient:   "linear-gradient(135deg,#15803d,#0f172a,#1f2937)",
				Tags:            []string{"剧集", "都市"},
				CastNames:       []string{"唐青"},
				RelatedVideoIDs: []string{"xinghe"},
			},
			{
				ID:              "vip-a",
				Title:           "零点航线",
				Subtitle:        "电影会员抢先看",
				Description:     "跨海航班",
				Score:           "8.4",
				Heat:            "热度 9000",
				Update:          "会员抢先看",
				Category:        "电影 / 灾难",
				Year:            "2026",
				Region:          "中国大陆",
				TotalEpisodes:   1,
				Quality:         "4K",
				Badge:           "会员抢先看",
				Progress:        "会员抢先看",
				Duration:        "39:46",
				SourceURL:       "/assets/video/staticTest.mp4",
				CoverGradient:   "linear-gradient(135deg,#be123c,#312e81,#111827)",
				Tags:            []string{"电影", "会员"},
				CastNames:       []string{"江川"},
				RelatedVideoIDs: []string{"xinghe", "movie-a", "drama-a", "doc-a"},
			},
			{
				ID:              "doc-a",
				Title:           "海岸来信",
				Subtitle:        "纪录片",
				Description:     "海岛人文",
				Score:           "8.8",
				Heat:            "热度 7388",
				Update:          "全 6 集已上线",
				Category:        "纪录片",
				Year:            "2026",
				Region:          "中国大陆",
				TotalEpisodes:   6,
				Quality:         "4K HDR",
				Badge:           "口碑纪录",
				Progress:        "4K HDR",
				Duration:        "48:08",
				SourceURL:       "/assets/video/staticTest.mp4",
				CoverGradient:   "linear-gradient(135deg,#0f766e,#172554,#111827)",
				Tags:            []string{"纪录片", "自然"},
				CastNames:       []string{"周闻"},
				RelatedVideoIDs: []string{"xinghe"},
			},
		},
	}
}

// idsOf 将视频列表转成逗号分隔 id；videos 为待检查的视频列表。
func idsOf(videos []video.Video) string {
	ids := ""
	for index, item := range videos {
		if index > 0 {
			ids += ","
		}
		ids += item.ID
	}

	return ids
}
