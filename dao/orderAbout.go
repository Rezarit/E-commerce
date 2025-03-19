package dao

import (
	"database/sql"
	"github.com/Rezarit/E-commerce/domain"
)

func Order(req domain.Order) (sql.Result, error) {
	//执行插入语句
	cmd := "INSERT INTO orders (address, total, user_id) VALUES (?, ?, ?)"
	result, err := DB.Exec(cmd, req.Address, req.Total, req.UserID)
	if err != nil {
		return result, err
	}
	return result, nil
}
