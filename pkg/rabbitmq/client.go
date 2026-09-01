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

// 死信交换机名称
const dlxExchangeName = "dlx_exchange"

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

	// 开启 publisher confirm：发布消息后等 MQ 确认，确保消息真正到达
	if err = channel.Confirm(false); err != nil {
		logger.Sugar.Errorf("开启 publisher confirm 失败: %v", err)
		return err
	}

	// 声明死信交换机（fanout：死信消息广播给所有绑定的死信队列）
	if err = channel.ExchangeDeclare(dlxExchangeName, "fanout", true, false, false, false, nil); err != nil {
		logger.Sugar.Errorf("无法声明死信交换机 %s: %v", dlxExchangeName, err)
		return err
	}

	// 声明队列
	declaredQueues = make(map[string]string)
	// 遍历队列
	for key, queueName := range queues {
		var args amqp.Table

		// order 主队列绑定死信交换机：消息变死信后自动路由到死信队列
		if key == "order" {
			args = amqp.Table{
				"x-dead-letter-exchange": dlxExchangeName,
			}
		}

		_, err = channel.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			args,
		)
		if err != nil {
			logger.Sugar.Errorf("无法声明队列 %s: %v", queueName, err)
			return err
		}

		// 死信队列绑定到死信交换机
		if key == "order_dlx" {
			if err = channel.QueueBind(queueName, "", dlxExchangeName, false, nil); err != nil {
				logger.Sugar.Errorf("无法绑定死信队列 %s 到死信交换机: %v", queueName, err)
				return err
			}
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
