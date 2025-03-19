package domain

type User struct {
	UserID       int    `json:"user_id"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Phone        int    `json:"phone"`
	QQ           int    `json:"qq"`
	Gender       string `json:"gender"`
	Email        string `json:"email"`
	Birthday     string `json:"birthday"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type Product struct {
	ProductID   int    `json:"product_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CommentNum  string `json:"comment_num"`
	Price       string `json:"price"`
	IsaddedCart bool   `json:"is_added_cart"`
	Cover       string `json:"cover"`
	PublishTime int    `json:"publish_time"`
	Link        string `json:"link"`
}

type Order struct {
	OrderID int     `json:"order_id"`
	Address string  `json:"address"`
	Total   float32 `json:"total"`
	UserID  string  `json:"user_id"`
}

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

var ErrorStatus int

var ErrorMsg string
