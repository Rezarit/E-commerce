package auth

import (
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/Rezarit/go-seckill-system/pkg/response"
	"github.com/gin-gonic/gin"
	"log"
)

func GetAuthHeader(client *gin.Context) string {
	authHeader := client.GetHeader("Authorization")
	if authHeader == "" {
		log.Printf("[登录态验证失败] 请求头Authorization为空 | 请求路径：%s", client.FullPath())
		response.Fail(client, domain.ErrCodeTokenEmpty, "请先登录")
		return ""
	}
	return authHeader
}
