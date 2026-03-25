package consumers

import (
	"github.com/Rezarit/go-seckill-system/pkg/rabbitmq"
	"github.com/Rezarit/go-seckill-system/service"
	"log"
)

// InitConsumer 初始化消费者
func InitConsumer(name string, handler service.MessageHandler) {
	ch := rabbitmq.GetChannel()
	if ch == nil {
		log.Fatalf("无法获取RabbitMQ通道，%s消费者启动失败", name)
		return
	}

	// 获取队列
	q := rabbitmq.GetQueueName(name)
	if q == "" {
		log.Fatalf("获取%s队列失败，%s消费者启动失败", name, name)
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
		log.Fatalf("消费订单队列失败: %s", err)
	}

	// 使用一个 goroutine 来处理消息，防止阻塞主线程
	go func() {
		log.Printf("[%s] 消费者已启动，等待消息中...", q)
		for d := range msgs {
			log.Printf("从队列 [%s] 收到一条消息", q)
			// 处理消息
			err = handler(d.Body)
			if err == nil {
				// 手动发送 Ack
				err = d.Ack(false)
				if err != nil {
					log.Printf("手动发送 Ack 失败: %v", err)
				} else {
					log.Printf("[%s] 消息处理成功，已发送 Ack", q)
				}
			} else {
				// 处理失败，发送 Nack
				err = d.Nack(false, false)
				if err != nil {
					log.Printf("手动发送 Nack 失败: %v", err)
				} else {
					log.Printf("[%s] 消息处理失败: %v，已发送 Nack，消息将被丢弃或死信", q, err)
				}
			}
		}
	}()
}
