package rabbitmq

import (
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"github.com/streadway/amqp"
)

var (
	conn           *amqp.Connection
	channel        *amqp.Channel
	declaredQueues map[string]string
)

// InitRabbitMQ 初始化RabbitMQ连接和通道
func InitRabbitMQ(url string, queues map[string]string) error {
	var err error
	// 连接到RabbitMQ服务器
	conn, err = amqp.Dial(url)
	if err != nil {
		logger.Sugar.Errorf("无法连接到RabbitMQ: %v", err)
		return err
	}

	// 创建通道
	channel, err = conn.Channel()
	if err != nil {
		logger.Sugar.Errorf("无法打开通道: %v", err)
		return err
	}

	// 声明队列
	declaredQueues = make(map[string]string)
	// 遍历队列
	for key, queueName := range queues {
		_, err = channel.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			logger.Sugar.Errorf("无法声明队列 %s: %v", queueName, err)
			return err
		}
		declaredQueues[key] = queueName
	}

	logger.Sugar.Info("RabbitMQ 初始化成功，所有队列已声明！")
	return nil
}

// GetQueueName 根据 key 获取真实队列名
func GetQueueName(key string) string {
	return declaredQueues[key]
}

// GetChannel 返回当前通道
func GetChannel() *amqp.Channel {
	return channel
}

// Close 关闭RabbitMQ连接和通道
func Close() {
	logger.Sugar.Info("正在关闭RabbitMQ连接...")
	if channel != nil {
		if err := channel.Close(); err != nil {
			logger.Sugar.Errorf("关闭RabbitMQ通道失败: %v", err)
		}
	}
	if conn != nil {
		if err := conn.Close(); err != nil {
			logger.Sugar.Errorf("关闭RabbitMQ连接失败: %v", err)
		}
	}
	logger.Sugar.Info("RabbitMQ连接已成功关闭。")
}
