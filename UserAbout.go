package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Register(client *gin.Context) {
	//绑定数据
	var user User
	if err := client.BindJSON(&user); err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10002,
			"JSON解析失败")
		return
	}

	//检查用户名长度
	if 1 > len(user.Username) {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"用户名过短")
		return
	}

	//检验密码长度
	if 8 > len(user.Password) {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"密码过短")
		return
	}

	//检验是否用户名是否存在
	var count int
	cmd := "SELECT COUNT(*) FROM users WHERE username =?"
	err := db.QueryRow(cmd, user.Username).Scan(&count)
	if err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10003,
			fmt.Sprintf("数据未能成功填入数据库: %v", err))
		return
	}
	if count > 0 {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"用户名已存在")
		return
	}

	//密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10003,
			fmt.Sprintf("密码加密失败: %v", err))
		return
	}

	//执行插入指令
	cmd = "INSERT INTO users(nickname,username,password) VALUES (?,?,?);"
	_, err = db.Exec(cmd, user.Nickname, user.Username, string(hashedPassword))
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
		"info":   "用户创建成功"})
}

func Password(client *gin.Context) {
	//绑定数据
	var user User
	if err := client.BindJSON(&user); err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"JSON解析失败")
		return
	}

	//检验长度
	if 8 > len(user.Password) {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"密码过短")
	}

	//加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10003,
			fmt.Sprintf("密码加密失败: %v", err))
		return
	}

	//执行更新指令
	cmd := "UPDATE users SET password=? WHERE Id=?;"
	_, err = db.Exec(cmd, hashedPassword, user.ID)
	if err != nil {
		sendErrorResponse(client,
			http.StatusInternalServerError,
			10001,
			fmt.Sprintf("密码更新失败: %v", err))
		return
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
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"用户 ID 不能为空")
		return
	}

	//执行查询指令
	var user User
	cmd := "SELECT id, avatar, nickname, introduction, phone, qq, gender, email, birthday FROM users WHERE id=?;"
	err := db.QueryRow(cmd, userID).Scan(
		&user.ID,
		&user.Avatar,
		&user.Nickname,
		&user.Introduction,
		&user.Phone,
		&user.QQ,
		&user.Gender,
		&user.Email,
		&user.Birthday)
	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(client,
				http.StatusNotFound,
				10001,
				"未找到该用户")
		} else {
			sendErrorResponse(client,
				http.StatusInternalServerError,
				10001,
				fmt.Sprintf("数据库查询错误: %v", err))
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
