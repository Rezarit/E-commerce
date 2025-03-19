package service

import (
	"database/sql"
	"github.com/Rezarit/E-commerce/dao"
	"github.com/Rezarit/E-commerce/domain"
)

func ListSearch(user domain.User) ([]domain.Product, error) {
	var products []domain.Product
	//执行查询指令
	rows, err := dao.SearchProduct()
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() //此处忽略了错误信息

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
			return nil, err
		}

		product.IsaddedCart, err = dao.IsAdded(user.UserID, product.ProductID)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func TypeSearch(user domain.User, rows *sql.Rows) ([]domain.Product, error) {
	var products []domain.Product

	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
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

		product.IsaddedCart, err = dao.IsAdded(user.UserID, product.ProductID)
		if err != nil {
			return products, err
		}

		products = append(products, product)
	}

	err := rows.Err()
	if err != nil {
		return products, err
	}
	return products, nil
}
