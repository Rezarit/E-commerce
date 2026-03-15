package api

import (
	"github.com/Rezarit/go-seckill-system/api/common"
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/Rezarit/go-seckill-system/pkg/response"
	"github.com/Rezarit/go-seckill-system/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CreatProduct 添加商品
func CreatProduct(client *gin.Context) {
	var product domain.ProductCreatRequest
	isPass := common.BindRequest(client, &product)
	if !isPass {
		return
	}

	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	productID, err := service.CreatProduct(product, userID)
	if !common.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "商品添加成功", domain.ProductCreateResponse{ProductID: productID})
}

// UpdateProduct 更新商品
func UpdateProduct(client *gin.Context) {
	var product domain.ProductUpdateRequest
	isPass := common.BindRequest(client, &product)
	if !isPass {
		return
	}

	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	err := service.UpdateProduct(product, userID)
	if !common.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "商品更新成功", product)
}

// DeleteProduct 删除商品
func DeleteProduct(client *gin.Context) {
	var product domain.ProductDeleteRequest
	isPass := common.BindRequest(client, &product)
	if !isPass {
		return
	}

	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	err := service.DeleteProduct(product.ProductID, userID)
	if !common.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "商品删除成功", nil)
}

// ShowProductList 获取商品列表
func ShowProductList(client *gin.Context) {
	products, err := service.GetProductList()
	if !common.HandleBusinessError(client, err) {
		return
	}
	response.Success(client, "获取商品列表成功", products)
}

// SearchProduct 搜索商品
func SearchProduct(client *gin.Context) {
	var searchRequest domain.ProductSearchRequest

	isPass := common.BindRequest(client, &searchRequest)
	if !isPass {
		return
	}

	products, err := service.SearchProduct(searchRequest.Keyword)
	if !common.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "搜索商品成功", products)
}

// ParseProductID 解析商品ID
func ParseProductID(client *gin.Context) int64 {
	productID := client.Param("product_id")
	if productID == "" {
		response.Fail(client, domain.ErrCodeParamInvalid, "商品ID不能为空")
		return 0
	}

	// 转换商品ID为int64
	pid, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		response.Fail(client, domain.ErrCodeParamInvalid, "商品ID格式错误")
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
	if !common.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "获取商品详情成功", product)
}
