package api

import (
	"fmt"
	"github.com/Rezarit/E-commerce/dao"
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddCart(client *gin.Context) {
	//绑定数据
	var user domain.User
	if err := client.BindJSON(&user); err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"JSON解析失败")
		return
	}

	//检查用户ID
	if user.UserID == 0 {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"用户 ID 不能为空")
		return
	}

	//获取商品ID
	productIDStr := client.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"无效的 product_id")
		return
	}

	//检查ID
	if productIDStr == "" {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"商品 ID 不能为空")
		return
	}

	err = dao.AddintoCart(user.UserID, productID)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			fmt.Sprintf("加入购物车失败: %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success"})
}

func ShowCart(client *gin.Context) {
	//绑定数据
	var user domain.User
	if err := client.BindJSON(&user); err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"JSON解析失败")
		return
	}

	//检查用户ID
	if user.UserID == 0 {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"用户 ID 不能为空")
		return
	}

	rows, err := dao.CartSearchByID(user.UserID)
	if err != nil {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10003,
			fmt.Sprintf("数据库查询错误: %v", err))
		return
	}

	productIDs, err := service.CartSearch(rows)
	if err != nil {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10003,
			fmt.Sprintf("查询错误: %v", err))
		return
	}

	var args []any
	for _, id := range productIDs {
		args = append(args, id)
	}

	if len(productIDs) == 0 {
		client.JSON(http.StatusOK, gin.H{
			"status": 10000,
			"info":   "success",
			"data":   "空的购物车"})
		return
	}

	products, err := service.SearchCart(productIDs, args)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10003,
			fmt.Sprintf(": %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success",
		"data":   products})
}
