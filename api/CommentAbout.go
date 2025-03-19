package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Rezarit/E-commerce/dao"
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetComment(client *gin.Context) {
	productIDstr := client.Param("product_id")
	productID, err := strconv.Atoi(productIDstr)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusInternalServerError,
			10001,
			fmt.Sprintf("获取产品信息失败: %v", err))
	}

	rows, err := dao.GetCommentbyProductID(productID)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("数据库查询失败: %v", err))
	}

	comments, err := service.PutCommentMsg(rows)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("获取评论信息失败: %v", err))
	}

	client.JSON(http.StatusOK, gin.H{
		"status":   10000,
		"info":     "success",
		"comments": comments})
}

func PutComment(client *gin.Context) {
	productIDstr := client.Param("product_id")
	productID, err := strconv.Atoi(productIDstr)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"无效产品 ID ")
	}

	// 验证商品 ID 是否为空
	if productID == 0 {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10002,
			"商品 ID 不能为空")
		return
	}

	var comment domain.Comment
	if err := client.BindJSON(&comment); err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10002,
			"JSON解析失败")
		return
	}

	if comment.Content == "" {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10002,
			"评论内容不能为空")
		return
	}

	err = dao.InsertComment(comment, productID)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("数据未能成功填入数据库: %v", err))
		return
	}

	//插入成功响应
	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "评论发表成功"})
}

func DeleteComment(client *gin.Context) {
	comentIDstr := client.Param("comment_id")
	comentID, err := strconv.Atoi(comentIDstr)
	if err != nil {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"无效评论 ID")
		return
	}

	if comentID == 0 {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"评论ID不能为空")
		return
	}

	err = dao.DeleteComment(comentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.SendErrorResponse(
				client,
				http.StatusBadRequest,
				10003,
				"未找到该评论")
			return
		} else {
			service.SendErrorResponse(
				client,
				http.StatusInternalServerError,
				10003,
				fmt.Sprintf("数据库查询失败: %v", err))
			return
		}
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success"})
}

func UpdateComment(client *gin.Context) {
	commentIDstr := client.Param("comment_id")
	commentID, err := strconv.Atoi(commentIDstr)
	if commentID == 0 {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10002,
			"评论ID不能为空")
		return
	}

	var comment domain.Comment
	if err := client.BindJSON(&comment); err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10002,
			"JSON解析失败")
		return
	}

	err = dao.UpdataComment(comment, commentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.SendErrorResponse(
				client,
				http.StatusBadRequest,
				10003,
				"未找到该评论")
			return
		} else {
			service.SendErrorResponse(
				client,
				http.StatusInternalServerError,
				10003,
				fmt.Sprintf("数据库查询失败: %v", err))
			return
		}
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success"})
}

func Praise(client *gin.Context) {
	var comment domain.Comment
	if err := client.BindJSON(&comment); err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10002,
			"JSON解析失败")
		return
	}

	err := dao.PraiseComment(comment)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("数据未能成功填入数据库: %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success"})
}
