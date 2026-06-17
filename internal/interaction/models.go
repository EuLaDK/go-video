package interaction

type CommentSort string

const (
	CommentSortLatest CommentSort = "latest"
	CommentSortHot    CommentSort = "hot"
)

type CommentInput struct {
	Content string `json:"content"`
}

type CommentItem struct {
	ID        string `json:"id"`
	VideoID   string `json:"videoId"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	LikedByMe bool   `json:"likedByMe"`
	Likes     int    `json:"likes"`
	CreatedAt int64  `json:"createdAt"`
}

type DanmakuInput struct {
	Content string `json:"content"`
	Color   string `json:"color"`
}

type DanmakuItem struct {
	ID        string `json:"id"`
	VideoID   string `json:"videoId"`
	Content   string `json:"content"`
	Color     string `json:"color"`
	CreatedAt int64  `json:"createdAt"`
}
