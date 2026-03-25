package consumers

import (
	"encoding/json"
	"github.com/Rezarit/go-seckill-system/internal/dao"
	"github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/internal/service"
	"log"
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
		log.Printf("解析购物车消息失败: %v", err)
		return err
	}
	// 清空购物车
	err := service.ClearCartInRedis(msg.UserID) // 清空 Redis 购物车
	if err != nil {
		return err
	}
	err = dao.ClearCart(msg.UserID) // 清空 MySQL 购物车
	if err != nil {
		return err
	}
	return nil
}
