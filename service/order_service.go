package service

import (
	"github.com/Rezarit/go-seckill-system/dao"
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"log"
)

// MakeOrder 下单
func MakeOrder(userID int64, address string) error {
	log.Printf("[Service] 开始下单 | 用户ID：%d", userID)

	carts, err := GetCartItems(userID)
	if err != nil {
		return err
	}
	err = CheckCart(carts)
	if err != nil {
		return err
	}

	// 创建订单
	err = createOrder(userID, address)
	if err != nil {
		return err
	}

	log.Printf("[Service] 下单成功 | 用户ID：%d", userID)
	return nil
}

// GetCartItems 获取用户购物车商品
func GetCartItems(userID int64) ([]domain.Cart, error) {
	carts, err := cartService.GetCartRedis(userID)
	if err != nil {
		log.Printf("[Service] 获取购物车失败 | 用户ID：%d | 错误：%v", userID, err)
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "获取购物车失败",
		}
	}
	return carts, nil
}

// ExecuteOrderCreation 创建订单
func ExecuteOrderCreation(userID int64, address string, carts []domain.Cart) (int64, error) {
	// 开启事务
	tx := dao.DB.Begin()
	if tx.Error != nil {
		log.Printf("[Service] 开启事务失败: %v", tx.Error)
		return 0, &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "无法处理订单"}
	}

	// 确保事务会被处理
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[Service] 事务处理中发生 Panic: %v", r)
		}
	}()

	// 创建订单
	order := domain.Order{
		UserID:  userID,
		Address: address,
		Total:   calculateTotalAmount(carts),
	}

	if err := dao.CreateOrder(tx, &order); err != nil {
		log.Printf("[Service] 创建订单失败 | 用户ID：%d | 错误：%v", userID, err)
		tx.Rollback()
		return 0, &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "创建订单失败",
		}
	}

	// 创建订单商品并处理库存
	if err := createOrderItemsAndUpdateStock(tx, order.OrderID, carts); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("[Service] 提交事务失败: %v", err)
		tx.Rollback()
		return 0, &domain.BusinessError{Code: domain.ErrCodeDBError, Msg: "订单最终提交失败"}
	}

	log.Printf("[Service] 订单 %d 创建成功，事务已提交", order.OrderID)
	return order.OrderID, nil
}

// createOrder 创建订单
func createOrder(userID int64, address string) error {
	log.Printf("[Service] 接收到下单请求 | 用户ID: %d, 地址: %s", userID, address)

	// 创建订单消息
	orderMsg := &domain.OrderMessage{
		UserID:  userID,
		Address: address,
	}

	// 发送消息到MQ
	err := SendMessage(orderMsg, "order")
	if err != nil {
		return err
	}

	log.Printf("[Service] 下单请求已成功发送到MQ | 用户ID: %d", userID)
	return nil // 立刻返回成功
}

// calculateTotalAmount 计算订单总金额
func calculateTotalAmount(carts []domain.Cart) decimal.Decimal {
	total := decimal.NewFromInt(0)
	for _, cart := range carts {
		product, err := dao.GetProductByID(cart.ProductID)
		if err == nil {
			itemTotal := product.Price.Mul(decimal.NewFromInt(int64(cart.Quantity)))
			total = total.Add(itemTotal)
		}
	}
	return total
}

// createOrderItemsAndUpdateStock 创建订单商品并更新库存
func createOrderItemsAndUpdateStock(tx *gorm.DB, orderID int64, carts []domain.Cart) error {
	for _, cart := range carts {
		if err := processCartItem(tx, orderID, cart); err != nil {
			return err
		}
	}
	return nil
}

// getProductInfo 获取商品信息
func getProductInfo(productID int64) (*domain.Product, error) {
	product, err := dao.GetProductByID(productID)
	if err != nil {
		log.Printf("[Service] 获取商品信息失败 | 商品ID：%d | 错误：%v", productID, err)
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "获取商品信息失败",
		}
	}
	return product, nil
}

// checkStock 检查库存
func checkStock(product *domain.Product, quantity int) error {
	if product.Stock < quantity {
		log.Printf("[Service] 商品库存不足 | 商品ID：%d | 库存：%d | 需求：%d", product.ProductID, product.Stock, quantity)
		return &domain.BusinessError{
			Code: domain.ErrCodeParamInvalid,
			Msg:  "商品库存不足",
		}
	}
	return nil
}

