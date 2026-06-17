package video

type Channel struct {
	Slug        string   `json:"slug"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Accent      string   `json:"accent"`
}

type Episode struct {
	Episode  int    `json:"episode"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
	Status   string `json:"status,omitempty"`
}

type ReleaseCalendarItem struct {
	Time   string `json:"time"`
	Detail string `json:"detail"`
	Active bool   `json:"active,omitempty"`
}

type Video struct {
	ID              string                `json:"id"`
	Title           string                `json:"title"`
	Subtitle        string                `json:"subtitle"`
	Description     string                `json:"description"`
	Score           string                `json:"score"`
	Heat            string                `json:"heat"`
	Update          string                `json:"update"`
	Category        string                `json:"category"`
	Year            string                `json:"year"`
	Region          string                `json:"region"`
	TotalEpisodes   int                   `json:"totalEpisodes"`
	Quality         string                `json:"quality"`
	Badge           string                `json:"badge"`
	Progress        string                `json:"progress"`
	Duration        string                `json:"duration"`
	SourceURL       string                `json:"sourceUrl"`
	CoverGradient   string                `json:"coverGradient"`
	Tags            []string              `json:"tags"`
	CastNames       []string              `json:"castNames"`
	ReleaseCalendar []ReleaseCalendarItem `json:"releaseCalendar"`
	Episodes        []Episode             `json:"episodes"`
	RelatedVideoIDs []string              `json:"relatedVideoIds"`
}

type HomePageData struct {
	FeaturedVideo        Video   `json:"featuredVideo"`
	HotVideos            []Video `json:"hotVideos"`
	RankVideos           []Video `json:"rankVideos"`
	RecommendationVideos []Video `json:"recommendationVideos"`
}

type ChannelPageData struct {
	Channel   Channel `json:"channel"`
	HeroVideo Video   `json:"heroVideo"`
	Videos    []Video `json:"videos"`
}

type SearchPageData struct {
	HotSearchKeywords    []string `json:"hotSearchKeywords"`
	RecommendationVideos []Video  `json:"recommendationVideos"`
	Videos               []Video  `json:"videos"`
}

type WatchPageData struct {
	RelatedVideos []Video `json:"relatedVideos"`
	Video         Video   `json:"video"`
}

type RankQuery struct {
	Channel string
	Sort    string
}

type ChannelQuery struct {
	Slug string
	Type string
	Year string
	Sort string
}

type SearchQuery struct {
	Q       string
	Channel string
	Quality string
	Type    string
	Year    string
	Sort    string
}
