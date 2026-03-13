package domain

type Cart struct {
	CartID    int64 `json:"id" gorm:"primaryKey;autoIncrement"` // 主键
	UserID    int64 `json:"user_id" gorm:"index"`               // 用户ID
	ProductID int64 `json:"product_id" gorm:"index"`            // 商品ID
	Quantity  int   `json:"quantity" gorm:"default:1"`          // 商品数量
}

type AddToCartRequest struct {
	Quantity int `json:"quantity"`
}
