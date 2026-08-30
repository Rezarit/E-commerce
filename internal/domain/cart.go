package domain

// CartItem 购物车项：一个用户购物车里的一个商品
type CartItem struct {
	CartID    int64 `json:"id" gorm:"primaryKey;autoIncrement"` // 主键
	UserID    int64 `json:"user_id" gorm:"index"`               // 用户ID
	ProductID int64 `json:"product_id" gorm:"index"`            // 商品ID
	Quantity  int   `json:"quantity" gorm:"default:1"`          // 商品数量
}

// TableName 指定表名为 carts（与 init.sql 保持一致，避免 GORM 默认推导成 cart_items）
func (CartItem) TableName() string {
	return "carts"
}

type AddToCartRequest struct {
	Quantity int `json:"quantity"`
}
