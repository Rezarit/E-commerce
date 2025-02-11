package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func Login(client *gin.Context) {
	//绑定数据
	var user User
	err := client.BindJSON(&user)
	if err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"status": 10000,
			"info":   "JSON解析失败"})
		return
	}

	password := client.Query("password")
	if password != user.Password {
		client.JSON(http.StatusBadRequest, gin.H{
			"status": 10000,
			"info":   "密码错误"})
		return
	}

	//生成token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["Id"] = user.ID
	claims["Username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	secret := []byte("114514")
	signedToken, err := token.SignedString(secret)
	if err != nil {
		client.JSON(10001, gin.H{"error": "密钥转换失败"})
		return
	}

	client.JSON(http.StatusOK, gin.H{"token": signedToken})
}

func IsLogin() gin.HandlerFunc {
	return func(client *gin.Context) {
		//读取头文件
		authHeader := client.Request.Header.Get("Authorization")
		if authHeader == "" {
			client.JSON(10001, gin.H{"error": "Authorization为空"})
			client.Abort()
			return
		}

		//提取Jwt
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			client.JSON(10001, gin.H{"error": "Jwt分解错误"})
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
			client.JSON(10001, gin.H{"error": "无效的token"})
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
		client.JSON(10001, gin.H{"error": "Authorization为空"})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		client.JSON(10001, gin.H{"erroe": "Jwt分解错误"})
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
		client.JSON(10001, gin.H{"error": "无效的token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		client.JSON(10001, gin.H{"error": "无效的token claims"})
		return
	}
	if exp, ok := claims["exp"]; ok {
		if float64(time.Now().Unix()) > exp.(float64) {
			client.JSON(10001, gin.H{"error": "claims已过期"})
			client.Abort()
			return
		}
	}

	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	newToken := jwt.New(jwt.SigningMethodHS256)
	secret := []byte("114514")
	signedToken, err := newToken.SignedString(secret)
	if err != nil {
		client.JSON(10001, gin.H{"error": "密钥转换失败"})
		return
	}

	client.JSON(10000, gin.H{"token": signedToken})
}
