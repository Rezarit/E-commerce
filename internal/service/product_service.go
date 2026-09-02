package service

import (
	"math/rand"
	"time"

	dao2 "github.com/Rezarit/go-seckill-system/internal/dao"
	domain2 "github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"github.com/Rezarit/go-seckill-system/pkg/redis"
	"github.com/Rezarit/go-seckill-system/pkg/validator"
)

// CreatProduct 添加商品
func CreatProduct(product domain2.ProductCreatRequest, userID int64) (int64, error) {
	// 检查商品名是否符合要求
	if err := CheckProductName(product.ProductName); err != nil {
		return 0, err
	}
	// 检查商品名是否存在
	if err := CheckProductNameExists(product.ProductName); err != nil {
		return 0, err
	}

	// 获取商户ID
	merchant, err := dao2.GetMerchantByUserID(userID)
	if err != nil {
		logger.Sugar.Errorf("[Service] 查询商户信息失败 | 用户ID：%d | 错误：%v", userID, err)
		return 0, &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "查询商户信息失败"}
	}
	// 整合商品信息
	productToInsert := domain2.Product{
		MerchantID:  merchant.MerchantID,
		ProductName: product.ProductName,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Cover:       product.Cover,
		Link:        product.Link,
	}

	// 插入商品相关数据
	logger.Sugar.Infof("[Service] 开始添加商品 | 商品名：%s", product.ProductName)
	if err = dao2.InsertProduct(&productToInsert); err != nil {
		logger.Sugar.Errorf("[Service] 添加商品失败 | 商品名：%s | 错误：%v", product.ProductName, err)
		return 0, &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "添加商品失败"}
	}
	logger.Sugar.Infof("[Service] 添加商品成功 | 商品名：%s", product.ProductName)
	return productToInsert.ProductID, nil
}

// CheckProductName 检查商品名是否符合要求
func CheckProductName(productName string) error {
	trimmedName, err := validator.TrimAndCheckEmpty(productName, "商品名")
	if err != nil {
		logger.Sugar.Errorf("[Service] 商品名不能为空 | 商品名：%s | 错误：%v", productName, err)
		return &domain2.BusinessError{Code: domain2.ErrCodeParamInvalid, Msg: err.Error()}
	}
	if err = validator.CheckLengthRange(trimmedName, "商品名", 1, 200); err != nil {
		logger.Sugar.Errorf("[Service] 商品名长度需在1-200字节之间 | 商品名：%s | 错误：%v", productName, err)
		return &domain2.BusinessError{Code: domain2.ErrCodeParamInvalid, Msg: err.Error()}
	}
	return nil
}

// CheckProductNameExists 检查商品名是否存在
func CheckProductNameExists(productName string) error {
	//查询商品名是否存在
	exists, err := dao2.CheckProductNameExists(productName)
	if err != nil {
		logger.Sugar.Errorf("[Service] 检查商品名是否存在失败 | 商品名：%s | 错误：%v", productName, err)
		return &domain2.BusinessError{Code: domain2.ErrCodeParamInvalid, Msg: err.Error()}
	}
	if exists {
		logger.Sugar.Infof("[Service] 商品名已存在 | 商品名：%s", productName)
		return &domain2.BusinessError{Code: domain2.ErrCodeProductExists, Msg: "商品名已存在"}
	}
	return nil
}

// UpdateProduct 更新商品
func UpdateProduct(productID int64, product domain2.ProductUpdateRequest, userID int64) error {
	// 检查商品归属权
	if err := CheckProductOwnership(productID, userID); err != nil {
		return err
	}
	// 检查商品名是否符合要求
	if err := CheckProductName(product.ProductName); err != nil {
		return err
	}

	productToUpdate := domain2.Product{
		ProductName: product.ProductName,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Cover:       product.Cover,
		Link:        product.Link,
	}

	// 更新商品相关数据
	logger.Sugar.Infof("[Service] 开始更新商品 | 商品ID：%d", productID)
	if err := dao2.UpdateProduct(&productToUpdate); err != nil {
		logger.Sugar.Errorf("[Service] 更新商品失败 | 商品ID：%d | 错误：%v", productID, err)
		return &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "更新商品失败"}
	}
	logger.Sugar.Infof("[Service] 更新商品成功 | 商品ID：%d", productID)
	return nil
}

