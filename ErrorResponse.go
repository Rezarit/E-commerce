package main

import "github.com/gin-gonic/gin"

// 发送错误响应的通用函数
func sendErrorResponse(c *gin.Context, statusCode int, errorStatus int, errorMsg string) {
	c.JSON(statusCode, gin.H{
		"status": errorStatus,
		"info":   errorMsg,
	})
}
