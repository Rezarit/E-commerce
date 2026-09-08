package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	domain2 "github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	myredis "github.com/Rezarit/go-seckill-system/pkg/redis"
	"github.com/go-redis/redis/v8"
)

// SetOrderResult 写入订单处理结果（键 = msg_id）。
// 调用时机：MySQL 建单成功之后——result 是「建单成功的盖章回执」，查它不会误报成功。
func (s *OrderService) SetOrderResult(msgID string, result domain2.OrderResult, expiration time.Duration) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyOrderResult, msgID)

	data, err := json.Marshal(result)
	if err != nil {
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheSerializeError,
			Msg:  "序列化订单结果失败: " + err.Error(),
		}
	}

	err = s.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
			Msg:  "写入订单结果失败: " + err.Error(),
		}
	}
	return nil
}

// GetOrderResult 读取订单处理结果。返回 ok=false 表示键不存在（查无 = processing 态）。
func (s *OrderService) GetOrderResult(msgID string) (domain2.OrderResult, bool, error) {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyOrderResult, msgID)

	val, err := s.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return domain2.OrderResult{}, false, nil
	}
	if err != nil {
		return domain2.OrderResult{}, false, &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
			Msg:  "获取订单结果失败: " + err.Error(),
		}
	}

	var result domain2.OrderResult
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return domain2.OrderResult{}, false, &domain2.BusinessError{
			Code: domain2.ErrCodeCacheDeserializeError,
			Msg:  "反序列化订单结果失败: " + err.Error(),
		}
	}
	return result, true, nil
}

// ClaimDeadProcessed 抢死信处理幂等标记（键 = msg_id）。
// 标记语义 =「该单库存释放已完成」：返回 true = 我是第一个、可执行精确释放；
// 返回 false = 释放已完成过，跳过（防 MQ 重投/重复消费 → 重复释放 → 反向超卖）。
func (s *OrderService) ClaimDeadProcessed(msgID string) (bool, error) {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyDeadProcessed, msgID)
	return s.client.SetNX(ctx, key, 1, myredis.DefaultSessionTTL).Result()
}

// ReleaseDeadProcessed 撤销死信处理标记：释放动作失败时调用，
// 给可能到来的 MQ 重投留补救机会（撤销 = 「还没干完，别拦着重投重试」）。
func (s *OrderService) ReleaseDeadProcessed(msgID string) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyDeadProcessed, msgID)
	return s.client.Del(ctx, key).Err()
}

// CacheNullProduct 缓存空商品到Redis
func (s *CacheService) CacheNullProduct(productID int64, expiration time.Duration) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyProductNull, productID)

	err := s.client.Set(ctx, key, productID, expiration).Err()
	if err != nil {
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
			Msg:  "缓存空商品失败: " + err.Error(),
		}
	}
	return nil
}

// CacheProduct 缓存商品详情
func (s *CacheService) CacheProduct(product *domain2.Product, expiration time.Duration) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyProductDetail, product.ProductID)

	data, err := json.Marshal(product)
	if err != nil {
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheSerializeError,
			Msg:  "商品序列化失败: " + err.Error(),
		}
	}

	err = s.client.Set(ctx, key, string(data), expiration).Err()
	if err != nil {
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
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
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
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
			return &domain2.BusinessError{
				Code: domain2.ErrCodeCacheError,
				Msg:  fmt.Sprintf("批量缓存商品库存失败 | 商品ID: %d | 错误: %v", productID, err),
			}
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
			Msg:  "执行批量缓存失败: " + err.Error(),
		}
	}
	return nil
}

// GetProductFromCache 从缓存获取商品详情
// 返回 found：true=命中返回商品，false=缓存未命中（key 不存在，需查 MySQL）
func (s *CacheService) GetProductFromCache(productID int64) (*domain2.Product, bool, error) {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyProductDetail, productID)

	result, err := s.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil // 未命中，不是错误
	}
	if err != nil {
		return nil, false, &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
			Msg:  "获取缓存商品失败: " + err.Error(),
		}
	}

	var product domain2.Product
	err = json.Unmarshal([]byte(result), &product)
	if err != nil {
		return nil, false, &domain2.BusinessError{
			Code: domain2.ErrCodeCacheDeserializeError,
			Msg:  "商品反序列化失败: " + err.Error(),
		}
	}
	return &product, true, nil
}

// TryLockProduct 尝试获取商品缓存重建锁（SETNX 原子抢锁，带 TTL 防死锁）
// 返回 true = 抢到锁（调用方负责用完释放）
func (s *CacheService) TryLockProduct(productID int64) (bool, error) {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyProductLock, productID)

	locked, err := s.client.SetNX(ctx, key, 1, myredis.DefaultSeckillLockTTL).Result()
	if err != nil {
		return false, err
	}
	return locked, nil
}

