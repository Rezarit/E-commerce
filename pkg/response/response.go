package response

import (
	domain2 "github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, domain2.Response{
		Code: domain2.ErrCodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, domain2.Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
