package consumers

import (
	"encoding/json"

	"github.com/Rezarit/go-seckill-system/internal/domain"
	service2 "github.com/Rezarit/go-seckill-system/internal/service"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
)

// InitDLXConsumer 初始化死信队列消费者
// 职责：① 告警留痕 ② 刷新式补偿 Redis 库存（让 Redis 对齐 MySQL，防超卖）
func InitDLXConsumer() {
	InitConsumer("order_dlx", handleDLXMessage)
}

// handleDLXMessage 死信消息处理
// 刷新式补偿：死信 = 订单没建成 = MySQL 库存正确，逐个商品把 MySQL 库存覆盖回 Redis。
// 以 MySQL 为权威，天然幂等、不信任消息内容（假消息也只是把 Redis 设为 MySQL 值，不会超卖）。
func handleDLXMessage(body []byte) error {
	logger.Sugar.Errorf("⚠️⚠️⚠️ [死信队列] 消息重试超限进死信，需要人工介入！消息体: %s", string(body))

	// 解析消息，拿到商品清单（items 里的商品入口已扣过 Redis，需补偿）
	var msg domain.OrderMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logger.Sugar.Errorf("[死信] 消息解析失败，无法补偿: %v", err)
		return nil // 返回 nil → Ack，不再重试（死信是终点）
	}

	// 刷新式补偿：逐个商品，查 MySQL 库存并覆盖写回 Redis
	for _, item := range msg.Items {
		if err := service2.RefreshStockFromDB(item.ProductID); err != nil {
			// 补偿失败（MySQL/Redis 异常）：打日志，靠 Redis TTL 过期后从 MySQL 重新预热兜底
			logger.Sugar.Errorf("[死信] 库存补偿失败 | 商品ID: %d | 错误: %v（将靠 TTL 兜底）", item.ProductID, err)
		}
	}

	return nil
}
