package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

// 订单状态常量
const (
	OrderStatusPending   = "pending"   // 待支付
	OrderStatusPaid      = "paid"      // 已支付
	OrderStatusShipping  = "shipping"  // 已发货
	OrderStatusCompleted = "completed" // 已完成
	OrderStatusCancelled = "cancelled" // 已取消
)

type Order struct {
	OrderID   int64           `json:"order_id" gorm:"primaryKey;autoIncrement"`
	UserID    int64           `json:"user_id" gorm:"index"`
	Address   string          `json:"address" gorm:"not null"`
	Total     decimal.Decimal `json:"total" gorm:"type:decimal(10,2);default:0"`
	Status    string          `json:"status" gorm:"default:'pending'"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
}

type OrderItem struct {
	OrderItemID int64           `json:"order_item_id" gorm:"primaryKey;autoIncrement"`
	OrderID     int64           `json:"order_id" gorm:"index"`
	ProductID   int64           `json:"product_id"`
	ProductName string          `json:"product_name" gorm:"not null"`
	Quantity    int             `json:"quantity" gorm:"not null"`
	Price       decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not null"`
}

type OrderCreateRequest struct {
	Address string `json:"address" binding:"required"`
}