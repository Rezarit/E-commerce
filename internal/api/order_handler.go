package api

import (
	common2 "github.com/Rezarit/go-seckill-system/internal/api/common"
	"github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/internal/service"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"github.com/Rezarit/go-seckill-system/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// MakeOrder 下单
func MakeOrder(client *gin.Context) {
	// 获取用户ID
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	// 绑定请求参数
	var req domain.OrderCreateRequest
	isPass := common2.BindRequest(client, &req)
	if !isPass {
		return
	}

	// 执行下单操作
	acc, err := service.MakeOrder(userID, req.Address)
	if !common2.HandleBusinessError(client, err) {
		return
	}

	// 返回受理凭证：msg_id + 冻结的商品清单（前端/压测拿 msg_id 轮询 /order/result）
	response.Success(client, "下单受理成功，请用 msg_id 查询处理结果", acc)
}

// GetOrderList 获取订单列表
func GetOrderList(client *gin.Context) {
	// 获取用户ID
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	// 执行获取订单列表操作
	orders, err := service.GetOrderList(userID)
	if !common2.HandleBusinessError(client, err) {
		return
	}

	// 返回成功响应
	response.Success(client, "获取订单列表成功", orders)
}

// ParseOrderID 获取订单ID
func ParseOrderID(client *gin.Context) int64 {
	orderIDStr := client.Param("order_id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		logger.Sugar.Errorf("[API] 解析订单ID失败 | 错误：%v", err)
		return 0
	}
	return orderID
}

// GetOrderDetail 获取订单详情
func GetOrderDetail(client *gin.Context) {
	// 获取用户ID
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	// 获取订单ID
	orderID := ParseOrderID(client)
	if orderID == 0 {
		return
	}

	// 执行获取订单详情操作
	order, orderItems, err := service.GetOrderDetail(orderID, userID)
	if !common2.HandleBusinessError(client, err) {
		return
	}

	// 返回成功响应
	response.Success(client,
		"获取订单详情成功",
		gin.H{
			"order":       order,
			"order_items": orderItems,
		})
}

// QueryOrderResult 查询下单处理结果（三态：processing / success / failed）
// 前端/压测拿下单返回的 msg_id 轮询此接口，直到 success（成交）或 failed（最终失败）
func QueryOrderResult(client *gin.Context) {
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	msgID := client.Query("msg_id")
	if msgID == "" {
		common2.HandleBusinessError(client, &domain.BusinessError{
			Code: domain.ErrCodeParamInvalid,
			Msg:  "缺少 msg_id 参数",
		})
		return
	}

	result, err := service.QueryOrderResult(msgID, userID)
	if !common2.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "查询成功", result)
}
