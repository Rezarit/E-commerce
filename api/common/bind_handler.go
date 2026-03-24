package common

import (
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/Rezarit/go-seckill-system/pkg/response"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
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

func ParamID(client *gin.Context, key string) int64 {
	IDStr := client.Param(key)
	ID, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		response.Fail(client, domain.ErrCodeParamInvalid, "无效的"+key)
		return 0
	}
	return ID
}
