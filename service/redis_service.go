package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rezarit/go-seckill-system/domain"
	myredis "github.com/Rezarit/go-seckill-system/pkg/redis"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

// CacheNullProduct 缓存空商品到Redis
func (s *CacheService) CacheNullProduct(productID int64, expiration time.Duration) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyProductNull, productID)

	err := s.client.Set(ctx, key, productID, expiration).Err()
	if err != nil {
		return &domain.BusinessError{
			Code: domain.ErrCodeCacheError,
			Msg:  "缓存空商品失败: " + err.Error(),
		}
	}
	return nil
}

// CacheProduct 缓存商品详情
func (s *CacheService) CacheProduct(product *domain.Product, expiration time.Duration) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyProductDetail, product.ProductID)

	data, err := json.Marshal(product)
	if err != nil {
		return &domain.BusinessError{
			Code: domain.ErrCodeCacheSerializeError,
			Msg:  "商品序列化失败: " + err.Error(),
		}
	}

	err = s.client.Set(ctx, key, string(data), expiration).Err()
	if err != nil {
		return &domain.BusinessError{
			Code: domain.ErrCodeCacheError,
			Msg:  "缓存商品失败: " + err.Error(),
		}
	}
	return nil
}

// CacheProductStock 缓存商品库存（库存减扣专用）
func (s *CacheService) CacheProductStock(productID int64, stock int32, expiration time.Duration) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeySeckillStock, productID)

	err := s.client.Set(ctx, key, stock, expiration).Err()
	if err != nil {
		return &domain.BusinessError{
			Code: domain.ErrCodeCacheError,
			Msg:  "缓存商品库存失败: " + err.Error(),
		}
	}
	return nil
}

// BatchCacheProductStocks 批量缓存商品库存（使用Pipeline）
func (s *CacheService) BatchCacheProductStocks(stocks map[int64]int32, expiration time.Duration) error {
	ctx := context.Background()

	pipe := s.client.Pipeline()
	for productID, stock := range stocks {
		key := myredis.BuildKey(myredis.KeySeckillStock, productID)
		err := pipe.Set(ctx, key, stock, expiration).Err()
		if err != nil {
			return &domain.BusinessError{
				Code: domain.ErrCodeCacheError,
				Msg:  fmt.Sprintf("批量缓存商品库存失败 | 商品ID: %d | 错误: %v", productID, err),
			}
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return &domain.BusinessError{
			Code: domain.ErrCodeCacheError,
			Msg:  "执行批量缓存失败: " + err.Error(),
		}
	}
	return nil
}

// GetProductFromCache 从缓存获取商品详情
func (s *CacheService) GetProductFromCache(productID int64) (*domain.Product, error) {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyProductDetail, productID)

	result, err := s.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeCacheError,
			Msg:  "获取缓存商品失败: " + err.Error(),
		}
	}

	var product domain.Product
	err = json.Unmarshal([]byte(result), &product)
	if err != nil {
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeCacheDeserializeError,
			Msg:  "商品反序列化失败: " + err.Error(),
		}
	}
	return &product, nil
}

// InitAllProductStock 预热所有商品库存
func InitAllProductStock() error {
	// 获取所有商品
	products, err := GetProductList()
	if err != nil {
		return err
	}

	// 构建商品库存映射
	stocks := make(map[int64]int32)
	for _, product := range products {
		stocks[product.ProductID] = int32(product.Stock)
	}

	// 使用批量缓存函数
	if err = cacheService.BatchCacheProductStocks(stocks, myredis.DefaultSessionTTL); err != nil {
		log.Printf("[Service] 批量预热商品库存失败: %v", err)
		return err
	}
	log.Printf("商品库存预热完成 | 商品数量: %d", len(products))
	return nil
}

// AddToCartRedis 将商品加入Redis购物车
func (s *CartService) AddToCartRedis(userID, productID int64, quantity int) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyCart, userID)

	// 执行预加载的Lua脚本
	result, err := s.luaScript.Run(ctx, s.client, []string{key}, productID, quantity).Result()
	if err != nil {
		log.Printf("[RedisService] 加入购物车失败 | 用户ID: %d | 商品ID: %d | 错误: %v", userID, productID, err)
		return &domain.BusinessError{
			Code: domain.ErrCodeCacheError,
			Msg:  "加入购物车失败: " + err.Error(),
		}
	}

	// 处理Lua脚本返回结果
	if errMsg, ok := result.(string); ok && errMsg != "" {
		log.Printf("[RedisService] Lua脚本返回错误: %s", errMsg)
		return &domain.BusinessError{
			Code: domain.ErrCodeParamInvalid,
			Msg:  errMsg,
		}
	}

	newQuantity, _ := result.(int64)
	log.Printf("[RedisService] 加入购物车成功 | 用户ID: %d | 商品ID: %d | 新数量: %d", userID, productID, newQuantity)
	return nil
}

// GetCartRedis 获取Redis购物车内容
func (s *CartService) GetCartRedis(userID int64) ([]domain.Cart, error) {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyCart, userID)

	// 获取购物车所有商品
	result, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("[RedisService] 获取购物车失败 | 用户ID: %d | 错误: %v", userID, err)
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeCacheError,
			Msg:  "获取购物车失败: " + err.Error(),
		}
	}

	// 如果购物车为空
	if len(result) == 0 {
		log.Printf("[RedisService] 购物车为空 | 用户ID: %d", userID)
		return []domain.Cart{}, nil
	}

	// 转换Redis Hash到Cart结构
	var carts []domain.Cart
	for productIDStr, quantityStr := range result {
		productID, quantity, err := ParseCartItem(productIDStr, quantityStr)
		if err != nil {
			log.Printf("[RedisService] 解析购物车商品失败 | 用户ID: %d | 商品ID: %s | 数量: %s | 错误: %v", userID, productIDStr, quantityStr, err)
			continue // 跳过无效的商品
		}

		carts = append(carts, domain.Cart{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		})
	}

	log.Printf("[RedisService] 获取购物车成功 | 用户ID：%d | 商品数量：%d", userID, len(carts))
	return carts, nil
}

// ParseCartItem 解析商品ID和数量
func ParseCartItem(productIDStr, quantityStr string) (int64, int, error) {
	var productID int64
	var quantity int

	// 解析商品ID
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		log.Printf("[RedisService] 解析商品ID失败: %s | 错误: %v", productIDStr, err)
	}

	// 解析数量
	quantity, err = strconv.Atoi(quantityStr)
	if err != nil {
		log.Printf("[RedisService] 解析数量失败: %s | 错误: %v", quantityStr, err)
	}
	return productID, quantity, err
}

// ClearCartRedis 清空Redis购物车
func (s *CartService) ClearCartRedis(userID int64) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyCart, userID)

	// 删除整个购物车
	_, err := s.client.Del(ctx, key).Result()
	if err != nil {
		return &domain.BusinessError{
			Code: domain.ErrCodeCacheError,
			Msg:  "清空购物车失败: " + err.Error(),
		}
	}

	return nil
}

// RemoveFromCartRedis 从Redis购物车移除商品
func (s *CartService) RemoveFromCartRedis(userID, productID int64) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyCart, userID)

	// 删除指定商品
	_, err := s.client.HDel(ctx, key, fmt.Sprintf("%d", productID)).Result()
	if err != nil {
		return &domain.BusinessError{
			Code: domain.ErrCodeCacheError,
			Msg:  "移除购物车商品失败: " + err.Error(),
		}
	}

	return nil
}
