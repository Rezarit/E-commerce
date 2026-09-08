package consumers

import (
	"encoding/json"
	"time"

	"github.com/Rezarit/go-seckill-system/internal/domain"
	service2 "github.com/Rezarit/go-seckill-system/internal/service"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"github.com/Rezarit/go-seckill-system/pkg/redis"
)

// InitDLXConsumer 初始化死信队列消费者
// 职责：① 告警留痕 ② 刷新式补偿 Redis 库存（让 Redis 对齐 MySQL，防超卖）
func InitDLXConsumer() {
	InitConsumer("order_dlx", handleDLXMessage)
}

// handleDLXMessage 死信消息处理（失败终点收尾三件套：精确补偿库存 + 写 failed 态 + 告警）
// 死信 = 该单最终没建成 = MySQL 库存没错 = 只需把入口扣掉的 Redis 份额精确还回可卖池。
// 用精确 AddBackStock（按死信单 items），不用刷新式整 key 覆盖——避免误放在途单份额（笔记 10.9 修订）。
func handleDLXMessage(body []byte) error {
	// 解析消息，拿到商品清单（items 里的商品入口已扣过 Redis，需补偿）
	var msg domain.OrderMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logger.Sugar.Errorf("[死信] 消息解析失败，无法补偿: %v", err)
		return nil // 返回 nil → Ack，不再重试（死信是终点）
	}

	// ① 精确补偿：先抢「库存释放已完成」幂等标记，防 MQ 重投/重复消费导致重复释放（反向超卖）
	claimed, err := service2.Order.ClaimDeadProcessed(msg.MsgID)
	if err != nil {
		logger.Sugar.Errorf("[死信] 抢幂等标记失败 | msg_id: %s | 错误: %v（Ack，靠 TTL/对账兜）", msg.MsgID, err)
	} else if !claimed {
		// 已释放过：跳过释放。failed 态仍会走下面②幂等补写（SET 覆盖可重复）
		logger.Sugar.Infof("[死信] 该单库存已释放过，跳过释放 | msg_id: %s", msg.MsgID)
	} else if relErr := service2.ReleaseDeadOrderStock(msg.Items); relErr != nil {
		// 释放没成功：撤销标记，让可能到来的 MQ 重投能再补救（标记 ≠ 抢到，是「干完了」）
		logger.Sugar.Errorf("[死信] 库存释放失败 | msg_id: %s | 错误: %v（撤销标记，靠重投/TTL 补救）", msg.MsgID, relErr)
		_ = service2.Order.ReleaseDeadProcessed(msg.MsgID)
	}

	// ② 写最终失败态到 Redis result，让前端/用户轮询到明确失败。
	// 写失败也先 Ack，不套娃；极端情况用户看到 processing 直到 TTL 过期，由 MySQL 对账兜底。
	if err := service2.Order.SetOrderResult(msg.MsgID, domain.OrderResult{
		Status:    domain.OrderResultFailed,
		UserID:    msg.UserID,
		Reason:    "订单处理失败（重试耗尽，已进死信）",
		CreatedAt: time.Now(),
	}, redis.DefaultSessionTTL); err != nil {
		logger.Sugar.Errorf("[死信] 写入失败态失败 | msg_id: %s | 错误: %v", msg.MsgID, err)
	}

	// ③ 告警留痕，人工介入
	logger.Sugar.Errorf("⚠️⚠️⚠️ [死信队列] 消息重试超限进死信，需要人工介入！msg_id: %s | 消息体: %s", msg.MsgID, string(body))

	return nil
}
