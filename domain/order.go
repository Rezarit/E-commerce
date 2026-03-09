package domain

type Order struct {
	OrderID int     `json:"order_id"`
	Address string  `json:"address"`
	Total   float32 `json:"total"`
	UserID  string  `json:"user_id"`
}
