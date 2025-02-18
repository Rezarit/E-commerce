package main

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func generateToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix() // 设置有效期为 1 分钟
	// 将密钥转换为字节切片
	key := []byte("114514")
	return token.SignedString(key)
}

func generateRToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // 设置有效期为 1 小时
	// 将密钥转换为字节切片
	key := []byte("114514")
	return token.SignedString(key)
}

func Login(client *gin.Context) {
	//绑定数据
	var user User
	err := client.BindJSON(&user)
	if err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10002,
			"JSON解析失败")
		return
	}

	//验证是否为空
	if user.ID == "" {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"用户 ID 不能为空")
		return
	}

	//密码获取
	if user.Password == "" {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"密码不能为空")
		return
	}

	var password string
	//数据库密码
	cmd := "SELECT password From users WHERE id=?"
	err = db.QueryRow(cmd, user.ID).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(client,
				http.StatusNotFound,
				10003,
				"用户不存在")
		} else {
			sendErrorResponse(client,
				http.StatusInternalServerError,
				10003,
				"数据库查询出错，请稍后重试")
		}
		return
	}

	//核对密码
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"密码错误")
		return
	}

	//生成token和refresh_token
	Token, err := generateToken(user.ID)
	if err != nil {
		sendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("Token密钥转换失败: %v", err))
		return
	}
	RToken, err := generateRToken(user.ID)
	if err != nil {
		sendErrorResponse(client,
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

func IsLogin() gin.HandlerFunc {
	return func(client *gin.Context) {
		//读取头文件
		authHeader := client.Request.Header.Get("Authorization")
		if authHeader == "" {
			sendErrorResponse(
				client,
				http.StatusBadRequest,
				10001,
				"Authorization为空")
			client.Abort()
			return
		}

		//提取Jwt
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			sendErrorResponse(
				client,
				http.StatusBadRequest,
				10003,
				"Jwt分解错误")
			client.Abort()
			return
		}
		tokenstring := parts[1]

		//解析Jwt
		token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte("114514"), nil
		})
		if err != nil || !token.Valid {
			sendErrorResponse(
				client,
				http.StatusBadRequest,
				10002,
				"无效的token")
			client.Abort()
			return
		}

		//验证通过
		client.Set("token", token)
		client.Next()
	}
}

func Refresh(client *gin.Context) {
	authHeader := client.Request.Header.Get("Authorization")
	if authHeader == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"Authorization为空")
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"Jwt分解错误")
		return
	}
	tokenstring := parts[1]

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte("114514"), nil
	})
	if err != nil || !token.Valid {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"无效的token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"无效的token claims")
		return
	}
	if exp, ok := claims["exp"]; ok {
		if float64(time.Now().Unix()) > exp.(float64) {
			sendErrorResponse(
				client,
				http.StatusBadRequest,
				10001,
				"claims已过期")
			client.Abort()
			return
		}
	}

	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	newToken := jwt.New(jwt.SigningMethodHS256)
	secret := []byte("114514")
	signedToken, err := newToken.SignedString(secret)
	if err != nil {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"密钥转换失败")
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "刷新成功",
		"token":  signedToken})
}
