package dao

import (
	"github.com/Rezarit/go-seckill-system/internal/domain"
)

func AddToCart(cart domain.Cart) error {
	err := InsertRecord[domain.Cart](&cart)
	if err != nil {
		return err
	}
	return nil
}

func ShowCart(userID int64) ([]domain.Cart, error) {
	var carts []domain.Cart
	err := GetRecordsByField[domain.Cart]("user_id", userID, &carts)
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func RemoveFromCart(userID, productID int64) error {
	if err := DB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&domain.Cart{}).Error; err != nil {
		return err
	}
	return nil
}

func ClearCart(userID int64) error {
	if err := DB.Where("user_id = ?", userID).Delete(&domain.Cart{}).Error; err != nil {
		return err
	}
	return nil
}

// CheckCartItemExists 检查购物车商品是否存在
func CheckCartItemExists(userID, productID int64) (bool, error) {
	var count int64
	if err := DB.Model(&domain.Cart{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
