package consumers

import (
	"encoding/json"
	"github.com/Rezarit/go-seckill-system/internal/dao"
	"github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/internal/service"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
)

// InitCartConsumer 初始化购物车消费者
func InitCartConsumer() {
	InitConsumer("cart", handleCartMessage)
}

// handleCartMessage 处理购物车消息的函数
func handleCartMessage(body []byte) error {
	// 解析消息
	var msg domain.CartMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logger.Sugar.Errorf("解析购物车消息失败: %v", err)
		return err
	}
	// 清空购物车：先删 MySQL（持久化），再删 Redis（缓存）
	if err := dao.ClearCart(msg.UserID); err != nil {
		logger.Sugar.Errorf("清空 MySQL 购物车失败: %v", err)
		return err
	}
	if err := service.ClearCartInRedis(msg.UserID); err != nil {
		logger.Sugar.Errorf("清空 Redis 购物车失败: %v", err)
		return err
	}
	return nil
}