// UnlockProduct 释放商品缓存重建锁
func (s *CacheService) UnlockProduct(productID int64) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyProductLock, productID)

	return s.client.Del(ctx, key).Err()
}

// GetNullProduct 检查是否命中商品空值缓存（防穿透）
func (s *CacheService) GetNullProduct(productID int64) (bool, error) {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyProductNull, productID)

	_, err := s.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
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
		logger.Sugar.Errorf("[Service] 批量预热商品库存失败: %v", err)
		return err
	}
	logger.Sugar.Infof("商品库存预热完成 | 商品数量: %d", len(products))
	return nil
}

// AddToCartRedis 将商品加入Redis购物车
func (s *CartService) AddToCartRedis(userID, productID int64, quantity int) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyCart, userID)

	// 执行预加载的Lua脚本
	result, err := s.luaScript.Run(ctx, s.client, []string{key}, productID, quantity).Result()
	if err != nil {
		logger.Sugar.Errorf("[RedisService] 加入购物车失败 | 用户ID: %d | 商品ID: %d | 错误: %v", userID, productID, err)
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
			Msg:  "加入购物车失败: " + err.Error(),
		}
	}

	// 处理Lua脚本返回结果
	if errMsg, ok := result.(string); ok && errMsg != "" {
		logger.Sugar.Errorf("[RedisService] Lua脚本返回错误: %s", errMsg)
		return &domain2.BusinessError{
			Code: domain2.ErrCodeParamInvalid,
			Msg:  errMsg,
		}
	}

	newQuantity, _ := result.(int64)
	logger.Sugar.Infof("[RedisService] 加入购物车成功 | 用户ID: %d | 商品ID: %d | 新数量: %d", userID, productID, newQuantity)
	return nil
}

// GetCartRedis 获取Redis购物车内容
func (s *CartService) GetCartRedis(userID int64) ([]domain2.CartItem, error) {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyCart, userID)

	// 获取购物车所有商品
	result, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		logger.Sugar.Errorf("[RedisService] 获取购物车失败 | 用户ID: %d | 错误: %v", userID, err)
		return nil, &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
			Msg:  "获取购物车失败: " + err.Error(),
		}
	}

	// 如果购物车为空
	if len(result) == 0 {
		logger.Sugar.Infof("[RedisService] 购物车为空 | 用户ID: %d", userID)
		return []domain2.CartItem{}, nil
	}

	// 转换Redis Hash到Cart结构
	var items []domain2.CartItem
	for productIDStr, quantityStr := range result {
		productID, quantity, err := ParseCartItem(productIDStr, quantityStr)
		if err != nil {
			logger.Sugar.Errorf("[RedisService] 解析购物车商品失败 | 用户ID: %d | 商品ID: %s | 数量: %s | 错误: %v", userID, productIDStr, quantityStr, err)
			continue // 跳过无效的商品
		}

		items = append(items, domain2.CartItem{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		})
	}

	logger.Sugar.Infof("[RedisService] 获取购物车成功 | 用户ID：%d | 商品数量：%d", userID, len(items))
	return items, nil
}

// CacheCartItems 将购物车商品列表回填到 Redis Hash
func (s *CartService) CacheCartItems(userID int64, items []domain2.CartItem, expiration time.Duration) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyCart, userID)

	pipe := s.client.Pipeline()
	for _, item := range items {
		pipe.HSet(ctx, key, strconv.FormatInt(item.ProductID, 10), item.Quantity)
	}
	pipe.Expire(ctx, key, expiration)
	if _, err := pipe.Exec(ctx); err != nil {
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
			Msg:  "回填购物车缓存失败",
		}
	}
	return nil
}

// CacheNullCart 缓存空购物车标记（防穿透）
func (s *CartService) CacheNullCart(userID int64, expiration time.Duration) error {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyCartNull, userID)
	if err := s.client.Set(ctx, key, 1, expiration).Err(); err != nil {
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
			Msg:  "缓存空购物车失败",
		}
	}
	return nil
}

// IsNullCart 检查是否命中空购物车标记
func (s *CartService) IsNullCart(userID int64) (bool, error) {
	ctx := context.Background()
	key := myredis.BuildKey(myredis.KeyCartNull, userID)
	_, err := s.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ParseCartItem 解析商品ID和数量
func ParseCartItem(productIDStr, quantityStr string) (int64, int, error) {
	var productID int64
	var quantity int

	// 解析商品ID
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		logger.Sugar.Errorf("[RedisService] 解析商品ID失败: %s | 错误: %v", productIDStr, err)
	}

	// 解析数量
	quantity, err = strconv.Atoi(quantityStr)
	if err != nil {
		logger.Sugar.Errorf("[RedisService] 解析数量失败: %s | 错误: %v", quantityStr, err)
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
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
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
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCacheError,
			Msg:  "移除购物车商品失败: " + err.Error(),
		}
	}

	return nil
}
