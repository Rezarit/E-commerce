package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	myredis "github.com/Rezarit/go-seckill-system/pkg/redis"
)

// DeductStock 扣减商品库存
func (s *StockDeductService) DeductStock(productID int64, quantity int) (int, error) {
	stockKey := myredis.BuildKey(myredis.KeySeckillStock, productID)

	ctx := context.Background()
	result, err := s.luaScript.Run(ctx, s.client, []string{stockKey}, quantity).Int()
	if err != nil {
		return 0, fmt.Errorf("库存减扣操作失败: %v", err)
	}

	switch result {
	case -1:
		return 0, errors.New("商品不存在")
	case -2:
		return 0, errors.New("购买数量必须为正整数")
	case -3:
		return 0, errors.New("库存不足")
	default:
		logger.Sugar.Infof("[StockService] 库存减扣成功 | 商品ID: %d | 数量: %d | 新库存: %d",
			productID, quantity, result)
	}
	return result, nil
}

// AddBackStock 加回库存（补偿：扣减成功后下游失败时调用）
func (s *StockDeductService) AddBackStock(productID int64, quantity int) error {
	stockKey := myredis.BuildKey(myredis.KeySeckillStock, productID)

	ctx := context.Background()
	result, err := s.addBackScript.Run(ctx, s.client, []string{stockKey}, quantity).Int()
	if err != nil {
		return fmt.Errorf("库存加回操作失败: %v", err)
	}

	switch result {
	case -1:
		return errors.New("加回数量必须为正整数")
	case -2:
		// 库存 key 已过期并重新预热，无需补偿
		logger.Sugar.Infof("[StockService] 库存 key 不存在，跳过补偿 | 商品ID: %d", productID)
	default:
		logger.Sugar.Infof("[StockService] 库存加回成功 | 商品ID: %d | 数量: %d | 新库存: %d",
			productID, quantity, result)
	}
	return nil
}
