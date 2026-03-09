package common

import (
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/pkg/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
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

// ParseParam 通用解析参数
func ParseParam(param string, client *gin.Context) (interface{}, bool) {
	paramStr := client.Param(param)
	paramValue, err := url.ParseQuery(paramStr)
	if err != nil {
		log.Printf("[参数解析失败] | 参数：%s | 错误：%v", client.FullPath(), err)
		response.Fail(client, domain.ErrCodeJSONParseFailed, "参数解析失败，请检查请求格式")
		return nil, false
	}
	return paramValue, true
}
