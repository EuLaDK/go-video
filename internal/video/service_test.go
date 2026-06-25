package video_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"next-video-golang/internal/video"
)

type fakeRepository struct {
	channels        []video.Channel
	playbackSources map[string][]video.PlaybackSource
	videos          []video.Video
}

// ListChannels 返回测试用频道列表；ctx 为调用上下文。
func (repo fakeRepository) ListChannels(ctx context.Context) ([]video.Channel, error) {
	return repo.channels, nil
}

// ListVideos 返回测试用视频列表；ctx 为调用上下文。
func (repo fakeRepository) ListVideos(ctx context.Context) ([]video.Video, error) {
	return repo.videos, nil
}

// ListPlaybackSources 返回测试用播放源列表；ctx 为调用上下文，videoID 为当前视频 id。
func (repo fakeRepository) ListPlaybackSources(ctx context.Context, videoID string) ([]video.PlaybackSource, error) {
	return repo.playbackSources[videoID], nil
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

func TestServiceWatchPageReturnsPlaybackConfig(t *testing.T) {
	svc := video.NewService(sampleRepository())

	page, err := svc.WatchPage(context.Background(), "vip-a")
	if err != nil {
		t.Fatal(err)
	}

	var payload struct {
		Playback struct {
			Sources []struct {
				Quality   string `json:"quality"`
				Label     string `json:"label"`
				SourceURL string `json:"sourceUrl"`
				MimeType  string `json:"mimeType"`
			} `json:"sources"`
			DefaultQuality string `json:"defaultQuality"`
			RequiresVIP    bool   `json:"requiresVip"`
			CanPlay        bool   `json:"canPlay"`
			TrialSeconds   int    `json:"trialSeconds"`
			Message        string `json:"message"`
		} `json:"playback"`
	}
	body, err := json.Marshal(page)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatal(err)
	}

	if len(payload.Playback.Sources) != 1 {
		t.Fatalf("len(playback.sources) = %d, want 1", len(payload.Playback.Sources))
	}
	source := payload.Playback.Sources[0]
	if source.SourceURL != "/assets/video/staticTest.mp4" {
		t.Fatalf("sourceUrl = %q, want /assets/video/staticTest.mp4", source.SourceURL)
	}
	if source.Quality != "4K" || source.Label != "4K" {
		t.Fatalf("source quality/label = %q/%q, want 4K/4K", source.Quality, source.Label)
	}
	if source.MimeType != "video/mp4" {
		t.Fatalf("mimeType = %q, want video/mp4", source.MimeType)
	}
	if payload.Playback.DefaultQuality != "4K" {
		t.Fatalf("defaultQuality = %q, want 4K", payload.Playback.DefaultQuality)
	}
	if !payload.Playback.RequiresVIP {
		t.Fatalf("requiresVip = false, want true")
	}
	if payload.Playback.CanPlay {
		t.Fatalf("canPlay = true, want false for non-VIP viewer")
	}
	if payload.Playback.TrialSeconds != 360 {
		t.Fatalf("trialSeconds = %d, want 360", payload.Playback.TrialSeconds)
	}
	if payload.Playback.Message == "" {
		t.Fatalf("message is empty, want VIP playback prompt")
	}
}

func TestServiceWatchPageAllowsVipViewerToPlayVipContent(t *testing.T) {
	svc := video.NewService(sampleRepository())
	ctx := video.ContextWithPlaybackViewer(context.Background(), video.PlaybackViewer{
		IsVIP: true,
	})

	page, err := svc.WatchPage(ctx, "vip-a")
	if err != nil {
		t.Fatal(err)
	}

	if !page.Playback.RequiresVIP {
		t.Fatalf("requiresVip = false, want true")
	}
	if !page.Playback.CanPlay {
		t.Fatalf("canPlay = false, want true for VIP viewer")
	}
	if page.Playback.TrialSeconds != 0 {
		t.Fatalf("trialSeconds = %d, want 0 for VIP viewer", page.Playback.TrialSeconds)
	}
	if page.Playback.Message == "" {
		t.Fatalf("message is empty, want VIP playback message")
	}
}

func TestServiceWatchPageReturnsResumePointFromContext(t *testing.T) {
	svc := video.NewService(sampleRepository())
	ctx := video.ContextWithPlaybackViewer(context.Background(), video.PlaybackViewer{
		Resume: video.PlaybackResume{
			CanResume:       true,
			Episode:         2,
			WatchSeconds:    90,
			DurationSeconds: 2700,
		},
	})

	page, err := svc.WatchPage(ctx, "xinghe")
	if err != nil {
		t.Fatal(err)
	}

	if !page.Playback.Resume.CanResume {
		t.Fatalf("resume canResume = false, want true")
	}
	if page.Playback.Resume.Episode != 2 || page.Playback.Resume.WatchSeconds != 90 || page.Playback.Resume.DurationSeconds != 2700 {
		t.Fatalf("resume = %#v", page.Playback.Resume)
	}
}

func TestServiceWatchPageUsesRepositoryPlaybackSources(t *testing.T) {
	repository := sampleRepository()
	repository.playbackSources = map[string][]video.PlaybackSource{
		"vip-a": {
			{
				Quality:   "4K",
				Label:     "4K 超清",
				SourceURL: "/media/vip-a-4k.mp4",
				MimeType:  "video/mp4",
			},
			{
				Quality:   "1080P",
				Label:     "1080P 高清",
				SourceURL: "/media/vip-a-1080p.mp4",
				MimeType:  "video/mp4",
			},
		},
	}
	svc := video.NewService(repository)

	page, err := svc.WatchPage(context.Background(), "vip-a")
	if err != nil {
		t.Fatal(err)
	}

	if len(page.Playback.Sources) != 2 {
		t.Fatalf("len(playback.sources) = %d, want 2", len(page.Playback.Sources))
	}
	if page.Playback.Sources[0].SourceURL != "/media/vip-a-4k.mp4" {
		t.Fatalf("first sourceUrl = %q, want repository source", page.Playback.Sources[0].SourceURL)
	}
	if page.Playback.Sources[1].Quality != "1080P" {
		t.Fatalf("second quality = %q, want 1080P", page.Playback.Sources[1].Quality)
	}
	if page.Playback.DefaultQuality != "4K" {
		t.Fatalf("defaultQuality = %q, want 4K", page.Playback.DefaultQuality)
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
