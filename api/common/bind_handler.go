package common

import (
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/Rezarit/go-seckill-system/pkg/response"
	"github.com/gin-gonic/gin"
	"log"
)

// BindRequest 绑定信息
func BindRequest(client *gin.Context, req interface{}) bool {
	if err := client.ShouldBindJSON(req); err != nil {
		log.Printf("[参数绑定失败] | 路径：%s | 错误：%v", client.FullPath(), err)
		response.Fail(client, domain.ErrCodeJSONParseFailed, "参数解析失败，请检查请求格式")
		return false
	}
	return true
}
