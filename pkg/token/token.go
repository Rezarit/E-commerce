package token

import (
	"errors"
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/pkg/config"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

func GetTokenFromAuthHeader(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("无效的请求头")
	}
	return parts[1], nil
}

// ParseToken 解析token
//func ParseToken(tokenStr string) ，(jwt.Token，error) {
func ParseToken(tokenStr string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			hmacMethod, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok || hmacMethod.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.Cfg.JWT.Secret), nil
		},
	)

	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			switch {
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return domain.ErrTokenExpired // 过期
			case ve.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				return jwt.ErrSignatureInvalid // 签名无效
			}
		}
		return err
	}

	// 验证 Token 整体有效性
	if !token.Valid {
		return domain.ErrTokenInvalid
	}

	return nil
}

func ValidateRefreshToken(claims *domain.RefreshTokenClaims) error {
	nowUnix := time.Now().Unix()
	if claims.ExpiresAt < nowUnix {
		return domain.ErrTokenExpired
	}
	return nil
}
