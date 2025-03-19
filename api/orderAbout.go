package api

import (
	"fmt"
	"github.com/Rezarit/E-commerce/dao"
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MakeOrder(client *gin.Context) {
	var req domain.Order
	// 绑定请求体中的 JSON 数据
	if err := client.BindJSON(&req); err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"请求参数格式错误")
		return
	}

	// 输入验证
	if req.Address == "" || req.Total <= 0 || req.UserID == "" {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"地址、总价或用户 ID 不能为空")
		return
	}

	result, err := dao.Order(req)

	// 获取插入的订单 ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		service.SendErrorResponse(
			client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("获取订单 ID 失败: %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status":   10000,
		"info":     "下单成功",
		"order_id": fmt.Sprintf("%d", lastInsertID)})
}
