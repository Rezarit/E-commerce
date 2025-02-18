package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AddCart(client *gin.Context) {
	//绑定数据
	var user User
	if err := client.BindJSON(&user); err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"JSON解析失败")
		return
	}

	//检查用户ID
	if user.ID == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"用户 ID 不能为空")
		return
	}

	//获取商品ID
	productID := client.Param("product_id")

	//检查ID
	if productID == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"商品 ID 不能为空")
		return
	}

	//执行插入语句
	cmd := "INSERT INTO shopping_carts (user_id, product_id, quantity) VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE quantity = quantity + 1"
	_, err := db.Exec(cmd, user.ID, productID)
	if err != nil {
		sendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			fmt.Sprintf("加入购物车失败: %v", err))
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success"})
}

func ShowCart(client *gin.Context) {
	//绑定数据
	var user User
	if err := client.BindJSON(&user); err != nil {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"JSON解析失败")
		return
	}

	//检查用户ID
	if user.ID == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"用户 ID 不能为空")
		return
	}

	//执行查询操作
	cmd := "SELECT product_id FROM shopping_carts WHERE user_id = ?"
	rows, err := db.Query(cmd, user.ID)
	if err != nil {
		sendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			"数据库查询错误1")
		return
	}
	defer func() { _ = rows.Close() }() //此处忽略了错误信息

	var products []Product
	var productIDs []string

	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID)
		if err != nil {
			sendErrorResponse(client,
				http.StatusInternalServerError,
				10003,
				"获取商品信息失败")
			return
		}
		products = append(products, product)
		productIDs = append(productIDs, product.ID)
	}

	err = rows.Err()
	if err != nil {
		sendErrorResponse(
			client,
			http.StatusInternalServerError,
			10003,
			"数据库查询错误2")
		return
	}

	var args []any
	for _, id := range productIDs {
		args = append(args, id)
	}

	if len(productIDs) == 0 {
		client.JSON(http.StatusOK, gin.H{
			"status": 10000,
			"info":   "success",
			"data":   "空的购物车"})
		return
	}

	placeholder := make([]string, len(productIDs))
	for i := range placeholder {
		placeholder[i] = "?"
	}
	inClause := fmt.Sprintf("(%s)", strings.Join(placeholder, ","))

	cmd = fmt.Sprintf("SELECT * FROM products WHERE product_id IN %s", inClause)
	rows, err = db.Query(cmd, args...)
	if err != nil {
		sendErrorResponse(
			client,
			http.StatusInternalServerError,
			10003,
			"商品信息获取失败")
		return
	}

	//使用映射关联商品 ID 和商品信息
	productMap := make(map[string]*Product)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	for rows.Next() {
		var product Product
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Type,
			&product.CommentNum,
			&product.Price,
			&product.Cover,
			&product.PublishTime,
			&product.Link)
		if err != nil {
			sendErrorResponse(
				client,
				http.StatusInternalServerError,
				10003,
				"获取商品详细信息失败")
			return
		}
		product.IsaddedCart = true
		if p, ok := productMap[product.ID]; ok {
			*p = product
		}
	}

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success",
		"data":   products})
}

func isAdded(userID, productID string) (bool, error) {
	var cartID int
	cmd := "SELECT cart_id FROM shopping_carts WHERE user_id = ? AND product_id = ?"
	err := db.QueryRow(cmd, userID, productID).Scan(&cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
