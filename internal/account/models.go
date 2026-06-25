package account

type UserProfile struct {
	ID         string `json:"id"`
	AvatarURL  string `json:"avatarUrl"`
	Email      string `json:"email"`
	IsLoggedIn bool   `json:"isLoggedIn"`
	IsVip      bool   `json:"isVip"`
	Nickname   string `json:"nickname"`
	Phone      string `json:"phone"`
	VipUntil   string `json:"vipUntil"`
}

type LoginInput struct {
	AvatarURL string `json:"avatarUrl"`
	Contact   string `json:"contact"`
	Nickname  string `json:"nickname"`
}

type VipInput struct {
	VipUntil string `json:"vipUntil"`
}

type FavoriteInput struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Category      string `json:"category"`
	Progress      string `json:"progress"`
	CoverGradient string `json:"coverGradient"`
	Description   string `json:"description"`
}

type FavoriteItem struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Category      string `json:"category"`
	Progress      string `json:"progress"`
	CoverGradient string `json:"coverGradient"`
	Description   string `json:"description"`
	AddedAt       int64  `json:"addedAt"`
}

type WatchHistoryInput struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Category        string `json:"category"`
	Progress        string `json:"progress"`
	CoverGradient   string `json:"coverGradient"`
	Episode         *int   `json:"episode,omitempty"`
	WatchSeconds    *int   `json:"watchSeconds,omitempty"`
	DurationSeconds *int   `json:"durationSeconds,omitempty"`
}

type WatchHistoryItem struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Category        string `json:"category"`
	Progress        string `json:"progress"`
	CoverGradient   string `json:"coverGradient"`
	Episode         *int   `json:"episode,omitempty"`
	WatchSeconds    *int   `json:"watchSeconds,omitempty"`
	DurationSeconds *int   `json:"durationSeconds,omitempty"`
	WatchedAt       int64  `json:"watchedAt"`
}
