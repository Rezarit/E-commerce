package service

import "log"

// 全局服务实例
var (
	redisService       *RedisService
	stockDeductService *StockDeductService
)

// InitServiceInstances 初始化所有服务实例
func InitServiceInstances() error {
	log.Println("[Service] 开始初始化服务实例...")

	// 初始化添加到购物车服务
	var err error
	redisService, err = NewAddToCartService()
	if err != nil {
		log.Printf("[Service] 购物车服务初始化失败: %v", err)
		return err
	}
	log.Println("[Service] 购物车服务初始化完成")

	// 初始化库存减扣服务
	stockDeductService, err = NewStockDeductService()
	if err != nil {
		log.Printf("[Service] 库存减扣服务初始化失败: %v", err)
		return err
	}
	log.Println("[Service] 库存减扣服务初始化完成")

	log.Println("[Service] 所有服务实例初始化完成")
	return nil
}
