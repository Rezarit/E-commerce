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

func ShowList(client *gin.Context) {
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

	//创建产品列表
	var products []domain.Product

	products, err := service.ListSearch(user)
	if err != nil {
		service.SendErrorResponse(
			client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("获取列表失败: %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status":   10000,
		"info":     "success",
		"products": products,
	})
}

func SearchProduct(client *gin.Context) {
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

	//验证是否为空
	if productIDStr == "" {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"商品ID不能为空")
	}

	var product domain.Product

	product, err = dao.SearchProductbyID(productID)

	product.IsaddedCart, err = dao.IsAdded(user.UserID, productID)
	if err != nil {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10003,
			fmt.Sprintf("数据库查询错误: %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status":   10000,
		"info":     "success",
		"products": product,
	})
}

func ProductDetail(client *gin.Context) {
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

	var product domain.Product

	//实现查询功能
	product, err = dao.SearchProductbyID(productID)
	if err != nil {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10003,
			fmt.Sprintf("数据库查询错误: %v", err))
		return
	}

	product.IsaddedCart, err = dao.IsAdded(user.UserID, productID)
	if err != nil {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10003,
			fmt.Sprintf("数据库查询错误: %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status":   10000,
		"info":     "success",
		"products": product})
}

func GetType(client *gin.Context) {
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

	//获取type
	productType := client.Param("type")

	//检验type
	if productType == "" {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"类型不能为空")
		return
	}

	rows, err := dao.SearchProductbyType(productType)
	if err != nil {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			fmt.Sprintf("数据库查询错误: %v", err))
		return
	}

	products, err := service.TypeSearch(user, rows)
	if err != nil {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			fmt.Sprintf("查询失败: %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success",
		"data":   products})
}
