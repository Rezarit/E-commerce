package common

import (
	"errors"
	"github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/pkg/response"
	"github.com/gin-gonic/gin"
)

func HandleBusinessError(client *gin.Context, err error) bool {
	if err == nil {
		return true
	}

	var bizErr *domain.BusinessError
	if !errors.As(err, &bizErr) {
		bizErr = &domain.BusinessError{
			Code: domain.ErrCodeUnknown,
			Msg:  "未知错误",
		}
	}
	response.Fail(client, bizErr.Code, bizErr.Msg)
	return false
}
