package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MakeOrder(client *gin.Context) {
	var req Order
	// 绑定请求体中的 JSON 数据
	if err := client.BindJSON(&req); err != nil {
		sendErrorResponse(client, http.StatusBadRequest, 10001, "请求参数格式错误")
		return
	}

	// 输入验证
	if req.Address == "" || req.Total <= 0 || req.UserID == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"地址、总价或用户 ID 不能为空")
		return
	}

	//执行插入语句
	cmd := "INSERT INTO orders (address, total, user_id) VALUES (?, ?, ?)"
	result, err := db.Exec(cmd, req.Address, req.Total, req.UserID)
	if err != nil {
		sendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("下单失败: %v", err))
		return
	}

	// 获取插入的订单 ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		sendErrorResponse(
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
