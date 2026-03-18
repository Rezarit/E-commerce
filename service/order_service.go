package service

import (
	"github.com/Rezarit/go-seckill-system/dao"
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/shopspring/decimal"
	"log"
)

// MakeOrder 下单
func MakeOrder(userID int64, address string) (int64, error) {
	log.Printf("[Service] 开始下单 | 用户ID：%d", userID)

	// 获取用户购物车
	carts, err := getCartItems(userID)
	if err != nil {
		return 0, err
	}

	// 创建订单
	orderID, err := createOrder(userID, address, carts)
	if err != nil {
		return 0, err
	}

	// 清空购物车
	err = clearCart(userID)
	if err != nil {
		return 0, err
	}

	log.Printf("[Service] 下单成功 | 订单ID：%d | 用户ID：%d", orderID, userID)
	return orderID, nil
}

// getCartItems 获取用户购物车商品
func getCartItems(userID int64) ([]domain.Cart, error) {
	carts, err := cartService.GetCartRedis(userID)
	if err != nil {
		log.Printf("[Service] 获取购物车失败 | 用户ID：%d | 错误：%v", userID, err)
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "获取购物车失败",
		}
	}

	if len(carts) == 0 {
		log.Printf("[Service] 购物车为空 | 用户ID：%d", userID)
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeParamInvalid,
			Msg:  "购物车为空",
		}
	}

	return carts, nil
}

// createOrder 创建订单
func createOrder(userID int64, address string, carts []domain.Cart) (int64, error) {
	// 计算订单总金额
	total := calculateTotalAmount(carts)

	// 创建订单
	order := domain.Order{
		UserID:  userID,
		Address: address,
		Total:   total,
		Status:  domain.OrderStatusPending,
	}

	if err := dao.CreateOrder(&order); err != nil {
		log.Printf("[Service] 创建订单失败 | 用户ID：%d | 错误：%v", userID, err)
		return 0, &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "创建订单失败",
		}
	}

	// 创建订单商品并处理库存
	if err := createOrderItemsAndUpdateStock(order.OrderID, carts); err != nil {
		return 0, err
	}

	return order.OrderID, nil
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
func createOrderItemsAndUpdateStock(orderID int64, carts []domain.Cart) error {
	for _, cart := range carts {
		if err := processCartItem(orderID, cart); err != nil {
			return err
		}
	}
	return nil
}

// processCartItem 处理单个购物车商品
func processCartItem(orderID int64, cart domain.Cart) error {
	// 获取商品信息
	product, err := getProductInfo(cart.ProductID)
	if err != nil {
		return err
	}

	// 检查库存
	if err = checkStock(product, cart.Quantity); err != nil {
		return err
	}

	// 创建订单商品
	if err = createOrderItem(orderID, cart, product); err != nil {
		return err
	}

	// 扣减库存
	if err = updateProductStock(product, cart.Quantity); err != nil {
		return err
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
func createOrderItem(orderID int64, cart domain.Cart, product *domain.Product) error {
	orderItem := domain.OrderItem{
		OrderID:     orderID,
		ProductID:   cart.ProductID,
		ProductName: product.ProductName,
		Quantity:    cart.Quantity,
		Price:       product.Price,
	}

	if err := dao.CreateOrderItem(&orderItem); err != nil {
		log.Printf("[Service] 创建订单商品失败 | 订单ID：%d | 商品ID：%d | 错误：%v", orderID, cart.ProductID, err)
		return &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "创建订单商品失败",
		}
	}
	return nil
}

// updateProductStock 更新商品库存
func updateProductStock(product *domain.Product, quantity int) error {
	newStock, err := stockDeductService.DeductStock(product.ProductID, quantity)
	if err != nil {
		log.Printf("扣减库存失败: %v", err)
		return &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "更新商品库存失败",
		}
	}

	log.Printf("库存减扣成功 | 商品ID: %d | 数量: %d | 新库存: %d", product.ProductID, quantity, newStock)
	return nil
}

// clearCart 清空购物车
func clearCart(userID int64) error {
	// 使用Redis服务直接清空整个购物车
	err := cartService.ClearCartRedis(userID)
	if err != nil {
		log.Printf("[Service] 清空购物车失败 | 用户ID：%d | 错误：%v", userID, err)
		return err
	}

	log.Printf("[Service] 清空购物车成功 | 用户ID：%d", userID)
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
