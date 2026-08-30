package domain

// OrderItemMsg 订单消息里的商品项（入口已扣减的商品清单）
type OrderItemMsg struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type OrderMessage struct {
	UserID  int64          `json:"user_id"`
	Address string         `json:"address"`
	Items   []OrderItemMsg `json:"items"` // 已扣减的商品清单
}

type CartMessage struct {
	UserID int64 `json:"user_id"`
}
