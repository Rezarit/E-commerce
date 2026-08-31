package service

import (
	dao2 "github.com/Rezarit/go-seckill-system/internal/dao"
	domain2 "github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"gorm.io/gorm"
)

// AddToCart 加入购物车
func AddToCart(userID, productID int64, quantity int) error {
	logger.Sugar.Infof("[Service] 加入购物车 | 用户ID：%d | 商品ID：%d | 数量：%d", userID, productID, quantity)

	if quantity <= 0 {
		return &domain2.BusinessError{
			Code: domain2.ErrCodeParamInvalid,
			Msg:  "数量必须大于0",
		}
	}

	// 双写：先写 MySQL（持久源），再写 Redis（缓存，供高频读写）
	// 顺序有讲究：若 Redis 失败，读兜底能从 MySQL 回填自愈；反之若 MySQL 失败则数据会丢
	item := domain2.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
	if err := dao2.AddToCart(item); err != nil {
		logger.Sugar.Errorf("[Service] 加入购物车(MySQL)失败 | 用户ID：%d | 商品ID：%d | 错误：%v", userID, productID, err)
		return &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "加入购物车失败"}
	}

	// 再写 Redis（缓存，供高频读写）
	err := cartService.AddToCartRedis(userID, productID, quantity)
	if err != nil {
		logger.Sugar.Errorf("[Service] 加入购物车(Redis)失败 | 用户ID：%d | 商品ID：%d | 错误：%v", userID, productID, err)
		return err
	}

	logger.Sugar.Infof("[Service] 加入购物车成功 | 用户ID：%d | 商品ID：%d", userID, productID)
	return nil
}

// RemoveFromCart 从购物车移除商品
func RemoveFromCart(userID, productID int64) error {
	logger.Sugar.Infof("[Service] 从购物车移除商品 | 用户ID：%d | 商品ID：%d", userID, productID)

	// 先删 MySQL（持久化）
	if err := dao2.RemoveFromCart(userID, productID); err != nil {
		logger.Sugar.Errorf("[Service] 从购物车移除商品(MySQL)失败 | 用户ID：%d | 商品ID：%d | 错误：%v", userID, productID, err)
		return &domain2.BusinessError{
			Code: domain2.ErrCodeDBError,
			Msg:  "从购物车移除商品失败",
		}
	}

	// 再删 Redis（缓存）
	if err := cartService.RemoveFromCartRedis(userID, productID); err != nil {
		logger.Sugar.Errorf("[Service] 从购物车移除商品(Redis)失败 | 用户ID：%d | 商品ID：%d | 错误：%v", userID, productID, err)
		return &domain2.BusinessError{
			Code: domain2.ErrCodeDBError,
			Msg:  "从购物车移除商品失败",
		}
	}

	logger.Sugar.Infof("[Service] 从购物车移除商品成功 | 用户ID：%d | 商品ID：%d", userID, productID)
	return nil
}

// processCartItem 处理单个订单商品
func processCartItem(tx *gorm.DB, orderID int64, item domain2.OrderItemMsg) error {
	// 获取商品信息
	product, err := getProductInfo(item.ProductID)
	if err != nil {
		return err
	}

	// 创建订单商品
	if err = createOrderItem(tx, orderID, item, product); err != nil {
		return err
	}

	// 扣减库存（只扣 MySQL，Redis 已在入口扣减）
	if err = updateProductStock(tx, product, item.Quantity); err != nil {
		return err
	}

	return nil
}

// ClearCartInRedis 清空购物车
func ClearCartInRedis(userID int64) error {
	err := cartService.ClearCartRedis(userID)
	if err != nil {
		logger.Sugar.Errorf("[Service] 清空购物车失败 | 用户ID：%d | 错误：%v", userID, err)
		return err
	}

	logger.Sugar.Infof("[Service] 清空购物车成功 | 用户ID：%d", userID)
	return nil
}

// CheckCart 检查购物车是否为空
func CheckCart(items []domain2.CartItem) error {
	if len(items) == 0 {
		return &domain2.BusinessError{
			Code: domain2.ErrCodeCartEmpty,
			Msg:  "购物车为空",
		}
	}
	return nil
}
