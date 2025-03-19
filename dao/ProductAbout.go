package dao

import (
	"database/sql"
	"github.com/Rezarit/E-commerce/domain"
)

func SearchProduct() (*sql.Rows, error) {
	cmd := "SELECT * FROM products"
	rows, err := DB.Query(cmd)
	if err != nil {
		return rows, err
	}
	return rows, nil
}

func SearchProductbyID(productID int) (domain.Product, error) {
	var product domain.Product
	//执行查询指令
	cmd := "SELECT * FROM products WHERE product_id=?"
	err := DB.QueryRow(cmd, productID).Scan(
		&product.ProductID,
		&product.Name,
		&product.Description,
		&product.Type,
		&product.CommentNum,
		&product.Price,
		&product.Cover,
		&product.PublishTime,
		&product.Link)
	if err != nil {
		return product, err
	}
	return product, nil
}
func SearchProductbyType(productType string) (*sql.Rows, error) {
	//执行查询指令
	cmd := "SELECT * FROM product WHERE type=?"
	rows, err := DB.Query(cmd, productType)
	if err != nil {
		return rows, err
	}
	defer func() { _ = rows.Close() }() //此处忽略了错误信息

	return rows, nil
}
