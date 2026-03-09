package middleware

import (
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/pkg/response"
	"github.com/Rezarit/E-commerce/service"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

// LoginRequired 登录态验证中间件
func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取Authorization请求头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("[登录态验证失败] 请求头Authorization为空 | 请求路径：%s", c.FullPath())
			response.Fail(c, domain.ErrCodeTokenEmpty, "请先登录")
			c.Abort()
			return
		}

		// 拆分Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && strings.EqualFold(parts[0], "Bearer")) {
			log.Printf("[登录态验证失败] Authorization格式错误 | 请求路径：%s | 错误值：%s", c.FullPath(), authHeader)
			response.Fail(c, domain.ErrCodeTokenFormatError, "登录态格式错误，请重新登录")
			c.Abort()
			return
		}
		tokenStr := parts[1] // 提取纯Token字符串

		// 解析Token
		claims, err := service.ParseAccessToken(tokenStr)
		if err != nil {
			switch err.Error() {
			case "token_expired":
				log.Printf("[登录态验证失败] Token已过期 | 请求路径：%s | Token：%s", c.FullPath(), tokenStr)
				response.Fail(c, domain.ErrCodeTokenExpired, "登录态已过期，请重新登录")
			case "token_signature_invalid":
				log.Printf("[登录态验证失败] Token签名错误 | 请求路径：%s | Token：%s", c.FullPath(), tokenStr)
				response.Fail(c, domain.ErrCodeTokenInvalid, "登录态无效，请重新登录")
			default:
				log.Printf("[登录态验证失败] Token解析失败 | 请求路径：%s | 错误：%v | Token：%s", c.FullPath(), err, tokenStr)
				response.Fail(c, domain.ErrCodeTokenInvalid, "登录态无效，请重新登录")
			}
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		log.Printf("[登录态验证成功] 用户ID：%d | 请求路径：%s", claims.UserID, c.FullPath())

		c.Next()
	}
}
