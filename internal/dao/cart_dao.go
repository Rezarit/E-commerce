package dao

import (
	"github.com/Rezarit/go-seckill-system/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AddToCart 加入购物车（存在则累加数量，不存在则插入）
func AddToCart(item domain.CartItem) error {
	err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "product_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"quantity": gorm.Expr("quantity + ?", item.Quantity),
		}),
	}).Create(&item).Error
	if err != nil {
		return err
	}
	return nil
}

func ShowCart(userID int64) ([]domain.CartItem, error) {
	var items []domain.CartItem
	err := GetRecordsByField[domain.CartItem]("user_id", userID, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func RemoveFromCart(userID, productID int64) error {
	if err := DB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&domain.CartItem{}).Error; err != nil {
		return err
	}
	return nil
}

func ClearCart(userID int64) error {
	if err := DB.Where("user_id = ?", userID).Delete(&domain.CartItem{}).Error; err != nil {
		return err
	}
	return nil
}
