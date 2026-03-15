package middleware

import (
	"github.com/Rezarit/go-seckill-system/api"
	"github.com/Rezarit/go-seckill-system/api/auth"
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/Rezarit/go-seckill-system/pkg/response"
	"github.com/Rezarit/go-seckill-system/pkg/token"
	"github.com/Rezarit/go-seckill-system/service"
	"github.com/gin-gonic/gin"
	"log"
)

// LoginRequired 登录态验证中间件
func LoginRequired() gin.HandlerFunc {
	return func(client *gin.Context) {
		// 读取Authorization请求头
		authHeader := auth.GetAuthHeader(client)
		if authHeader == "" {
			client.Abort()
			return
		}

		// 拆分Bearer Token
		tokenStr, err := token.GetTokenFromAuthHeader(authHeader)
		if err != nil {
			response.Fail(client, domain.ErrCodeTokenFormatError, err.Error())
			client.Abort()
			return
		}

		// 解析Token
		claims, err := service.ParseAccessToken(tokenStr)
		if err != nil {
			switch err.Error() {
			case "token_expired":
				log.Printf("[登录态验证失败] Token已过期 | 请求路径：%s | Token：%s", client.FullPath(), tokenStr)
				response.Fail(client, domain.ErrCodeTokenExpired, "登录态已过期，请重新登录")
			case "token_signature_invalid":
				log.Printf("[登录态验证失败] Token签名错误 | 请求路径：%s | Token：%s", client.FullPath(), tokenStr)
				response.Fail(client, domain.ErrCodeTokenInvalid, "登录态无效，请重新登录")
			default:
				log.Printf("[登录态验证失败] Token解析失败 | 请求路径：%s | 错误：%v | Token：%s", client.FullPath(), err, tokenStr)
				response.Fail(client, domain.ErrCodeTokenInvalid, "登录态无效，请重新登录")
			}
			client.Abort()
			return
		}

		client.Set("user_id", claims.UserID)
		log.Printf("[登录态验证成功] 用户ID：%d | 请求路径：%s", claims.UserID, client.FullPath())
		client.Next()
	}
}

// MerchantRequired 商家权限验证中间件
func MerchantRequired() gin.HandlerFunc {
	return func(client *gin.Context) {
		//获取用户ID
		userID := api.ParseUserID(client)
		if userID == 0 {
			log.Printf("[商户权限验证失败] 获取用户信息失败 | 用户ID：%d", userID)
			response.Fail(client, domain.ErrCodeTokenInvalid, "登录态无效，请重新登录")
			return
		}

		info, err := service.GetUserInfoById(userID)
		if err != nil {
			log.Printf("[商户权限验证失败] 查询用户信息失败 | 用户ID：%d | 当前角色：%d", userID, info.Role)
			response.Fail(client, domain.ErrCodeDBError, err.Error())
			client.Abort()
			return
		}

		if info.Role != domain.RoleMerchant {
			log.Printf("[商户权限验证失败] 用户角色不符 | 用户ID：%d | 当前角色：%d | 所需角色：%d", userID, info.Role, domain.RoleMerchant)
			response.Fail(client, domain.ErrCodePermissionDenied, "您没有权限执行此操作")
			client.Abort()
			return
		}

		client.Set("user_id", userID)
		log.Printf("[登录态验证成功] 用户ID：%d | 请求路径：%s", userID, client.FullPath())
		client.Next()
	}
}
