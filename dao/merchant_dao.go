package dao

import "github.com/Rezarit/go-seckill-system/domain"

// CheckMerchantNameExists 检查商户名是否已存在
func CheckMerchantNameExists(merchantName string) (bool, error) {
	exists, err := CheckFieldExists[domain.Merchant]("merchant_name", merchantName)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetMerchantByUserID 根据用户ID获取商户信息
func GetMerchantByUserID(userID int64) (domain.Merchant, error) {
	merchant := domain.Merchant{}
	err := GetRecordByField[domain.Merchant, int64]("user_id", userID, &merchant)
	if err != nil {
		return domain.Merchant{}, err
	}
	return merchant, nil
}

// CreateMerchant 创建商户记录
func CreateMerchant(merchantApply *domain.MerchantApplication) error {
	return InsertRecord[domain.MerchantApplication](merchantApply)
}
