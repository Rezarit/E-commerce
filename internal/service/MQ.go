package service

import (
	"encoding/json"
	"time"

	"github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"github.com/Rezarit/go-seckill-system/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

// SendMessage 将信息打包发送到MQ
func SendMessage[T any](msg T, queueName string) error {
	// 将消息序列化成JSON
	msgBody, err := json.Marshal(msg)
	if err != nil {
		logger.Sugar.Errorf("[Service] 信息序列化失败: %v", err)
		return &domain.BusinessError{
			Code: domain.ErrCodeSystemError,
			Msg:  "信息序列化失败，请稍后再试",
		}
	}

	// 获取RabbitMQ通道
	ch := rabbitmq.GetChannel()
	if ch == nil {
		logger.Sugar.Info("[Service] 无法获取RabbitMQ通道，请检查MQ连接")
		return &domain.BusinessError{
			Code: domain.ErrCodeSystemError,
			Msg:  "服务繁忙，请稍后再试",
		}
	}

	// 注册确认监听（publisher confirm：确保消息真正到达 MQ）
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

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
		logger.Sugar.Errorf("[Service] 发布信息到MQ失败: %v", err)
		return &domain.BusinessError{
			Code: domain.ErrCodeMQPublishError,
			Msg:  "信息发布失败，请稍后再试",
		}
	}

	// 等 MQ 确认（带超时），超时/未确认视为发布失败
	select {
	case confirmed := <-confirms:
		if !confirmed.Ack {
			logger.Sugar.Errorf("[Service] 消息被 MQ 拒绝（nack）")
			return &domain.BusinessError{
				Code: domain.ErrCodeMQPublishError,
				Msg:  "信息发布失败，请稍后再试",
			}
		}
	case <-time.After(5 * time.Second):
		logger.Sugar.Errorf("[Service] 等待 MQ 确认超时，消息可能未到达")
		return &domain.BusinessError{
			Code: domain.ErrCodeMQPublishError,
			Msg:  "信息发布失败，请稍后再试",
		}
	}

	return nil
}
