package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetComment(client *gin.Context) {
	productID := client.Param("product_id")

	var comments []Comment

	cmd := "SELECT * FROM comments WHERE product_id=?"
	rows, err := db.Query(cmd, productID)
	if err != nil {
		sendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			"商品列表查询失败")
		return
	}
	defer func() { _ = rows.Close() }() //此处忽略了错误信息

	for rows.Next() {
		var comment Comment

		err = rows.Scan(
			&comment.PostID,
			&comment.PublishTime,
			&comment.Content,
			&comment.UserID,
			&comment.Avatar,
			&comment.Nickname,
			&comment.PraiseCount,
			&comment.ProductID)
		if err != nil {
			sendErrorResponse(client,
				http.StatusInternalServerError,
				10003,
				"获取商品信息失败")
			return
		}

		comment.IsPraised, err = isPraised(comment.PostID, comment.UserID)
		if err != nil {
			sendErrorResponse(
				client,
				http.StatusInternalServerError,
				10003,
				"数据查询错误")
			return
		}

		comments = append(comments, comment)
	}

	client.JSON(http.StatusOK, gin.H{
		"status":   10000,
		"info":     "success",
		"comments": comments})
}

func PutComment(client *gin.Context) {
	productID := client.Param("product_id")

	// 验证商品 ID 是否为空
	if productID == "" {
		sendErrorResponse(client, http.StatusBadRequest, 10002, "商品 ID 不能为空")
		return
	}

	var comment Comment

	if err := client.BindJSON(&comment); err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10002,
			"JSON解析失败")
		return
	}

	if comment.Content == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10002,
			"评论内容不能为空")
		return
	}

	cmd := "INSERT INTO comments (post_id, publish_time, content, user_id, avatar, nickname, praise_count, product_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(cmd,
		comment.PostID,
		comment.PublishTime,
		comment.Content,
		comment.UserID,
		comment.Avatar,
		comment.Nickname,
		comment.PraiseCount,
		productID)
	if err != nil {
		sendErrorResponse(client,
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
	comentID := client.Param("comment_id")
	if comentID == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10002,
			"评论ID不能为空")
		return
	}

	cmd := "DELETE FROM comments WHERE post_id=?"
	_, err := db.Exec(cmd, comentID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(
				client,
				http.StatusInternalServerError,
				10003,
				"未找到该评论")
			return
		} else {
			sendErrorResponse(
				client,
				http.StatusInternalServerError,
				10003,
				"数据库查询失败")
			return
		}
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success"})
}

func UpdateComment(client *gin.Context) {
	commentID := client.Param("comment_id")
	if commentID == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10002,
			"评论ID不能为空")
		return
	}

	var comment Comment
	if err := client.BindJSON(&comment); err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10002,
			"JSON解析失败")
		return
	}

	cmd := "UPDATE comments SET content=?,publish_time=? WHERE post_id=?"
	_, err := db.Exec(cmd, comment.Content, comment.PublishTime, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(
				client,
				http.StatusInternalServerError,
				10003,
				"未找到该评论")
			return
		} else {
			sendErrorResponse(
				client,
				http.StatusInternalServerError,
				10003,
				"数据库查询失败")
			return
		}
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success"})
}

func Praise(client *gin.Context) {
	var comment Comment
	if err := client.BindJSON(&comment); err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10002,
			"JSON解析失败")
		return
	}

	cmd := "INSERT INTO comments (user_id,post_id,like_status) VALUES (?, ?, ?)"
	_, err := db.Exec(cmd, comment.UserID, comment.PostID, comment.IsPraised)
	if err != nil {
		sendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("数据未能成功填入数据库: %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success"})
}

func isPraised(postID, userID string) (int, error) {
	var LikeStatu int
	cmd := "SELECT like_status FROM comment_likes WHERE user_id = ? AND product_id = ?"
	err := db.QueryRow(cmd, postID, userID).Scan(&LikeStatu)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else {
			return -1, err
		}
	}

	if LikeStatu == 1 {
		return 1, nil
	} else if LikeStatu == 2 {
		return 2, nil
	} else if LikeStatu == 0 {
		return 0, nil
	}

	return 0, nil
}
