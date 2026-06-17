package video

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var defaultHotSearchKeywords = []string{"星河回响", "电影", "悬疑", "纪录片", "会员", "科幻"}

type Repository interface {
	ListChannels(ctx context.Context) ([]Channel, error)
	ListVideos(ctx context.Context) ([]Video, error)
}

type Service struct {
	repository Repository
}

// NewService 创建视频服务；repository 提供频道和视频的持久化读取能力。
func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

// HomePage 获取首页聚合数据；ctx 为请求上下文。
func (service *Service) HomePage(ctx context.Context) (HomePageData, error) {
	videos, err := service.repository.ListVideos(ctx)
	if err != nil {
		return HomePageData{}, err
	}

	featured := firstVideo(videos)

	return HomePageData{
		FeaturedVideo:        featured,
		HotVideos:            sliceVideos(videos, 5, 10),
		RankVideos:           limitVideos(sortRankVideos(videos, "hot"), 3),
		RecommendationVideos: sliceVideos(videos, 1, 5),
	}, nil
}

// RankVideos 获取排行榜视频；ctx 为请求上下文，query 描述榜单和频道筛选。
func (service *Service) RankVideos(ctx context.Context, query RankQuery) ([]Video, error) {
	channels, videos, err := service.loadCatalog(ctx)
	if err != nil {
		return nil, err
	}

	filtered := filterVideos(videos, func(item Video) bool {
		return matchesVideoChannel(item, query.Channel, channels)
	})

	return sortRankVideos(filtered, withDefault(query.Sort, "hot")), nil
}

// ChannelPage 获取频道页数据；ctx 为请求上下文，query 描述频道和筛选条件。
func (service *Service) ChannelPage(ctx context.Context, query ChannelQuery) (ChannelPageData, error) {
	channels, videos, err := service.loadCatalog(ctx)
	if err != nil {
		return ChannelPageData{}, err
	}

	channel := getChannelBySlug(channels, query.Slug)
	channelVideos := getVideosByChannel(videos, channel)
	filtered := filterVideos(channelVideos, func(item Video) bool {
		return matchesVideoFilterKeyword(item, query.Type) &&
			(query.Year == "" || item.Year == query.Year)
	})
	sortedVideos := sortChannelVideos(filtered, withDefault(query.Sort, "default"))

	hero := firstVideo(sortedVideos)
	if hero.ID == "" {
		hero = firstVideo(channelVideos)
	}
	if hero.ID == "" {
		hero = firstVideo(videos)
	}

	return ChannelPageData{
		Channel:   channel,
		HeroVideo: hero,
		Videos:    sortedVideos,
	}, nil
}

// SearchPage 获取搜索页数据；ctx 为请求上下文，query 描述关键词和筛选条件。
func (service *Service) SearchPage(ctx context.Context, query SearchQuery) (SearchPageData, error) {
	channels, videos, err := service.loadCatalog(ctx)
	if err != nil {
		return SearchPageData{}, err
	}

	matched := searchVideos(videos, query.Q)
	filtered := filterVideos(matched, func(item Video) bool {
		return matchesVideoFilterKeyword(item, query.Type) &&
			matchesVideoChannel(item, query.Channel, channels) &&
			(query.Year == "" || item.Year == query.Year) &&
			(query.Quality == "" || item.Quality == query.Quality)
	})
	sortName := query.Sort
	if sortName == "relevance" {
		sortName = "default"
	}

	return SearchPageData{
		HotSearchKeywords:    append([]string(nil), defaultHotSearchKeywords...),
		RecommendationVideos: sliceVideos(videos, 1, 5),
		Videos:               sortChannelVideos(filtered, sortName),
	}, nil
}

// WatchPage 获取播放详情页数据；ctx 为请求上下文，videoID 为路由中的视频 id。
func (service *Service) WatchPage(ctx context.Context, videoID string) (WatchPageData, error) {
	videos, err := service.repository.ListVideos(ctx)
	if err != nil {
		return WatchPageData{}, err
	}

	currentVideo := getVideoByID(videos, videoID)
	relatedVideos := getRelatedVideos(videos, currentVideo, 4)

	return WatchPageData{
		RelatedVideos: relatedVideos,
		Video:         currentVideo,
	}, nil
}

// VideoIDs 获取静态路由使用的视频 id 列表；ctx 为请求上下文。
func (service *Service) VideoIDs(ctx context.Context) ([]string, error) {
	videos, err := service.repository.ListVideos(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(videos))
	for _, item := range videos {
		ids = append(ids, item.ID)
	}

	return ids, nil
}

// HealthCheck 验证服务依赖是否可读取；ctx 为请求上下文。
func (service *Service) HealthCheck(ctx context.Context) error {
	_, err := service.repository.ListVideos(ctx)
	return err
}