// DeleteProduct 删除商品
func DeleteProduct(productID int64, userID int64) error {
	// 检查商品归属权
	if err := CheckProductOwnership(productID, userID); err != nil {
		return err
	}

	// 删除商品相关数据
	logger.Sugar.Infof("[Service] 开始删除商品 | 商品ID：%d | 商户ID：%d", productID, userID)
	if err := dao2.DeleteProduct(productID); err != nil {
		logger.Sugar.Errorf("[Service] 删除商品失败 | 商品ID：%d | 错误：%v", productID, err)
		return &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "删除商品失败"}
	}
	logger.Sugar.Infof("[Service] 删除商品成功 | 商品ID：%d", productID)
	return nil
}

// CheckProductOwnership 检查商品归属权
func CheckProductOwnership(productID int64, userID int64) error {
	logger.Sugar.Infof("[Service] 开始检查商品归属权 | 商品ID：%d | 用户ID：%d", productID, userID)
	product, err := dao2.GetProductByID(productID)
	if err != nil {
		logger.Sugar.Errorf("[Service] 查询商品失败 | 商品ID：%d | 错误：%v", productID, err)
		return &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "查询商品失败"}
	}

	merchant, err := dao2.GetMerchantByUserID(userID)
	if err != nil {
		return &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "查询商户信息失败"}
	}

	if product.MerchantID != merchant.MerchantID {
		logger.Sugar.Errorf("[Service] 商品归属权错误 | 商品ID：%d | 商户ID：%d | 操作商户：%d", productID, product.MerchantID, merchant.MerchantID)
		return &domain2.BusinessError{
			Code: domain2.ErrCodePermissionDenied,
			Msg:  "无权操作此商品",
		}
	}
	logger.Sugar.Infof("[Service] 商品归属权验证成功 | 商品ID：%d | 用户ID：%d", productID, userID)
	return nil
}

// RefreshStockFromDB 刷新式补偿：查 MySQL 库存并覆盖写回 Redis
// （死信队列处理用：死信=订单没建成=MySQL 库存正确，让 Redis 对齐 MySQL，
//   以 MySQL 为权威，天然幂等、不会超卖、不信任消息内容）
func RefreshStockFromDB(productID int64) error {
	product, err := dao2.GetProductByID(productID)
	if err != nil {
		return err
	}
	if err := cacheService.CacheProductStock(productID, int32(product.Stock), redis.DefaultSessionTTL); err != nil {
		return err
	}
	logger.Sugar.Infof("[Service] 库存刷新补偿 | 商品ID: %d | 库存: %d", productID, product.Stock)
	return nil
}

// GetProductList 获取商品列表
func GetProductList() ([]domain2.Product, error) {
	logger.Sugar.Infof("[Service] 开始获取商品列表")
	products, err := dao2.GetProductList()
	if err != nil {
		logger.Sugar.Errorf("[Service] 获取商品列表失败 | 错误：%v", err)
		return nil, &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "获取商品列表失败"}
	}
	logger.Sugar.Infof("[Service] 获取商品列表成功 | 商品数量：%d", len(products))
	return products, nil
}

// SearchProduct 搜索商品
func SearchProduct(keyword string) ([]domain2.ProductSearchResponse, error) {
	logger.Sugar.Infof("[Service] 开始搜索商品 | 关键词：%s", keyword)
	products, err := dao2.SearchProduct(keyword)
	if err != nil {
		logger.Sugar.Errorf("[Service] 搜索商品失败 | 关键词：%s | 错误：%v", keyword, err)
		return nil, &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "搜索商品失败"}
	}

	// 转换为响应格式
	var resp []domain2.ProductSearchResponse
	for _, product := range products {
		resp = append(resp, domain2.ProductSearchResponse{
			ProductID:   product.ProductID,
			ProductName: product.ProductName,
			Price:       product.Price,
			Cover:       product.Cover,
		})
	}

	logger.Sugar.Infof("[Service] 搜索商品成功 | 关键词：%s | 商品数量：%d", keyword, len(products))
	return resp, nil
}

