package consumers

import (
	"encoding/json"

	"github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/internal/service"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"github.com/Rezarit/go-seckill-system/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

// MaxRetryCount 最大重试次数，超过后进死信队列
const MaxRetryCount = 3

// InitConsumer 初始化消费者
func InitConsumer(name string, handler service.MessageHandler) {
	ch := rabbitmq.GetChannel()
	if ch == nil {
		logger.Sugar.Fatalf("无法获取RabbitMQ通道，%s消费者启动失败", name)
		return
	}

	// 获取队列
	q := rabbitmq.GetQueueName(name)
	if q == "" {
		logger.Sugar.Fatalf("获取%s队列失败，%s消费者启动失败", name, name)
		return
	}

	// 消费消息
	msgs, err := ch.Consume(
		q, // queue
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Sugar.Fatalf("消费订单队列失败: %s", err)
	}

	// 使用一个 goroutine 来处理消息，防止阻塞主线程
	go func() {
		logger.Sugar.Infof("[%s] 消费者已启动，等待消息中...", q)
		for d := range msgs {
			logger.Sugar.Infof("从队列 [%s] 收到一条消息", q)
			// 处理消息
			err = handler(d.Body)
			if err == nil {
				// 处理成功，手动发送 Ack
				if err = d.Ack(false); err != nil {
					logger.Sugar.Errorf("手动发送 Ack 失败: %v", err)
				} else {
					logger.Sugar.Infof("[%s] 消息处理成功，已发送 Ack", q)
				}
				continue
			}

			// 处理失败：重试未超限则重新投递（retry_count+1），超限则 Nack 进死信
			delivery, requeue := handleRetry(ch, q, d)
			if requeue {
				// 重投 retry_count+1 的新消息，Ack 掉当前这条（避免无限重试）
				if err = ch.Publish("", q, false, false, amqp.Publishing{
					ContentType:  "application/json",
					DeliveryMode: amqp.Persistent,
					Body:         delivery,
				}); err != nil {
					logger.Sugar.Errorf("[%s] 重新投递消息失败: %v", q, err)
				} else if err = d.Ack(false); err != nil {
					logger.Sugar.Errorf("[%s] Ack 旧消息失败: %v", q, err)
				} else {
					logger.Sugar.Infof("[%s] 消息处理失败，已重投（retry_count+1）", q)
				}
			} else {
				// 重试超限 → Nack 进死信队列
				if err = d.Nack(false, false); err != nil {
					logger.Sugar.Errorf("手动发送 Nack 失败: %v", err)
				} else {
					logger.Sugar.Errorf("[%s] 消息处理失败且重试超限，已进死信队列: %v", q, err)
				}
			}
		}
	}()
}

// handleRetry 解析消息的 retry_count，决定重投还是进死信。
// 返回重投的消息体（retry_count+1）和是否重投。
func handleRetry(ch *amqp.Channel, queue string, d amqp.Delivery) ([]byte, bool) {
	// 解析消息体，拿到 retry_count（解析失败视为不可重试，直接进死信）
	var msg domain.OrderMessage
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		logger.Sugar.Errorf("[%s] 消息解析失败，直接进死信: %v", queue, err)
		return nil, false
	}

	if msg.RetryCount < MaxRetryCount {
		msg.RetryCount++
		newBody, err := json.Marshal(msg)
		if err != nil {
			logger.Sugar.Errorf("[%s] 重新序列化消息失败: %v", queue, err)
			return nil, false
		}
		return newBody, true
	}

	logger.Sugar.Errorf("[%s] 消息重试已达 %d 次上限，进死信队列 | msg_id: %s", queue, MaxRetryCount, msg.MsgID)
	return nil, false
}
