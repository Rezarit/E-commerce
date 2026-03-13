package api

import (
	"github.com/Rezarit/E-commerce/api/common"
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/pkg/response"
	"github.com/Rezarit/E-commerce/service"
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
	isPass := common.BindRequest(client, &req)
	if !isPass {
		return
	}

	// 执行注册指令
	err := service.RegisterMerchant(req, userID)
	if !common.HandleBusinessError(client, err) {
		return
	}

	// 成功响应
	resp := domain.MerchantApplyResponse{
		Status: domain.MerchantStatusPending,
	}
	response.Success(client, "商户注册申请已提交", resp)
}
