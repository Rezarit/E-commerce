package dao

import (
	"github.com/Rezarit/go-seckill-system/internal/domain"
	"gorm.io/gorm"
)

// CreateOrder 创建订单
func CreateOrder(tx *gorm.DB, order *domain.Order) error {
	return InsertRecord(order, tx)
}

// CreateOrderItem 创建订单商品
func CreateOrderItem(tx *gorm.DB, orderItem *domain.OrderItem) error {
	return InsertRecord(orderItem, tx)
}

// GetOrdersByUserID 获取用户订单列表
func GetOrdersByUserID(userID int64) ([]domain.Order, error) {
	var orders []domain.Order
	if err := DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrderByID 根据订单ID获取订单
func GetOrderByID(orderID int64) (*domain.Order, error) {
	var order domain.Order
	if err := GetRecordByField[domain.Order, int64]("order_id", orderID, &order); err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrderItemsByOrderID 根据订单ID获取订单商品
func GetOrderItemsByOrderID(orderID int64) ([]domain.OrderItem, error) {
	var orderItems []domain.OrderItem
	if err := GetRecordsByField[domain.OrderItem, int64]("order_id", orderID, &orderItems); err != nil {
		return nil, err
	}
	return orderItems, nil
}
