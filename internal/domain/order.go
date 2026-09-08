package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// 订单状态常量
const (
	OrderStatusCreated    = "created"    // 已创建（成交）：无支付的项目里下单即成交的终态
	OrderStatusProcessing = "processing" // 处理中
	OrderStatusPending    = "pending"    // 待支付（接入支付后的下一位，当前不用）
	OrderStatusPaid       = "paid"       // 已支付
	OrderStatusShipping   = "shipping"   // 已发货
	OrderStatusCompleted  = "completed"  // 已完成
	OrderStatusCancelled  = "cancelled"  // 已取消
)

type Order struct {
	OrderID   int64           `json:"order_id" gorm:"primaryKey;autoIncrement"`
	MsgID     string          `json:"msg_id" gorm:"uniqueIndex"` // 消息ID，用于幂等（一条消息只对应一个订单）
	UserID    int64           `json:"user_id" gorm:"index"`
	Address   string          `json:"address" gorm:"not null"`
	Total     decimal.Decimal `json:"total" gorm:"type:decimal(10,2);default:0"`
	Status    string          `json:"status" gorm:"default:'pending'"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
}

// 下单请求的「处理结果」状态（Redis result 用，独立于订单本身的生命周期状态）
const (
	OrderResultProcessing = "processing" // 查无 result 键 = 该态（受理后尚未处理完）
	OrderResultSuccess    = "success"    // 建单成交（MySQL 已落库）
	OrderResultFailed     = "failed"     // 最终失败（重试耗尽进死信）
)

// OrderResult 订单处理结果（Redis result 的值，键 = msg_id，贯穿 受理→查询）
type OrderResult struct {
	Status    string    `json:"status"`
	OrderID   int64     `json:"order_id,omitempty"`
	UserID    int64     `json:"user_id"`
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
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

// OrderItemSnapshot 下单受理时返回的「应下单商品」快照（订单建成后可用于对账）
type OrderItemSnapshot struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// OrderAcceptance 下单受理凭证：前端/压测拿 msg_id 去轮询 /order/result
// msg_id 同时是幂等键（唯一索引）与对外受理单号，贯穿 受理→消息→落库→查询
type OrderAcceptance struct {
	MsgID string              `json:"msg_id"`
	Items []OrderItemSnapshot `json:"items"`
}
