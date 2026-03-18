package service

import (
	"github.com/Rezarit/go-seckill-system/dao"
	"github.com/Rezarit/go-seckill-system/domain"
	"log"
)

// AddToCart 加入购物车
func AddToCart(userID, productID int64, quantity int) error {
	log.Printf("[Service] 加入购物车 | 用户ID：%d | 商品ID：%d | 数量：%d", userID, productID, quantity)

	if quantity <= 0 {
		return &domain.BusinessError{
			Code: domain.ErrCodeParamInvalid,
			Msg:  "数量必须大于0",
		}
	}

	// 使用Redis服务加入购物车
	err := cartService.AddToCartRedis(userID, productID, quantity)
	if err != nil {
		log.Printf("[Service] 加入购物车失败 | 用户ID：%d | 商品ID：%d | 错误：%v", userID, productID, err)
		return err
	}

	log.Printf("[Service] 加入购物车成功 | 用户ID：%d | 商品ID：%d", userID, productID)
	return nil
}

// ShowCart 获取购物车商品列表
func ShowCart(userID int64) ([]domain.Cart, error) {
	log.Printf("[Service] 获取购物车商品列表 | 用户ID：%d", userID)

	carts, err := cartService.GetCartRedis(userID)
	if err != nil {
		log.Printf("[Service] 获取购物车商品列表失败 | 用户ID：%d | 错误：%v", userID, err)
		return nil, err
	}

	log.Printf("[Service] 获取购物车商品列表成功 | 用户ID：%d | 商品数量：%d", userID, len(carts))
	return carts, nil
}

// RemoveFromCart 从购物车移除商品
func RemoveFromCart(userID, productID int64) error {
	log.Printf("[Service] 从购物车移除商品 | 用户ID：%d | 商品ID：%d", userID, productID)

	err := dao.RemoveFromCart(userID, productID)
	if err != nil {
		log.Printf("[Service] 从购物车移除商品失败 | 用户ID：%d | 商品ID：%d | 错误：%v", userID, productID, err)
		return &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "从购物车移除商品失败",
		}
	}

	log.Printf("[Service] 从购物车移除商品成功 | 用户ID：%d | 商品ID：%d", userID, productID)
	return nil
}
