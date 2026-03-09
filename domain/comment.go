package domain

type Comment struct {
	CommentID   int    `json:"post_id"`
	PublishTime int    `json:"publish_time"`
	Content     string `json:"content"`
	UserID      int    `json:"user_id"`
	Avatar      string `json:"avatar"`
	Nickname    string `json:"nickname"`
	PraiseCount string `json:"praise_count"`
	IsPraised   int    `json:"is_praised"`
	ProductID   string `json:"product_id"`
}
