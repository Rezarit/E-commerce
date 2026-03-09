package domain

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
