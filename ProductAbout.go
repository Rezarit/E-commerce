package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowList(client *gin.Context) {
	//获取用户ID
	userID := client.Param("user_id")

	//检查用户ID
	if userID == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"用户 ID 不能为空")
		return
	}
	//创建产品列表
	var products []Product

	//执行查询指令
	cmd := "SELECT * FROM products"
	rows, err := db.Query(cmd)
	if err != nil {
		sendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			"商品列表查询失败")
		return
	}
	defer func() { _ = rows.Close() }() //此处忽略了错误信息

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
			sendErrorResponse(client,
				http.StatusInternalServerError,
				10003,
				"获取商品信息失败")
			return
		}

		product.IsaddedCart, err = isAdded(userID, product.ID)
		if err != nil {
			sendErrorResponse(
				client,
				http.StatusBadRequest,
				10003,
				"数据库查询错误1")
			return
		}

		products = append(products, product)
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

	client.JSON(http.StatusOK, gin.H{
		"status":   10000,
		"info":     "success",
		"products": products,
	})
}

func SearchProduct(client *gin.Context) {
	//获取用户ID
	userID := client.Param("user_id")

	//检查用户ID
	if userID == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"用户 ID 不能为空")
		return
	}

	//获取商品ID
	productID := client.Param("product_id")

	//验证是否为空
	if productID == "" {
		sendErrorResponse(client,
			http.StatusBadRequest,
			10001,
			"商品ID不能为空")
	}

	var product Product

	//执行查询指令
	cmd := "SELECT * FROM products WHERE product_id=?"
	err := db.QueryRow(cmd, productID).Scan(
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
		if err == sql.ErrNoRows {
			sendErrorResponse(
				client,
				http.StatusBadRequest,
				10001,
				"未找到该商品")
		} else {
			sendErrorResponse(
				client,
				http.StatusInternalServerError,
				10003,
				"数据库查询失败")
		}
		return
	}

	product.IsaddedCart, err = isAdded(userID, productID)
	if err != nil {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10003,
			"数据库查询错误")
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status":   10000,
		"info":     "success",
		"products": product,
	})
}

func ProductDetail(client *gin.Context) {
	//获取用户ID
	userID := client.Param("user_id")

	//检查用户ID
	if userID == "" {
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

	var product Product

	//实现查询功能
	cmd := "SELECT * FROM product WHERE product_id=?"

	err := db.QueryRow(cmd, productID).Scan(
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
		if err == sql.ErrNoRows {
			sendErrorResponse(
				client,
				http.StatusBadRequest,
				10001,
				"未找到该商品")
		} else {
			sendErrorResponse(
				client,
				http.StatusInternalServerError,
				10003,
				"数据库查询失败")
		}
		return
	}

	product.IsaddedCart, err = isAdded(userID, productID)
	if err != nil {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10003,
			"数据库查询错误")
		return
	}

	client.JSON(http.StatusOK, gin.H{
		"status":   10000,
		"info":     "success",
		"products": product})
}

func GetType(client *gin.Context) {
	//获取用户ID
	userID := client.Param("user_id")

	//检查用户ID
	if userID == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"用户 ID 不能为空")
		return
	}

	//获取type
	productType := client.Param("type")

	//检验type
	if productType == "" {
		sendErrorResponse(
			client,
			http.StatusBadRequest,
			10001,
			"类型不能为空")
		return
	}

	//执行查询指令
	cmd := "SELECT * FROM product WHERE type=?"
	rows, err := db.Query(cmd, productType)
	if err != nil {
		sendErrorResponse(client,
			http.StatusInternalServerError,
			10003,
			"商品列表查询失败")
		return
	}
	defer func() { _ = rows.Close() }() //此处忽略了错误信息

	var products []Product

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
			sendErrorResponse(client,
				http.StatusInternalServerError,
				10003,
				"获取商品信息失败")
			return
		}

		product.IsaddedCart, err = isAdded(userID, product.ID)
		if err != nil {
			sendErrorResponse(
				client,
				http.StatusBadRequest,
				10003,
				"数据库查询错误1")
			return
		}

		products = append(products, product)
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

	client.JSON(http.StatusOK, gin.H{
		"status": 10000,
		"info":   "success",
		"data":   products})
}
