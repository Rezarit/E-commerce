package service

import (
	"github.com/Rezarit/go-seckill-system/dao"
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/Rezarit/go-seckill-system/pkg/validator"
	"log"
)

// CreatProduct 添加商品
func CreatProduct(product domain.ProductCreatRequest, userID int64) (int64, error) {
	// 检查商品名是否符合要求
	if err := CheckProductName(product.ProductName); err != nil {
		return 0, err
	}
	// 检查商品名是否存在
	if err := CheckProductNameExists(product.ProductName); err != nil {
		return 0, err
	}

	// 获取商户ID
	merchant, err := dao.GetMerchantByUserID(userID)
	if err != nil {
		log.Printf("[Service] 查询商户信息失败 | 用户ID：%d | 错误：%v", userID, err)
		return 0, &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "查询商户信息失败"}
	}
	// 整合商品信息
	productToInsert := domain.Product{
		MerchantID:  merchant.MerchantID,
		ProductName: product.ProductName,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Cover:       product.Cover,
		Link:        product.Link,
	}

	// 插入商品相关数据
	log.Printf("[Service] 开始添加商品 | 商品名：%s", product.ProductName)
	if err = dao.InsertProduct(&productToInsert); err != nil {
		log.Printf("[Service] 添加商品失败 | 商品名：%s | 错误：%v", product.ProductName, err)
		return 0, &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "添加商品失败"}
	}
	log.Printf("[Service] 添加商品成功 | 商品名：%s", product.ProductName)
	return productToInsert.ProductID, nil
}

// CheckProductName 检查商品名是否符合要求
func CheckProductName(productName string) error {
	trimmedName, err := validator.TrimAndCheckEmpty(productName, "商品名")
	if err != nil {
		log.Printf("[Service] 商品名不能为空 | 商品名：%s | 错误：%v", productName, err)
		return &domain.BusinessError{Code: domain.ErrCodeParamInvalid, Msg: err.Error()}
	}
	if err = validator.CheckLengthRange(trimmedName, "商品名", 1, 20); err != nil {
		log.Printf("[Service] 商品名长度需在1-20位之间 | 商品名：%s | 错误：%v", productName, err)
		return &domain.BusinessError{Code: domain.ErrCodeParamInvalid, Msg: err.Error()}
	}
	return nil
}

// CheckProductNameExists 检查商品名是否存在
func CheckProductNameExists(productName string) error {
	//查询商品名是否存在
	exists, err := dao.CheckProductNameExists(productName)
	if err != nil {
		log.Printf("[Service] 检查商品名是否存在失败 | 商品名：%s | 错误：%v", productName, err)
		return &domain.BusinessError{Code: domain.ErrCodeParamInvalid, Msg: err.Error()}
	}
	if exists {
		log.Printf("[Service] 商品名已存在 | 商品名：%s", productName)
		return &domain.BusinessError{Code: domain.ErrCodeProductExists, Msg: "商品名已存在"}
	}
	return nil
}

// UpdateProduct 更新商品
func UpdateProduct(product domain.ProductUpdateRequest, userID int64) error {
	// 检查商品是否存在
	if err := CheckProductIDExists(product.ProductID); err != nil {
		return err
	}
	// 检查商品归属权
	if err := CheckProductOwnership(product.ProductID, userID); err != nil {
		return err
	}
	// 检查商品名是否符合要求
	if err := CheckProductName(product.ProductName); err != nil {
		return err
	}

	productToUpdate := domain.Product{
		ProductName: product.ProductName,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Cover:       product.Cover,
		Link:        product.Link,
	}

	// 更新商品相关数据
	log.Printf("[Service] 开始更新商品 | 商品ID：%d", product.ProductID)
	if err := dao.UpdateProduct(&productToUpdate); err != nil {
		log.Printf("[Service] 更新商品失败 | 商品ID：%d | 错误：%v", product.ProductID, err)
		return &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "更新商品失败"}
	}
	log.Printf("[Service] 更新商品成功 | 商品ID：%d", product.ProductID)
	return nil
}

