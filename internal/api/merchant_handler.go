package api

import (
	common2 "github.com/Rezarit/go-seckill-system/internal/api/common"
	"github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/internal/service"
	"github.com/Rezarit/go-seckill-system/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterMerchant(client *gin.Context) {
	// 获取用户ID
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	// 绑定请求
	var req domain.MerchantApplyRequest
	isPass := common2.BindRequest(client, &req)
	if !isPass {
		return
	}

	// 执行注册指令
	err := service.RegisterMerchant(req, userID)
	if !common2.HandleBusinessError(client, err) {
		return
	}

	// 成功响应
	resp := domain.MerchantApplyResponse{
		Status: domain.MerchantStatusPending,
	}
	response.Success(client, "商户注册申请已提交", resp)
}
