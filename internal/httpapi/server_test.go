package httpapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"next-video-golang/internal/httpapi"
	"next-video-golang/internal/video"
)

type fakeVideoService struct {
	err          error
	rankQuery    video.RankQuery
	channelQuery video.ChannelQuery
	searchQuery  video.SearchQuery
	watchID      string
}

// HealthCheck 模拟健康检查；ctx 为请求上下文。
func (service *fakeVideoService) HealthCheck(ctx context.Context) error {
	return service.err
}

// HomePage 模拟首页数据；ctx 为请求上下文。
func (service *fakeVideoService) HomePage(ctx context.Context) (video.HomePageData, error) {
	if service.err != nil {
		return video.HomePageData{}, service.err
	}
	return video.HomePageData{
		FeaturedVideo: sampleVideo("xinghe"),
	}, nil
}

// RankVideos 模拟排行榜数据；ctx 为请求上下文，query 为筛选条件。
func (service *fakeVideoService) RankVideos(ctx context.Context, query video.RankQuery) ([]video.Video, error) {
	service.rankQuery = query
	return []video.Video{sampleVideo("rank-a")}, service.err
}

// ChannelPage 模拟频道页数据；ctx 为请求上下文，query 为筛选条件。
func (service *fakeVideoService) ChannelPage(ctx context.Context, query video.ChannelQuery) (video.ChannelPageData, error) {
	service.channelQuery = query
	return video.ChannelPageData{
		Channel:   video.Channel{Slug: query.Slug, Label: "电影"},
		HeroVideo: sampleVideo("hero-a"),
		Videos:    []video.Video{sampleVideo("channel-a")},
	}, service.err
}

// SearchPage 模拟搜索页数据；ctx 为请求上下文，query 为筛选条件。
func (service *fakeVideoService) SearchPage(ctx context.Context, query video.SearchQuery) (video.SearchPageData, error) {
	service.searchQuery = query
	return video.SearchPageData{
		HotSearchKeywords: []string{"星河回响"},
		Videos:            []video.Video{sampleVideo("search-a")},
	}, service.err
}

// WatchPage 模拟播放详情数据；ctx 为请求上下文，videoID 为视频 id。
func (service *fakeVideoService) WatchPage(ctx context.Context, videoID string) (video.WatchPageData, error) {
	service.watchID = videoID
	return video.WatchPageData{
		Video:         sampleVideo(videoID),
		RelatedVideos: []video.Video{sampleVideo("related-a")},
	}, service.err
}

// VideoIDs 模拟静态视频 id；ctx 为请求上下文。
func (service *fakeVideoService) VideoIDs(ctx context.Context) ([]string, error) {
	return []string{"xinghe", "movie-a"}, service.err
}

func TestServerHomeRouteReturnsJSON(t *testing.T) {
	srv := httpapi.NewServer(&fakeVideoService{})
	rec := get(srv, "/videos/home")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var body video.HomePageData
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}

	if body.FeaturedVideo.ID != "xinghe" {
		t.Fatalf("FeaturedVideo.ID = %q, want xinghe", body.FeaturedVideo.ID)
	}
}

func TestServerRoutesQueryParameters(t *testing.T) {
	service := &fakeVideoService{}
	srv := httpapi.NewServer(service)

	get(srv, "/videos/rank?sort=score&channel=movie")
	if service.rankQuery != (video.RankQuery{Sort: "score", Channel: "movie"}) {
		t.Fatalf("rankQuery = %#v", service.rankQuery)
	}

	get(srv, "/videos/channel/movie?type=%E7%94%B5%E5%BD%B1&year=2026&sort=hot")
	if service.channelQuery != (video.ChannelQuery{Slug: "movie", Type: "电影", Year: "2026", Sort: "hot"}) {
		t.Fatalf("channelQuery = %#v", service.channelQuery)
	}

	get(srv, "/videos/search?q=%E7%A7%91%E5%B9%BB&quality=4K&channel=vip&type=%E7%94%B5%E5%BD%B1&year=2026&sort=hot")
	wantSearch := video.SearchQuery{Q: "科幻", Quality: "4K", Channel: "vip", Type: "电影", Year: "2026", Sort: "hot"}
	if service.searchQuery != wantSearch {
		t.Fatalf("searchQuery = %#v", service.searchQuery)
	}

	get(srv, "/videos/xinghe")
	if service.watchID != "xinghe" {
		t.Fatalf("watchID = %q, want xinghe", service.watchID)
	}
}

func TestServerVideoIDsRouteReturnsArray(t *testing.T) {
	srv := httpapi.NewServer(&fakeVideoService{})
	rec := get(srv, "/videos/ids")

	var body []string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(body, []string{"xinghe", "movie-a"}) {
		t.Fatalf("body = %#v", body)
	}
}

func TestServerHealthAndErrorResponses(t *testing.T) {
	healthy := httpapi.NewServer(&fakeVideoService{})
	rec := get(healthy, "/health")
	if rec.Code != http.StatusOK {
		t.Fatalf("healthy status = %d, want 200", rec.Code)
	}

	unhealthy := httpapi.NewServer(&fakeVideoService{err: errors.New("database down")})
	rec = get(unhealthy, "/health")
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("unhealthy status = %d, want 503", rec.Code)
	}

	rec = get(unhealthy, "/videos/home")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("api error status = %d, want 500", rec.Code)
	}
}

func TestServerCORSAndMethodGuard(t *testing.T) {
	srv := httpapi.NewServer(&fakeVideoService{})

	optionsRequest := httptest.NewRequest(http.MethodOptions, "/videos/home", nil)
	optionsRecorder := httptest.NewRecorder()
	srv.ServeHTTP(optionsRecorder, optionsRequest)
	if optionsRecorder.Code != http.StatusNoContent {
		t.Fatalf("options status = %d, want 204", optionsRecorder.Code)
	}

	if optionsRecorder.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("missing CORS allow origin header")
	}

	postRequest := httptest.NewRequest(http.MethodPost, "/videos/home", nil)
	postRecorder := httptest.NewRecorder()
	srv.ServeHTTP(postRecorder, postRequest)
	if postRecorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("post status = %d, want 405", postRecorder.Code)
	}
}

// get 向测试服务发送 GET 请求；handler 为被测 HTTP 服务，target 为请求路径。
func get(handler http.Handler, target string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

// sampleVideo 构造测试视频；id 为视频唯一标识。
func sampleVideo(id string) video.Video {
	return video.Video{
		ID:        id,
		Title:     "测试视频",
		SourceURL: "/assets/video/staticTest.mp4",
		Tags:      []string{"测试"},
	}
}
