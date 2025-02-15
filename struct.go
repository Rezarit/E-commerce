package main

type User struct {
	ID           string `json:"id"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Phone        int64  `json:"phone"`
	QQ           int64  `json:"qq"`
	Gender       string `json:"gender"`
	Email        string `json:"email"`
	Birthday     string `json:"birthday"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type Product struct {
	ID          string `json:"product_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CommentNum  string `json:"comment_num"`
	Price       string `json:"price"`
	IsaddedCart bool   `json:"is_added_cart"`
	Cover       string `json:"cover"`
	PublishTime int64  `json:"publish_time"`
	Link        string `json:"link"`
}

type Order struct {
	OrderID string  `json:"order_id"`
	Address string  `json:"address"`
	Total   float32 `json:"total"`
	UserID  string  `json:"user_id"`
}

type Comment struct {
	PostID      string `json:"post_id"`
	PublishTime int64  `json:"publish_time"`
	Content     string `json:"content"`
	UserID      string `json:"user_id"`
	Avatar      string `json:"avatar"`
	Nickname    string `json:"nickname"`
	PraiseCount string `json:"praise_count"`
	IsPraised   int    `json:"is_praised"`
	ProductID   string `json:"product_id"`
}
