package token

import (
	"errors"
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/Rezarit/go-seckill-system/pkg/config"
	"github.com/dgrijalva/jwt-go"
	"log"
	"strings"
	"time"
)

func GetTokenFromAuthHeader(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && strings.EqualFold(parts[0], "Bearer")) {
		log.Printf("[登录态验证失败] Authorization格式错误 | 错误值：%s", authHeader)
		return "", errors.New("无效的请求头")
	}
	return parts[1], nil
}

// ParseToken 解析token
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