// loadCatalog 一次性读取频道和视频；ctx 为请求上下文。
func (service *Service) loadCatalog(ctx context.Context) ([]Channel, []Video, error) {
	channels, err := service.repository.ListChannels(ctx)
	if err != nil {
		return nil, nil, err
	}

	videos, err := service.repository.ListVideos(ctx)
	if err != nil {
		return nil, nil, err
	}

	return channels, videos, nil
}

// getChannelBySlug 根据 slug 获取频道；channels 为候选频道，slug 未命中时返回精选频道。
func getChannelBySlug(channels []Channel, slug string) Channel {
	for _, item := range channels {
		if item.Slug == slug {
			return item
		}
	}

	for _, item := range channels {
		if item.Slug == "featured" {
			return item
		}
	}

	if len(channels) == 0 {
		return Channel{Slug: "featured", Label: "精选"}
	}

	return channels[0]
}

// getVideosByChannel 获取频道视频；videos 为候选视频，channel 为频道配置。
func getVideosByChannel(videos []Video, channel Channel) []Video {
	if channel.Slug == "" || channel.Slug == "featured" {
		return append([]Video(nil), videos...)
	}

	filtered := filterVideos(videos, func(item Video) bool {
		return matchesChannelKeywords(item, channel.Keywords)
	})
	if len(filtered) > 0 {
		return filtered
	}

	return append([]Video(nil), videos...)
}

// matchesChannelKeywords 判断视频是否匹配频道关键词；item 为候选视频，keywords 为频道关键词。
func matchesChannelKeywords(item Video, keywords []string) bool {
	searchableText := strings.Join(videoKeywordFields(item), " ")
	for _, keyword := range keywords {
		if strings.Contains(searchableText, keyword) {
			return true
		}
	}

	return len(keywords) == 0
}

// matchesVideoFilterKeyword 判断视频是否命中筛选词；item 为候选视频，keyword 为空时默认命中。
func matchesVideoFilterKeyword(item Video, keyword string) bool {
	normalizedKeyword := strings.ToLower(strings.TrimSpace(keyword))
	if normalizedKeyword == "" {
		return true
	}

	searchableText := strings.ToLower(strings.Join(videoKeywordFields(item), " "))
	return strings.Contains(searchableText, normalizedKeyword)
}

// matchesVideoChannel 判断视频是否匹配频道筛选；item 为候选视频，channelSlug 为空时默认命中。
func matchesVideoChannel(item Video, channelSlug string, channels []Channel) bool {
	if channelSlug == "" || channelSlug == "all" || channelSlug == "featured" {
		return true
	}

	for _, channel := range channels {
		if channel.Slug == channelSlug {
			return matchesChannelKeywords(item, channel.Keywords)
		}
	}

	return true
}

// searchVideos 根据关键词搜索视频；videos 为候选视频，query 为空时返回空列表。
func searchVideos(videos []Video, query string) []Video {
	keyword := strings.ToLower(strings.TrimSpace(query))
	if keyword == "" {
		return []Video{}
	}

	return filterVideos(videos, func(item Video) bool {
		searchableText := strings.ToLower(strings.Join(videoSearchFields(item), " "))
		return strings.Contains(searchableText, keyword)
	})
}

// sortChannelVideos 按频道页排序规则排序；videos 为候选视频，sortName 为排序标识。
func sortChannelVideos(videos []Video, sortName string) []Video {
	sortedVideos := append([]Video(nil), videos...)

	switch sortName {
	case "new":
		sort.SliceStable(sortedVideos, func(first, second int) bool {
			return yearValue(sortedVideos[first]) > yearValue(sortedVideos[second]) ||
				(yearValue(sortedVideos[first]) == yearValue(sortedVideos[second]) &&
					heatValue(sortedVideos[first]) > heatValue(sortedVideos[second]))
		})
	case "hot":
		sort.SliceStable(sortedVideos, func(first, second int) bool {
			return heatValue(sortedVideos[first]) > heatValue(sortedVideos[second])
		})
	case "score":
		sort.SliceStable(sortedVideos, func(first, second int) bool {
			return scoreValue(sortedVideos[first]) > scoreValue(sortedVideos[second])
		})
	}

	return sortedVideos
}

// sortRankVideos 按排行榜规则排序；videos 为候选视频，sortName 为榜单标识。
func sortRankVideos(videos []Video, sortName string) []Video {
	sortedVideos := append([]Video(nil), videos...)
	if sortName == "vip" {
		sortedVideos = filterVideos(sortedVideos, isVIPVideoContent)
	}

	switch sortName {
	case "score":
		sort.SliceStable(sortedVideos, func(first, second int) bool {
			return scoreValue(sortedVideos[first]) > scoreValue(sortedVideos[second])
		})
	case "new", "rising":
		sort.SliceStable(sortedVideos, func(first, second int) bool {
			return yearValue(sortedVideos[first]) > yearValue(sortedVideos[second]) ||
				(yearValue(sortedVideos[first]) == yearValue(sortedVideos[second]) &&
					heatValue(sortedVideos[first]) > heatValue(sortedVideos[second]))
		})
	case "reputation":
		sort.SliceStable(sortedVideos, func(first, second int) bool {
			return scoreValue(sortedVideos[first]) > scoreValue(sortedVideos[second]) ||
				(scoreValue(sortedVideos[first]) == scoreValue(sortedVideos[second]) &&
					heatValue(sortedVideos[first]) > heatValue(sortedVideos[second]))
		})
	default:
		sort.SliceStable(sortedVideos, func(first, second int) bool {
			return heatValue(sortedVideos[first]) > heatValue(sortedVideos[second])
		})
	}

	return sortedVideos
}

