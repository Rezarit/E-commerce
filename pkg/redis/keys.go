package redis

import (
	"fmt"
	"time"
)

const (
	// 商品相关
	KeyProductDetail = "cache:product:detail:%d" // 商品详情缓存
	KeyProductNull   = "cache:product:null:%d"   // 商品空值缓存
	KeySeckillStock  = "seckill:stock:%d"        // 秒杀库存
	KeyOrderResult   = "user:order:%d:%d"        // 订单结果缓存

	// 购物车相关
	KeyCart = "cart:%d" // 用户购物车，使用Hash结构存储
)

// 默认过期时间常量
const (
	DefaultProductCacheTTL = 1 * time.Hour
	DefaultNullCacheTTL    = 5 * time.Minute
	DefaultSessionTTL      = 24 * time.Hour
	DefaultSeckillLockTTL  = 10 * time.Second
)

// BuildKey 构建Key的工具函数
func BuildKey(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
