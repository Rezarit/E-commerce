package response

import (
	"github.com/Rezarit/E-commerce/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, domain.Response{
		Code: domain.ErrCodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, domain.Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
