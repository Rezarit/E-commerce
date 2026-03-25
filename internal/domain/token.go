package domain

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}