// createOrderItem 创建订单商品
func createOrderItem(tx *gorm.DB, orderID int64, cart domain.Cart, product *domain.Product) error {
	orderItem := domain.OrderItem{
		OrderID:     orderID,
		ProductID:   cart.ProductID,
		ProductName: product.ProductName,
		Quantity:    cart.Quantity,
		Price:       product.Price,
	}

	if err := dao.CreateOrderItem(tx, &orderItem); err != nil {
		log.Printf("[Service] 创建订单商品失败 | 订单ID：%d | 商品ID：%d | 错误：%v", orderID, cart.ProductID, err)
		return &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "创建订单商品失败",
		}
	}
	return nil
}

// updateProductStock 更新商品库存
func updateProductStock(tx *gorm.DB, product *domain.Product, quantity int) error {
	newStock, err := stockDeductService.DeductStock(product.ProductID, quantity)
	if err != nil {
		log.Printf("redis扣减库存失败: %v", err)
		return &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "更新商品库存失败",
		}
	}
	// 更新数据库中的库存
	product.Stock = newStock
	if err = dao.DeductStock(tx, product.ProductID, quantity); err != nil {
		log.Printf("数据库更新库存失败: %v", err)
		return &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "更新商品库存失败",
		}
	}

	log.Printf("库存减扣成功 | 商品ID: %d | 数量: %d | 新库存: %d", product.ProductID, quantity, newStock)
	return nil
}

// GetOrderList 获取用户订单列表
func GetOrderList(userID int64) ([]domain.Order, error) {
	log.Printf("[Service] 获取用户订单列表 | 用户ID：%d", userID)

	orders, err := dao.GetOrdersByUserID(userID)
	if err != nil {
		log.Printf("[Service] 获取订单列表失败 | 用户ID：%d | 错误：%v", userID, err)
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "获取订单列表失败",
		}
	}

	log.Printf("[Service] 获取订单列表成功 | 用户ID：%d | 订单数量：%d", userID, len(orders))
	return orders, nil
}

// GetOrderDetail 获取订单详情
func GetOrderDetail(orderID, userID int64) (*domain.Order, []domain.OrderItem, error) {
	log.Printf("[Service] 获取订单详情 | 订单ID：%d | 用户ID：%d", orderID, userID)

	// 获取订单
	order, err := getOrderByID(orderID)
	if err != nil {
		return nil, nil, err
	}

	// 检查订单归属
	if err = checkOrderOwnership(order, userID); err != nil {
		return nil, nil, err
	}

	// 获取订单商品
	orderItems, err := getOrderItemsByOrderID(orderID)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("[Service] 获取订单详情成功 | 订单ID：%d | 用户ID：%d", orderID, userID)
	return order, orderItems, nil
}

// getOrderByID 根据订单ID获取订单
func getOrderByID(orderID int64) (*domain.Order, error) {
	order, err := dao.GetOrderByID(orderID)
	if err != nil {
		log.Printf("[Service] 获取订单失败 | 订单ID：%d | 错误：%v", orderID, err)
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "获取订单失败",
		}
	}
	return order, nil
}

// checkOrderOwnership 检查订单归属
func checkOrderOwnership(order *domain.Order, userID int64) error {
	if order.UserID != userID {
		log.Printf("[Service] 订单归属错误 | 订单ID：%d | 用户ID：%d | 订单所属用户：%d", order.OrderID, userID, order.UserID)
		return &domain.BusinessError{
			Code: domain.ErrCodePermissionDenied,
			Msg:  "无权访问此订单",
		}
	}
	return nil
}

// getOrderItemsByOrderID 根据订单ID获取订单商品
func getOrderItemsByOrderID(orderID int64) ([]domain.OrderItem, error) {
	orderItems, err := dao.GetOrderItemsByOrderID(orderID)
	if err != nil {
		log.Printf("[Service] 获取订单商品失败 | 订单ID：%d | 错误：%v", orderID, err)
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "获取订单商品失败",
		}
	}
	return orderItems, nil
}
