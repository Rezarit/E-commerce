package main

import (
	"database/sql"
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
	return token.SignedString("114514")
}

func generateRToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // 设置有效期为 1 小时
	return token.SignedString("114514")
}

func Login(client *gin.Context) {
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

	//密码获取
	password := client.Query("password")
	if password == "" {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"密码不能为空")
		return
	}

	//数据库密码
	cmd := "SELECT password From users WHERE id=?"
	err = db.QueryRow(cmd, userID).Scan(&user.Password)
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
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"密码错误")
		return
	}

	//生成token和refresh_token
	Token, err := generateToken(userID)
	if err != nil {
		sendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			"Token密钥转换失败")
		return
	}
	RToken, err := generateRToken(userID)
	if err != nil {
		sendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			"RToken密钥转换失败")
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
