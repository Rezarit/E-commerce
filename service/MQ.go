package service

import (
	"encoding/json"
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/Rezarit/go-seckill-system/pkg/rabbitmq"
	"github.com/streadway/amqp"
	"log"
)

// SendMessage 将信息打包发送到MQ
func SendMessage[T any](msg T, queueName string) error {
	// 将消息序列化成JSON
	msgBody, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[Service] 信息序列化失败: %v", err)
		return &domain.BusinessError{
			Code: domain.ErrCodeSystemError,
			Msg:  "信息序列化失败，请稍后再试",
		}
	}

	// 获取RabbitMQ通道
	ch := rabbitmq.GetChannel()
	if ch == nil {
		log.Println("[Service] 无法获取RabbitMQ通道，请检查MQ连接")
		return &domain.BusinessError{
			Code: domain.ErrCodeSystemError,
			Msg:  "服务繁忙，请稍后再试",
		}
	}

	// 发布消息到队列
	err = ch.Publish(
		"",
		rabbitmq.GetQueueName(queueName),
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // 持久化
			Body:         msgBody,
		},
	)

	if err != nil {
		log.Printf("[Service] 发布信息到MQ失败: %v", err)
		return &domain.BusinessError{
			Code: domain.ErrCodeMQPublishError,
			Msg:  "信息发布失败，请稍后再试",
		}
	}
	return nil
}