// isVIPVideoContent 判断视频是否属于会员权益内容；item 为候选视频。
func isVIPVideoContent(item Video) bool {
	searchableText := strings.Join([]string{
		item.Badge,
		item.Progress,
		item.Quality,
		item.Subtitle,
		strings.Join(item.Tags, " "),
	}, " ")

	return strings.Contains(searchableText, "会员") ||
		strings.Contains(searchableText, "独播") ||
		strings.Contains(searchableText, "4K") ||
		strings.Contains(strings.ToLower(searchableText), "vip")
}

// getVideoByID 根据 id 获取视频；videos 为候选视频，未命中时返回第一个视频。
func getVideoByID(videos []Video, videoID string) Video {
	for _, item := range videos {
		if item.ID == videoID {
			return item
		}
	}

	return firstVideo(videos)
}

// getRelatedVideos 获取相关推荐；videos 为候选视频，currentVideo 为当前视频，limit 为最大数量。
func getRelatedVideos(videos []Video, currentVideo Video, limit int) []Video {
	relatedVideos := make([]Video, 0, limit)
	seenIDs := map[string]bool{currentVideo.ID: true}

	for _, relatedID := range currentVideo.RelatedVideoIDs {
		for _, item := range videos {
			if item.ID == relatedID && !seenIDs[item.ID] {
				relatedVideos = append(relatedVideos, item)
				seenIDs[item.ID] = true
				break
			}
		}

		if len(relatedVideos) >= limit {
			return relatedVideos[:limit]
		}
	}

	for _, item := range videos {
		if !seenIDs[item.ID] {
			relatedVideos = append(relatedVideos, item)
			seenIDs[item.ID] = true
		}

		if len(relatedVideos) >= limit {
			return relatedVideos[:limit]
		}
	}

	return relatedVideos
}

// filterVideos 根据断言过滤视频；videos 为候选列表，keep 返回 true 表示保留。
func filterVideos(videos []Video, keep func(Video) bool) []Video {
	filtered := make([]Video, 0, len(videos))
	for _, item := range videos {
		if keep(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// firstVideo 返回第一个视频；videos 为空时返回零值视频。
func firstVideo(videos []Video) Video {
	if len(videos) == 0 {
		return Video{}
	}

	return videos[0]
}

// sliceVideos 返回视频子区间；videos 为候选列表，start 和 end 为半开区间边界。
func sliceVideos(videos []Video, start int, end int) []Video {
	if start >= len(videos) {
		return []Video{}
	}
	if end > len(videos) {
		end = len(videos)
	}

	return append([]Video(nil), videos[start:end]...)
}

// limitVideos 限制视频数量；videos 为候选列表，limit 为最大数量。
func limitVideos(videos []Video, limit int) []Video {
	if len(videos) <= limit {
		return videos
	}

	return append([]Video(nil), videos[:limit]...)
}

// videoKeywordFields 返回频道筛选字段；item 为待搜索视频。
func videoKeywordFields(item Video) []string {
	fields := []string{
		item.Title,
		item.Subtitle,
		item.Category,
		item.Badge,
		item.Progress,
	}
	fields = append(fields, item.Tags...)

	return fields
}

// videoSearchFields 返回搜索字段；item 为待搜索视频。
func videoSearchFields(item Video) []string {
	fields := []string{
		item.Title,
		item.Subtitle,
		item.Description,
		item.Category,
		item.Badge,
		item.Progress,
		item.Year,
		item.Region,
	}
	fields = append(fields, item.Tags...)
	fields = append(fields, item.CastNames...)

	return fields
}

// heatValue 提取热度数值；item 为待排序视频。
func heatValue(item Video) int {
	digits := strings.Builder{}
	for _, char := range item.Heat {
		if unicode.IsDigit(char) {
			digits.WriteRune(char)
		}
	}

	value, err := strconv.Atoi(digits.String())
	if err != nil {
		return 0
	}

	return value
}

// scoreValue 提取评分数值；item 为待排序视频。
func scoreValue(item Video) float64 {
	value, err := strconv.ParseFloat(item.Score, 64)
	if err != nil {
		return 0
	}

	return value
}

// yearValue 提取年份数值；item 为待排序视频。
func yearValue(item Video) int {
	value, err := strconv.Atoi(item.Year)
	if err != nil {
		return 0
	}

	return value
}

// withDefault 返回非空值或默认值；value 为候选值，fallback 为默认值。
func withDefault(value string, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
