package api

import (
	"database/sql"
	"fmt"
	"github.com/Rezarit/E-commerce/dao"
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Register(client *gin.Context) {
	//绑定数据
	var user domain.User
	if err := client.BindJSON(&user); err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10002,
			"JSON解析失败")
		return
	}

	//检查用户名长度
	domain.ErrorStatus = service.CheckUsernameLen(user.Username)
	if domain.ErrorStatus == 10001 {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"用户名过短")
	}

	//检查密码长度
	domain.ErrorStatus = service.CheckPasswordLen(user.Password)
	if domain.ErrorStatus == 10001 {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"密码过短")
	}

	//检验是否用户名是否存在
	var count string
	is, err := dao.CheckUsernameExists(count)
	if err != nil {
		if err != sql.ErrNoRows {
			service.SendErrorResponse(client,
				http.StatusBadRequest,
				10003,
				fmt.Sprintf("数据未能成功填入数据库1: %v", err))
			return
		}
	}
	if is == true {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"用户名已存在")
		return
	}

	//密码加密
	hashedPassword, err := service.HashedPassword(user.Password)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10003,
			fmt.Sprintf("密码加密失败: %v", err))
		return
	}

	//执行插入指令
	err = dao.InsertPassword(user.Nickname, user.Username, string(hashedPassword))
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("数据未能成功填入数据库2: %v", err))
		return
	}

	//插入成功响应
	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "用户创建成功"})
}

func Password(client *gin.Context) {
	//绑定数据
	var user domain.User
	if err := client.BindJSON(&user); err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"JSON解析失败")
		return
	}

	//检查密码长度
	domain.ErrorStatus = service.CheckPasswordLen(user.Password)
	if domain.ErrorStatus == 10001 {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"密码过短")
		return
	}

	//加密密码
	hashedPassword, err := service.HashedPassword(user.Password)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10003,
			fmt.Sprintf("密码加密失败: %v", err))
		return
	}

	//执行更新指令
	err = dao.UpdatePassword(user.UserID, string(hashedPassword))
	if err != nil {
		service.SendErrorResponse(client,
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
	userIDStr := client.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"无效的 user_id")
		return
	}

	//执行查询指令
	var user domain.User
	user, err = dao.SearchUserMsg(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			service.SendErrorResponse(client,
				http.StatusNotFound,
				10003,
				"未找到该用户")
		} else {
			service.SendErrorResponse(client,
				http.StatusInternalServerError,
				10003,
				fmt.Sprintf("数据库查询错误: %v", err))
		}
		return
	}

	//构造反应体
	type UserResponse struct {
		Status int    `json:"status"`
		Info   string `json:"info"`
		Data   struct {
			User domain.User `json:"user"`
		} `json:"data"`
	}

	response := UserResponse{
		Status: 10000,
		Info:   "success",
	}
	response.Data.User = user

	client.JSON(http.StatusOK, response)
}

func Info(client *gin.Context) {
	var user domain.User
	if err := client.BindJSON(&user); err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"JSON解析失败")
	}

	err := dao.UpdateUserMeg(user.UserID,
		user.Phone,
		user.QQ,
		user.Avatar,
		user.Nickname,
		user.Introduction,
		user.Gender,
		user.Email,
		user.Birthday)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("更新数据失败: %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "信息更新成功"})
}
