package dao

import (
	"database/sql"
	"errors"
	"fmt"
)

func IsAdded(userID, productID int) (bool, error) {
	var cartID int
	cmd := "SELECT cart_id FROM shopping_carts WHERE user_id = ? AND product_id = ?"
	err := DB.QueryRow(cmd, userID, productID).Scan(&cartID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func AddintoCart(userID, productID int) error {
	//执行插入语句
	cmd := "INSERT INTO shopping_carts (user_id, product_id, quantity) VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE quantity = quantity + 1"
	_, err := DB.Exec(cmd, userID, productID)
	if err != nil {
		return err
	}
	return nil
}

func CartSearchByID(userID int) (*sql.Rows, error) {
	//执行查询操作
	cmd := "SELECT product_id FROM shopping_carts WHERE user_id = ?"
	rows, err := DB.Query(cmd, userID)
	if err != nil {
		return rows, err
	}
	defer func() { _ = rows.Close() }() //此处忽略了错误信息

	return rows, nil
}

func SearchCart(inClause string, args []any) (*sql.Rows, error) {
	cmd := fmt.Sprintf("SELECT * FROM products WHERE product_id IN %s", inClause)
	rows, err := DB.Query(cmd, args...)
	if err != nil {
		return rows, err
	}
	defer func() { _ = rows.Close() }()
	return rows, nil
}
