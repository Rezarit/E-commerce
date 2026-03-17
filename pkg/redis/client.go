package redis

import (
	"context"
	"fmt"
	"github.com/Rezarit/go-seckill-system/pkg/config"
	"github.com/go-redis/redis/v8"
	"log"
)

var client *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.RedisConfig) error {
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis连接失败: %v", err)
	}

	log.Printf("redis连接成功: %s", cfg.Addr)
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return client
}
