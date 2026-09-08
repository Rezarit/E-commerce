package consumers

import (
	"encoding/json"
	"time"

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

	// 直接用消息里入口已扣减的商品清单（不再读购物车：购物车已在入口被读，且库存已扣）
	items := msg.Items
	if len(items) == 0 {
		logger.Sugar.Infof("消费者发现商品清单为空，用户ID: %d，消息将被丢弃", msg.UserID)
		return nil
	}

	// 执行数据库下单操作（只建订单 + 扣 MySQL 库存，Redis 库存已在 API 入口扣过）
	orderID, err := service2.ExecuteOrderCreation(msg.MsgID, msg.UserID, msg.Address, items)
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

	// 建单成功（orderID>0）才写「成交回执」到 Redis result：MySQL 先落库、Redis 后写，
	// 查 result 不会误报成功。orderID==0 = 幂等命中（订单此前已建过），跳过写，
	// 避免把已存 result 的 order_id 覆盖成 0。
	if orderID > 0 {
		err = service2.Order.SetOrderResult(msg.MsgID, domain.OrderResult{
			Status:    domain.OrderResultSuccess,
			OrderID:   orderID,
			UserID:    msg.UserID,
			CreatedAt: time.Now(),
		}, redis.DefaultSessionTTL)
		if err != nil {
			// 回执写失败：return err 触发消息重投，幂等命中后重跑会补写（不会重复建单）
			logger.Sugar.Errorf("写入订单结果失败 | msg_id: %s | 错误: %v", msg.MsgID, err)
			return err
		}
	}
	logger.Sugar.Infof("订单创建成功！OrderID: %d。结果已写入 Redis result", orderID)

	return nil
}
