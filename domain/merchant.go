package domain

import (
	"time"
)

// MerchantStatus 商户状态
type MerchantStatus string

const (
	MerchantStatusPending  = "pending"  // 待审核
	MerchantStatusActive   = "active"   // 已激活
	MerchantStatusRejected = "rejected" // 已拒绝
	MerchantStatusBanned   = "banned"   // 已封禁
)

// Merchant 商户信息表
type Merchant struct {
	MerchantID      int64          `json:"merchant_id" gorm:"primaryKey;autoIncrement"`
	UserID          int64          `json:"user_id" gorm:"uniqueIndex;not null"`
	MerchantName    string         `json:"merchant_name" gorm:"size:100;not null;uniqueIndex"`
	BusinessLicense string         `json:"business_license" gorm:"size:50;not null"`
	ContactPhone    string         `json:"contact_phone" gorm:"size:20;not null"`
	Address         string         `json:"address" gorm:"type:text;not null"`
	ShopDescription string         `json:"shop_description" gorm:"type:text"`
	Status          MerchantStatus `json:"status" gorm:"size:20;default:'active'"`
	CreateTime      time.Time      `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime      time.Time      `json:"update_time" gorm:"autoUpdateTime"`
}

type MerchantApplication struct { // 申请记录表
	ApplicationID   int64          `json:"application_id" gorm:"primaryKey;autoIncrement"`
	UserID          int64          `json:"user_id" gorm:"not null"`
	MerchantName    string         `json:"merchant_name" gorm:"size:100;not null"`
	BusinessLicense string         `json:"business_license" gorm:"size:50;not null"`
	ContactPhone    string         `json:"contact_phone" gorm:"size:20;not null"`
	Address         string         `json:"address" gorm:"type:text;not null"`
	Status          MerchantStatus `json:"status" gorm:"size:20;default:'pending'"`
	ApplyTime       time.Time      `json:"apply_time" gorm:"autoCreateTime"`
	AuditTime       *time.Time     `json:"audit_time"`
	AuditAdmin      *string        `json:"audit_admin" gorm:"size:50"`
	RejectReason    *string        `json:"reject_reason" gorm:"type:text"`
}

// MerchantApplyRequest 商户申请请求
type MerchantApplyRequest struct {
	MerchantName    string `json:"merchant_name" binding:"required"`
	BusinessLicense string `json:"business_license" binding:"required"`
	ContactPhone    string `json:"contact_phone" binding:"required"`
	Address         string `json:"address" binding:"required"`
}

// MerchantApplyResponse 商户申请响应
type MerchantApplyResponse struct {
	Status string `json:"status"` // 商户状态
}

// MerchantInfoResponse 商户信息响应
type MerchantInfoResponse struct {
	MerchantID      int64  `json:"merchant_id"`
	UserID          int64  `json:"user_id"`
	Username        string `json:"username"`
	MerchantName    string `json:"merchant_name"`
	BusinessLicense string `json:"business_license"`
	ContactPhone    string `json:"contact_phone"`
	Address         string `json:"address"`
	ShopDescription string `json:"shop_description"`
	Status          string `json:"status"`
	ApplyTime       string `json:"apply_time"`
	AuditTime       string `json:"audit_time,omitempty"`
}

// MerchantListResponse 商户列表项
type MerchantListResponse struct {
	MerchantID   int64  `json:"merchant_id"`
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	MerchantName string `json:"merchant_name"`
	Status       string `json:"status"`
	ApplyTime    string `json:"apply_time"`
	AuditTime    string `json:"audit_time,omitempty"`
}

// MerchantAuditRequest 商户审核请求
type MerchantAuditRequest struct {
	Status       string `json:"status" binding:"required"` // "active" 或 "rejected"
	RejectReason string `json:"reject_reason,omitempty"`   // 拒绝原因（当status=rejected时）
}

// MerchantUpdateRequest 商户信息更新请求
type MerchantUpdateRequest struct {
	MerchantName    string `json:"merchant_name,omitempty"`
	ContactPhone    string `json:"contact_phone,omitempty"`
	Address         string `json:"address,omitempty"`
	ShopDescription string `json:"shop_description,omitempty"`
}
