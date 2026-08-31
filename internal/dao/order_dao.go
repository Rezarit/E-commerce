package dao

import (
	"errors"
	"strings"

	"github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"gorm.io/gorm"
)

// ErrDuplicateMsgID 消息重复（幂等命中）：同一条消息已被处理过
var ErrDuplicateMsgID = errors.New("duplicate msg_id")

// CreateOrder 创建订单（靠数据库唯一索引做幂等：
// 同一条消息重复插入时撞 msg_id 唯一索引，返回 ErrDuplicateMsgID。
// 用索引而非"先查再插"是因为后者不是原子的，并发下两个消费者可能同时查到都不存在）
func CreateOrder(tx *gorm.DB, order *domain.Order) error {
	err := InsertRecord(order, tx)
	if err != nil && isDuplicateError(err) {
		logger.Sugar.Infof("[DAO] 消息重复（msg_id 已存在），订单创建被跳过 | msg_id: %s", order.MsgID)
		return ErrDuplicateMsgID
	}
	return err
}

// isDuplicateError 判断是否为 MySQL 唯一键冲突错误（1062 是 MySQL 唯一键冲突错误码）
func isDuplicateError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "1062"))
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
