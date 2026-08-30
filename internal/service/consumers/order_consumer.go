package consumers

import (
	"encoding/json"

	"github.com/Rezarit/go-seckill-system/internal/domain"
	service2 "github.com/Rezarit/go-seckill-system/internal/service"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"github.com/Rezarit/go-seckill-system/pkg/redis"
)

// InitOrderConsumer 初始化订单消费者
func InitOrderConsumer() {
	InitConsumer("order", handleOrderMessage)
}

// handleOrderMessage 处理订单消息的函数
func handleOrderMessage(body []byte) error {
	// 解析消息
	var msg domain.OrderMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logger.Sugar.Errorf("解析订单消息失败: %v", err)
		return err
	}

	// 直接用消息里入口已扣减的商品清单（不再读购物车）
	items := msg.Items
	if len(items) == 0 {
		logger.Sugar.Infof("消费者发现商品清单为空，用户ID: %d，消息将被丢弃", msg.UserID)
		return nil
	}

	// 执行数据库下单操作（只建订单 + 扣 MySQL 库存）
	orderID, err := service2.ExecuteOrderCreation(msg.UserID, msg.Address, items)
	if err != nil {
		logger.Sugar.Errorf("消费者创建订单失败: %v", err)
		return err
	}

	// 创建一个购物车消息
	cartMsg := domain.CartMessage{UserID: msg.UserID}

	// 发送清空购物车的消息到 MQ
	err = service2.SendMessage(cartMsg, "cart")
	if err != nil {
		logger.Sugar.Errorf("发送清空购物车消息失败: %v", err)
	}
	logger.Sugar.Infof("已发送清空购物车消息 | 用户ID: %d", msg.UserID)

	logger.Sugar.Infof("订单创建成功！OrderID: %d。准备清空用户 %d 的购物车...", orderID, msg.UserID)

	// 结果写入 Redis，供前端轮询
	err = service2.Order.OrderResult(orderID, msg.UserID, redis.DefaultSessionTTL)
	if err != nil {
		return err
	}
	logger.Sugar.Infof("准备将订单结果写入 Redis...")

	return nil
}
