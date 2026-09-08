package redis

import (
	"fmt"
	"time"
)

const (
	// 商品相关
	KeyProductDetail = "cache:product:detail:%d" // 商品详情缓存
	KeyProductNull   = "cache:product:null:%d"   // 商品空值缓存
	KeyProductLock   = "lock:product:%d"         // 商品详情缓存重建互斥锁（防击穿）
	KeySeckillStock  = "seckill:stock:%d"        // 秒杀库存
	KeyOrderResult   = "user:order:result:%s"    // 订单处理结果缓存（键=msg_id，值为三态 JSON）
	KeyDeadProcessed = "dead:processed:%s"       // 死信处理幂等标记（键=msg_id：库存释放已完成）

	// 购物车相关
	KeyCart     = "cart:%d"      // 用户购物车，使用Hash结构存储
	KeyCartNull = "cart:null:%d" // 购物车空值缓存（防穿透）
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
