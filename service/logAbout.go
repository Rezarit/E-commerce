package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func CheckIDIsNil(ID int) bool {
	if 0 == ID {
		return true
	}
	return false
}

func CheckPasswordIsNil(password string) bool {
	if password == "" {
		return true
	}
	return false
}

func ComparePassword(input, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(input))
	if err != nil {
		return false
	}
	return true
}

func GenerateToken(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix() // 设置有效期为 1 分钟
	// 将密钥转换为字节切片
	key := []byte("114514")
	return token.SignedString(key)
}

func GenerateRToken(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // 设置有效期为 1 小时
	// 将密钥转换为字节切片
	key := []byte("114514")
	return token.SignedString(key)
}

func CutauthHeader(authHeader string) string {
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return ""
	}
	return parts[1]
}

func ParseToken(tokenstring string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte("114514"), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token, nil
}

func ClaimToken(token *jwt.Token) string {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "无效的token claims"
	}
	if exp, ok := claims["exp"]; ok {
		if float64(time.Now().Unix()) > exp.(float64) {
			return "claims已过期"
		}
	}

	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	return ""
}

func NewToken() (string, string) {
	newToken := jwt.New(jwt.SigningMethodHS256)
	secret := []byte("114514")
	signedToken, err := newToken.SignedString(secret)
	if err != nil {
		return fmt.Sprintf("生成 JWT 签名出错: %v", err), ""
	}
	return "", signedToken
}

func IsLogin() gin.HandlerFunc {
	return func(client *gin.Context) {
		//读取头文件
		authHeader := client.Request.Header.Get("Authorization")
		if authHeader == "" {
			SendErrorResponse(
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
			SendErrorResponse(
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
			SendErrorResponse(
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
