package service

import (
	"context"
	"errors"
	"fmt"
	myredis "github.com/Rezarit/go-seckill-system/pkg/redis"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

type StockDeductService struct {
	client    *redis.Client
	luaScript *redis.Script
}

// NewStockDeductService 创建一个新的库存扣减服务实例
func NewStockDeductService() (*StockDeductService, error) {
	scriptContent, err := loadLuaScript("scripts/lua/deduct_stock.lua")
	if err != nil {
		return nil, err
	}

	return &StockDeductService{
		client:    myredis.GetClient(),
		luaScript: redis.NewScript(scriptContent),
	}, nil
}

// loadLuaScript 从文件加载Lua脚本
func loadLuaScript(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("加载Lua脚本失败: %v", err)
		return "", err
	}
	return string(content), nil
}

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
		log.Printf("[StockService] 库存减扣成功 | 商品ID: %d | 数量: %d | 新库存: %d",
			productID, quantity, result)
	}
	return result, nil
}