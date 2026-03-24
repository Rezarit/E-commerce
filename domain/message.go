package domain

type OrderMessage struct {
	UserID  int64  `json:"user_id"`
	Address string `json:"address"`
}

type CartMessage struct {
	UserID int64 `json:"user_id"`
}