// GetProductDetail 获取商品详情（缓存 + 击穿锁 + 穿透防护）
func GetProductDetail(productID int64) (*domain2.Product, error) {
	// 空值缓存命中：商品不存在，直接返回（防穿透）
	if isNull, err := cacheService.GetNullProduct(productID); err == nil && isNull {
		logger.Sugar.Infof("[Service] 命中空值缓存 | 商品ID：%d", productID)
		return nil, &domain2.BusinessError{
			Code: domain2.ErrCodeProductNotFound,
			Msg:  "商品不存在",
		}
	}

	// 从缓存获取商品详情（found=true 命中）
	product, found, err := cacheService.GetProductFromCache(productID)
	if err == nil && found {
		logger.Sugar.Infof("[Service] 从缓存获取商品详情成功 | 商品ID：%d", productID)
		return product, nil
	}

	// 缓存未命中：尝试获取重建锁（防击穿——热点 key 过期瞬间只有一个请求查 DB）
	locked, lockErr := cacheService.TryLockProduct(productID)
	if lockErr != nil {
		logger.Sugar.Errorf("[Service] 获取商品锁失败 | 商品ID：%d | 错误：%v", productID, lockErr)
	}
	if locked {
		defer func() {
			if err := cacheService.UnlockProduct(productID); err != nil {
				logger.Sugar.Errorf("[Service] 释放商品锁失败 | 商品ID：%d | 错误：%v", productID, err)
			}
		}()
		// 抢到锁：从数据库获取
		product, err = dao2.GetProductByID(productID)
		if err != nil {
			// 商品不存在：缓存空值防穿透
			logger.Sugar.Infof("[Service] 商品不存在，缓存空值 | 商品ID：%d", productID)
			_ = cacheService.CacheNullProduct(productID, redis.DefaultNullCacheTTL)
			return nil, &domain2.BusinessError{
				Code: domain2.ErrCodeProductNotFound,
				Msg:  "商品不存在",
			}
		}

		// 缓存商品详情（TTL 加随机偏移，防雪崩）
		ttl := cacheTTLWithJitter(redis.DefaultProductCacheTTL)
		if cacheErr := cacheService.CacheProduct(product, ttl); cacheErr != nil {
			logger.Sugar.Errorf("[Service] 缓存商品详情失败 | 商品ID：%d | 错误：%v", productID, cacheErr)
		}

		logger.Sugar.Infof("[Service] 获取商品详情成功 | 商品ID：%d", productID)
		return product, nil
	}

	// 没抢到锁：说明别的请求正在查库回填，短暂等待后重查缓存
	for i := 0; i < 3; i++ {
		time.Sleep(50 * time.Millisecond)
		product, found, err = cacheService.GetProductFromCache(productID)
		if err == nil && found {
			return product, nil
		}
		// 若对方查到的是空值，GetNullProduct 也检查下
		if isNull, _ := cacheService.GetNullProduct(productID); isNull {
			return nil, &domain2.BusinessError{
				Code: domain2.ErrCodeProductNotFound,
				Msg:  "商品不存在",
			}
		}
	}

	logger.Sugar.Errorf("[Service] 等待缓存重建超时 | 商品ID：%d", productID)
	return nil, &domain2.BusinessError{
		Code: domain2.ErrCodeDBError,
		Msg:  "系统繁忙，请稍后再试",
	}
}

// cacheTTLWithJitter 给 TTL 加随机偏移（±10%），防止大量 key 同时过期导致雪崩
func cacheTTLWithJitter(base time.Duration) time.Duration {
	// 随机 -10% ~ +10%
	offset := time.Duration(rand.Intn(20)-10) * base / 100
	return base + offset
}

func GetMerchantProductList(userID int64) ([]domain2.Product, error) {
	logger.Sugar.Infof("[Service] 开始获取商户商品列表 | 用户ID：%d", userID)
	merchantID, err := dao2.GetMerchantIDByUserID(userID)

	products, err := dao2.GetProductListByMerchantID(merchantID)
	if err != nil {
		logger.Sugar.Errorf("[Service] 获取商户商品列表失败 | 商户ID：%d | 错误：%v", merchantID, err)
		return nil, &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "获取商户商品列表失败"}
	}
	logger.Sugar.Infof("[Service] 获取商户商品列表成功 | 商户ID：%d | 商品数量：%d", merchantID, len(products))
	return products, nil
}
