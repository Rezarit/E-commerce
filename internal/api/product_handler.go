package api

import (
	common2 "github.com/Rezarit/go-seckill-system/internal/api/common"
	domain2 "github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/internal/service"
	"github.com/Rezarit/go-seckill-system/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CreatProduct 添加商品
func CreatProduct(client *gin.Context) {
	var product domain2.ProductCreatRequest
	isPass := common2.BindRequest(client, &product)
	if !isPass {
		return
	}

	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	productID, err := service.CreatProduct(product, userID)
	if !common2.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "商品添加成功", domain2.ProductCreateResponse{ProductID: productID})
}

// UpdateProduct 更新商品
func UpdateProduct(client *gin.Context) {
	productID := common2.ParamID(client, "product_id")
	if productID == 0 {
		return
	}

	var product domain2.ProductUpdateRequest
	isPass := common2.BindRequest(client, &product)
	if !isPass {
		return
	}

	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	err := service.UpdateProduct(productID, product, userID)
	if !common2.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "商品更新成功", product)
}

// DeleteProduct 删除商品
func DeleteProduct(client *gin.Context) {
	productID := common2.ParamID(client, "product_id")

	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	err := service.DeleteProduct(productID, userID)
	if !common2.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "商品删除成功", nil)
}

// ShowProductList 获取商品列表
func ShowProductList(client *gin.Context) {
	products, err := service.GetProductList()
	if !common2.HandleBusinessError(client, err) {
		return
	}
	response.Success(client, "获取商品列表成功", products)
}

// SearchProduct 搜索商品
func SearchProduct(client *gin.Context) {
	var searchRequest domain2.ProductSearchRequest

	isPass := common2.BindRequest(client, &searchRequest)
	if !isPass {
		return
	}

	products, err := service.SearchProduct(searchRequest.Keyword)
	if !common2.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "搜索商品成功", products)
}

// ParseProductID 解析商品ID
func ParseProductID(client *gin.Context) int64 {
	productID := client.Param("product_id")
	if productID == "" {
		response.Fail(client, domain2.ErrCodeParamInvalid, "商品ID不能为空")
		return 0
	}

	// 转换商品ID为int64
	pid, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		response.Fail(client, domain2.ErrCodeParamInvalid, "商品ID格式错误")
		return 0
	}
	return pid
}

// ProductDetail 获取商品详情
func ProductDetail(client *gin.Context) {
	productID := ParseProductID(client)
	if productID == 0 {
		return
	}

	product, err := service.GetProductDetail(productID)
	if !common2.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "获取商品详情成功", product)
}

func GetMerchantProductList(client *gin.Context) {
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	products, err := service.GetMerchantProductList(userID)
	if !common2.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "获取商品列表成功", products)
}