// CheckProductIDExists 检查商品ID是否存在
func CheckProductIDExists(productID int64) error {
	exists, err := dao.CheckProductIDExists(productID)
	if err != nil {
		log.Printf("[Service] 检查商品ID是否存在失败 | 商品ID：%d | 错误：%v", productID, err)
		return &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "检查商品ID失败"}
	}
	if !exists {
		log.Printf("[Service] 商品ID不存在 | 商品ID：%d", productID)
		return &domain.BusinessError{Code: domain.ErrCodeProductNotFound, Msg: "商品ID不存在"}
	}
	return nil
}

// DeleteProduct 删除商品
func DeleteProduct(productID int64, userID int64) error {
	// 检查商品是否存在
	if err := CheckProductIDExists(productID); err != nil {
		return err
	}
	// 检查商品归属权
	if err := CheckProductOwnership(productID, userID); err != nil {
		return err
	}

	// 删除商品相关数据
	log.Printf("[Service] 开始删除商品 | 商品ID：%d | 商户ID：%d", productID, userID)
	if err := dao.DeleteProduct(productID); err != nil {
		log.Printf("[Service] 删除商品失败 | 商品ID：%d | 错误：%v", productID, err)
		return &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "删除商品失败"}
	}
	log.Printf("[Service] 删除商品成功 | 商品ID：%d", productID)
	return nil
}

// CheckProductOwnership 检查商品归属权
func CheckProductOwnership(productID int64, userID int64) error {
	log.Printf("[Service] 开始检查商品归属权 | 商品ID：%d | 用户ID：%d", productID, userID)
	product, err := dao.GetProductByID(productID)
	if err != nil {
		log.Printf("[Service] 查询商品失败 | 商品ID：%d | 错误：%v", productID, err)
		return &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "查询商品失败"}
	}

	merchant, err := dao.GetMerchantByUserID(userID)
	if err != nil {
		return &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "查询商户信息失败"}
	}

	if product.MerchantID != merchant.MerchantID {
		log.Printf("[Service] 商品归属权错误 | 商品ID：%d | 商户ID：%d | 操作商户：%d", productID, product.MerchantID, merchant.MerchantID)
		return &domain.BusinessError{
			Code: domain.ErrCodePermissionDenied,
			Msg:  "无权操作此商品",
		}
	}
	log.Printf("[Service] 商品归属权验证成功 | 商品ID：%d | 用户ID：%d", productID, userID)
	return nil
}

// GetProductList 获取商品列表
func GetProductList() ([]domain.Product, error) {
	log.Printf("[Service] 开始获取商品列表")
	products, err := dao.GetProductList()
	if err != nil {
		log.Printf("[Service] 获取商品列表失败 | 错误：%v", err)
		return nil, &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "获取商品列表失败"}
	}
	log.Printf("[Service] 获取商品列表成功 | 商品数量：%d", len(products))
	return products, nil
}

// SearchProduct 搜索商品
func SearchProduct(keyword string) ([]domain.ProductSearchResponse, error) {
	log.Printf("[Service] 开始搜索商品 | 关键词：%s", keyword)
	products, err := dao.SearchProduct(keyword)
	if err != nil {
		log.Printf("[Service] 搜索商品失败 | 关键词：%s | 错误：%v", keyword, err)
		return nil, &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "搜索商品失败"}
	}

	// 转换为响应格式
	var resp []domain.ProductSearchResponse
	for _, product := range products {
		resp = append(resp, domain.ProductSearchResponse{
			ProductID:   product.ProductID,
			ProductName: product.ProductName,
			Price:       product.Price,
			Cover:       product.Cover,
		})
	}

	log.Printf("[Service] 搜索商品成功 | 关键词：%s | 商品数量：%d", keyword, len(products))
	return resp, nil
}

// GetProductDetail 获取商品详情
func GetProductDetail(productID int64) (*domain.Product, error) {
	log.Printf("[Service] 开始获取商品详情 | 商品ID：%d", productID)
	// 检查商品是否存在
	if err := CheckProductIDExists(productID); err != nil {
		return nil, err
	}

	product, err := dao.GetProductByID(productID)
	if err != nil {
		log.Printf("[Service] 获取商品详情失败 | 商品ID：%d | 错误：%v", productID, err)
		return nil, &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "获取商品详情失败"}
	}
	log.Printf("[Service] 获取商品详情成功 | 商品ID：%d", productID)
	return product, nil
}
