package service

import (
	"errors"
	"fmt"
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/pkg/config"
	"github.com/Rezarit/E-commerce/pkg/token"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

func getJWTSecret() []byte {
	return []byte(config.Cfg.JWT.Secret)
}

// ParseAccessToken 解析Access Token
func ParseAccessToken(tokenStr string) (*domain.CustomClaims, error) {
	claims := &domain.CustomClaims{}
	err := token.ParseToken(tokenStr, claims)
	if err != nil {
		tokenPreview := tokenStr
		if len(tokenStr) > 10 {
			tokenPreview = tokenStr[:10]
		}
		log.Printf("[解析AccessToken失败] Token：%s，错误：%v", tokenPreview, err)
		return nil, fmt.Errorf("access token解析失败：%w", err)
	}
	return claims, nil
}

// ParseRefreshToken 解析Refresh Token
func ParseRefreshToken(tokenStr string) (*domain.RefreshTokenClaims, error) {
	claims := &domain.RefreshTokenClaims{}
	err := token.ParseToken(tokenStr, claims)
	if err != nil {
		tokenPreview := tokenStr
		if len(tokenStr) > 10 {
			tokenPreview = tokenStr[:10]
		}
		log.Printf("[解析RefreshToken失败] Token：%s，错误：%v", tokenPreview, err)
		return nil, errors.New("refresh token解析失败：" + err.Error())
	}
	return claims, nil
}

func GenerateAccessToken(userID int64) (string, error) {
	expireHour := config.Cfg.JWT.AccessTokenExpireHour
	claims := domain.CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expireHour) * time.Hour).Unix(),
			Issuer:    "e-commerce",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString(getJWTSecret())
	if err != nil {
		log.Printf("[AccessToken生成失败] 用户ID：%d，错误：%v", userID, err)
		return "", &domain.BusinessError{
			Code: domain.ErrCodeSystemError,
			Msg:  "登录态生成失败，请重试",
		}
	}
	return tokenString, nil
}

func GenerateRefreshToken(userID int64) (string, error) {
	expireHour := config.Cfg.JWT.RefreshTokenExpireHour
	claims := domain.RefreshTokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expireHour) * time.Hour).Unix(),
			Issuer:    "e-commerce",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := refreshToken.SignedString(getJWTSecret())
	if err != nil {
		log.Printf("[RefreshToken生成失败] 用户ID：%d，错误：%v", userID, err)
		return "", &domain.BusinessError{
			Code: domain.ErrCodeSystemError,
			Msg:  "刷新登录态失败",
		}
	}
	return tokenString, nil
}
