package service

import (
	"database/sql"
	"fmt"
	"github.com/Rezarit/E-commerce/dao"
	"github.com/Rezarit/E-commerce/domain"
	"strings"
)

func CartSearch(rows *sql.Rows) ([]int, error) {
	var products []domain.Product
	var productIDs []int

	for rows.Next() {
		var product domain.Product
		err := rows.Scan(&product.ProductID)
		if err != nil {
			return productIDs, err
		}
		products = append(products, product)
		productIDs = append(productIDs, product.ProductID)
	}

	err := rows.Err()
	if err != nil {
		return productIDs, err
	}
	return productIDs, nil
}

func SearchCart(productIDs []int, args []any) ([]domain.Product, error) {
	placeholder := make([]string, len(productIDs))
	for i := range placeholder {
		placeholder[i] = "?"
	}
	inClause := fmt.Sprintf("(%s)", strings.Join(placeholder, ","))

	rows, err := dao.SearchCart(inClause, args)

	var products []domain.Product

	//使用映射关联商品 ID 和商品信息
	productMap := make(map[int]*domain.Product)
	for i := range products {
		productMap[products[i].ProductID] = &products[i]
	}

	for rows.Next() {
		var product domain.Product
		err = rows.Scan(
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
			return products, err
		}
		product.IsaddedCart = true
		if p, ok := productMap[product.ProductID]; ok {
			*p = product
		}
	}
	return products, nil
}
