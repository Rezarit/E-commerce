package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(client *gin.Context) {
	//绑定数据
	var user User
	if err := client.BindJSON(&user); err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"status": 10001,
			"info":   "JSON解析失败"})
		return
	}

	//检查用户名长度
	if 1 > len(user.Username) {
		client.JSON(http.StatusBadRequest, gin.H{
			"status": 10001,
			"info":   "用户名过短"})
		return
	}

	//检验密码长度
	if 8 > len(user.Password) {
		client.JSON(http.StatusBadRequest, gin.H{
			"status": 10001,
			"info":   "密码过短"})
		return
	}

	//检验是否用户名是否存在
	var count int
	cmd := "SELECT COUNT(*) FROM users WHERE username =?"
	err := db.QueryRow(cmd, user.Username).Scan(&count)
	if err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"status": 10001,
			"info":   fmt.Sprintf("数据未能成功填入数据库: %v", err)})
		return
	}
	if count > 0 {
		client.JSON(http.StatusBadRequest, gin.H{
			"status": 10001,
			"info":   "用户名已存在"})
		return
	}

	//执行插入指令
	cmd = "INSERT INTO users(nickname,username,password) VALUES (?,?,?);"
	_, err = db.Exec(cmd, user.Nickname, user.Username, user.Password)
	if err != nil {
		client.JSON(http.StatusInternalServerError, gin.H{
			"status": 10001,
			"info":   fmt.Sprintf("数据未能成功填入数据库: %v", err)})
		return
	}

	//插入成功响应
	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "用户创建成功"})
}

func Password(client *gin.Context) {
	//绑定数据
	var user User
	if err := client.BindJSON(&user); err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"status": 10001,
			"info":   "JSON解析失败"})
		return
	}

	//核对密码
	if client.Query("password") != user.Password {
		client.JSON(http.StatusBadRequest, gin.H{
			"status": 10001,
			"info":   "原密码错误"})
		return
	}

	//执行更新指令
	cmd := "UPDATE users SET password=? WHERE Id=?;"
	_, err := db.Exec(cmd, user.Password, user.ID)
	if err != nil {
		client.JSON(http.StatusInternalServerError, gin.H{
			"status": 10001,
			"info":   "密码更新失败"})
	}

	//成功响应
	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "密码更新成功"})
}

func GetInfoById(client *gin.Context) {
	//获取用户ID
	userID := client.Param("user_id")

	//验证是否为空
	if userID == "" {
		client.JSON(http.StatusBadRequest, gin.H{
			"status": 10001,
			"info":   "用户 ID 不能为空"})
	}

	//执行查询指令
	var user User
	cmd := "SELECT id, avatar, nickname, introduction, phone, qq, gender, email, birthday FROM users WHERE id=?;"
	err := db.QueryRow(cmd, userID).Scan(&user.ID, &user.Avatar, &user.Nickname, &user.Introduction, &user.Phone, &user.QQ, &user.Gender, &user.Email, &user.Birthday)
	if err != nil {
		if err == sql.ErrNoRows {
			client.JSON(http.StatusNotFound, gin.H{
				"status": 10001,
				"info":   "未找到该用户"})
		} else {
			client.JSON(http.StatusInternalServerError, gin.H{
				"status": 10001,
				"info":   fmt.Sprintf("数据库查询错误: %v", err)})
		}
		return
	}

	//构造反应体
	response := Response{
		Status: 10000,
		Info:   "success",
	}
	response.Data.User = user

	client.JSON(http.StatusOK, response)
}
