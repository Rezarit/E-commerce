package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type Product struct {
	ProductID   int64           `json:"product_id" gorm:"primaryKey;autoIncrement"`
	MerchantID  int64           `json:"merchant_id" gorm:"not null"`
	ProductName string          `json:"product_name" gorm:"not null"`
	Description string          `json:"description" gorm:"not null"`
	CommentNum  int             `json:"comment_num" gorm:"default:0"`
	Price       decimal.Decimal `json:"price" gorm:"type:decimal(10,2);default:0"`
	Stock       int             `json:"stock" gorm:"default:0"`
	Cover       string          `json:"cover" gorm:"not null"`
	PublishTime time.Time       `json:"publish_time" gorm:"autoCreateTime"`
	Link        string          `json:"link" gorm:"not null"`
}

type ProductCreatRequest struct {
	ProductName string          `json:"product_name" binding:"required"`
	Description string          `json:"description" `
	Price       decimal.Decimal `json:"price" `
	Stock       int             `json:"stock" `
	Cover       string          `json:"cover"`
	Link        string          `json:"link"`
}

type ProductUpdateRequest struct {
	ProductID   int64           `json:"product_id" binding:"required"`
	ProductName string          `json:"product_name" binding:"required"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price" `
	Stock       int             `json:"stock"`
	Cover       string          `json:"cover"`
	Link        string          `json:"link"`
}

type ProductCreateResponse struct {
	ProductID int64 `json:"product_id"`
}

type ProductDeleteRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
}

type ProductSearchRequest struct {
	Keyword string `json:"keyword"`
}

type ProductSearchResponse struct {
	ProductID   int64           `json:"product_id"`
	ProductName string          `json:"product_name"`
	Price       decimal.Decimal `json:"price" `
	Cover       string          `json:"cover"`
}
