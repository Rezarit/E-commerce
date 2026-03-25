package dao

import (
	"errors"
	domain2 "github.com/Rezarit/go-seckill-system/internal/domain"
	"gorm.io/gorm"
)

// CheckProductNameExists 商品名是否存在
func CheckProductNameExists(productName string) (bool, error) {
	return CheckFieldExists[domain2.Product, string]("product_name", productName)
}

// CheckProductIDExists 商品ID是否存在
func CheckProductIDExists(productID int64) (bool, error) {
	return CheckFieldExists[domain2.Product, int64]("product_id", productID)
}

// InsertProduct 商品相关数据插入数据库
func InsertProduct(product *domain2.Product) error {
	if err := InsertRecord(product); err != nil {
		return err
	}
	return nil
}

// UpdateProduct 更新商品相关数据
func UpdateProduct(product *domain2.Product) error {
	if err := UpdateRecord("product_id", product.ProductID, product); err != nil {
		return err
	}
	return nil
}

// GetProductByID 根据商品ID查询商品
func GetProductByID(productID int64) (*domain2.Product, error) {
	var product domain2.Product
	if err := GetRecordByField("product_id", productID, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

// DeleteProduct 删除商品
func DeleteProduct(productID int64) error {
	if err := DeleteRecord[domain2.Product]("product_id", productID); err != nil {
		return err
	}
	return nil
}

// GetProductList 获取商品列表
func GetProductList() ([]domain2.Product, error) {
	var products []domain2.Product
	if err := DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// SearchProduct 搜索商品
func SearchProduct(keyword string) ([]domain2.Product, error) {
	var products []domain2.Product
	if err := DB.Where("product_name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductListByMerchantID 根据商户ID获取商品列表
func GetProductListByMerchantID(merchantID int64) ([]domain2.Product, error) {
	var products []domain2.Product
	err := GetRecordsByField[domain2.Product, int64]("merchant_id", merchantID, &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// GetMerchantIDByUserID 根据用户ID获取商户ID
func GetMerchantIDByUserID(userID int64) (int64, error) {
	var merchant domain2.Merchant
	err := GetRecordByField[domain2.Merchant, int64]("user_id", userID, &merchant)
	if err != nil {
		return 0, err
	}
	return merchant.MerchantID, nil
}

// DeductStock 扣减商品库存
func DeductStock(tx *gorm.DB, productID int64, quantity int) error {
	result := tx.Model(&domain2.Product{}).
		Where("product_id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))

	// 检查是否有行受到影响
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("stock not enough")
	}
	return nil
}
