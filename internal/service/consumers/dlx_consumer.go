package consumers

import (
	"github.com/Rezarit/go-seckill-system/pkg/logger"
)

// InitDLXConsumer 初始化死信队列消费者（只告警 + 留痕，不做业务、不重试）
func InitDLXConsumer() {
	InitConsumer("order_dlx", handleDLXMessage)
}

// handleDLXMessage 死信消息处理：只打告警日志，人工介入
func handleDLXMessage(body []byte) error {
	logger.Sugar.Errorf("⚠️⚠️⚠️ [死信队列] 消息重试超限进死信，需要人工介入！消息体: %s", string(body))
	// 返回 nil → Ack，不再重试（死信是终点，不无限循环）
	return nil
}
