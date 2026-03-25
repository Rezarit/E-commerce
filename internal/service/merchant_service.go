package service

import (
	"github.com/Rezarit/go-seckill-system/internal/dao"
	domain2 "github.com/Rezarit/go-seckill-system/internal/domain"
	"log"
	"strings"
)

func RegisterMerchant(merchant domain2.MerchantApplyRequest, userID int64) error {
	// 检查商户名是否存在
	err := CheckMerchantNameExists(merchant.MerchantName)
	if err != nil {
		return err
	}

	// 检查该用户是否已经是商户
	err = CheckUserIsMerchant(userID)
	if err != nil {
		return err
	}

	// 创建商户申请记录
	merchantRecord := &domain2.MerchantApplication{
		UserID:          userID,
		MerchantName:    merchant.MerchantName,
		BusinessLicense: merchant.BusinessLicense,
		ContactPhone:    merchant.ContactPhone,
		Address:         merchant.Address,
		Status:          domain2.MerchantStatusPending,
	}

	if err = dao.CreateMerchant(merchantRecord); err != nil {
		log.Printf("[Service] 创建商户申请失败 | 用户ID：%d | 错误：%v", userID, err)
		return err
	}

	// 记录申请成功日志
	log.Printf("[Service] 商户申请提交成功 | 用户ID：%d | 商户名：%s | 状态：待审核",
		userID, merchant.MerchantName)

	return nil
}

// CheckMerchantNameExists 检查商户名是否存在
func CheckMerchantNameExists(merchantName string) error {
	exists, err := dao.CheckMerchantNameExists(merchantName)
	if err != nil {
		log.Printf("[Service] 检查商户名存在性失败 | 商户名：%s | 错误：%v", merchantName, err)
		return err
	}
	if exists {
		log.Printf("[Service] 商户名已存在 | 商户名：%s", merchantName)
		return &domain2.BusinessError{
			Code: domain2.ErrCodeMerchantExists,
			Msg:  "商户名已存在",
		}
	}
	return nil
}

// CheckUserIsMerchant 检查该用户是否已经是商户
func CheckUserIsMerchant(userID int64) error {
	log.Printf("[Service] 检查用户是否已是商户 | 用户ID：%d", userID)
	existingMerchant, err := dao.GetMerchantByUserID(userID)
	if err != nil {
		// "record not found"错误，说明用户不是商户
		if strings.Contains(err.Error(), "record not found") {
			log.Printf("[Service] 用户不是商户 | 用户ID：%d", userID)
			return nil // 正常返回，用户可以申请商户
		}
		// 其他数据库错误
		log.Printf("[Service] 查询商户信息失败 | 用户ID：%d | 错误：%v", userID, err)
		return err
	}

	if existingMerchant.MerchantID != 0 {
		log.Printf("[Service] 用户已是商户 | 用户ID：%d | 商户名：%s", userID, existingMerchant.MerchantName)
		return &domain2.BusinessError{
			Code: domain2.ErrCodeAlreadyMerchant,
			Msg:  "该用户已经是商户",
		}
	}

	log.Printf("[Service] 用户不是商户 | 用户ID：%d", userID)
	return nil
}
