package api

import (
	"database/sql"
	"fmt"
	"github.com/Rezarit/E-commerce/dao"
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(client *gin.Context) {
	//绑定数据
	var user domain.User
	err := client.BindJSON(&user)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10002,
			"JSON解析失败")
		return
	}

	//验证是否为空
	is := service.CheckIDIsNil(user.UserID)
	if is == true {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"用户 ID 不能为空")
		return
	}

	//密码获取
	is = service.CheckPasswordIsNil(user.Password)
	if is == true {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"密码不能为空")
		return
	}

	//数据库密码
	password, err := dao.SearchPassword(user.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			service.SendErrorResponse(client,
				http.StatusNotFound,
				10003,
				"用户不存在")
		} else {
			service.SendErrorResponse(client,
				http.StatusInternalServerError,
				10003,
				"数据库查询出错，请稍后重试")
		}
		return
	}

	//核对密码
	is = service.ComparePassword(password, user.Password)
	if is == false {
		service.SendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"密码错误")
		return
	}

	//生成token和refresh_token
	Token, err := service.GenerateToken(user.UserID)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("Token密钥转换失败: %v", err))
		return
	}
	RToken, err := service.GenerateRToken(user.UserID)
	if err != nil {
		service.SendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("RToken密钥转换失败: %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status":        10000,
		"info":          "登录成功",
		"refresh_token": RToken,
		"token":         Token})
}

func Refresh(client *gin.Context) {
	authHeader := client.Request.Header.Get("Authorization")
	if authHeader == "" {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"Authorization为空")
		return
	}

	tokenstring := service.CutauthHeader(authHeader)
	if tokenstring == "" {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"Jwt分解错误")
		return
	}

	token, err := service.ParseToken(authHeader)
	if err != nil {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"无效的token")
		return
	}

	domain.ErrorMsg = service.ClaimToken(token)
	if domain.ErrorMsg != "" {
		service.SendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"无效的token claims")
		return
	}

	var signedToken string
	domain.ErrorMsg, signedToken = service.NewToken()

	client.JSON(http.StatusOK, gin.H{
		"status":        10000,
		"info":          "刷新成功",
		"refresh_token": tokenstring,
		"token":         signedToken})
}
