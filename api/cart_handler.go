package api

import (
	"github.com/Rezarit/go-seckill-system/api/common"
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/Rezarit/go-seckill-system/pkg/response"
	"github.com/Rezarit/go-seckill-system/service"
	"github.com/gin-gonic/gin"
)

// AddToCart 加入购物车
func AddToCart(client *gin.Context) {
	// 从请求参数中获取用户ID和商品ID
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}
	productID := ParseProductID(client)
	if productID == 0 {
		return
	}

	var req domain.AddToCartRequest
	isPass := common.BindRequest(client, &req)
	if !isPass {
		return
	}

	// 执行加入购物车操作
	err := service.AddToCart(userID, productID, req.Quantity)
	if !common.HandleBusinessError(client, err) {
		return
	}

	// 返回成功响应
	response.Success(client, "加入购物车成功", nil)
}

// ShowCart 获取购物车商品列表
func ShowCart(client *gin.Context) {
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	carts, err := service.ShowCart(userID)
	if !common.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "获取购物车成功", carts)
}

// RemoveCart 从购物车移除商品
func RemoveCart(client *gin.Context) {
	// 从请求参数中获取用户ID和商品ID
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}
	productID := ParseProductID(client)
	if productID == 0 {
		return
	}

	// 执行从购物车移除商品操作
	err := service.RemoveFromCart(userID, productID)
	if !common.HandleBusinessError(client, err) {
		return
	}

	// 返回成功响应
	response.Success(client, "从购物车移除商品成功", nil)
}
