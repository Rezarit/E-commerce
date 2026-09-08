package service

import (
	"errors"

	dao2 "github.com/Rezarit/go-seckill-system/internal/dao"
	domain2 "github.com/Rezarit/go-seckill-system/internal/domain"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	myredis "github.com/Rezarit/go-seckill-system/pkg/redis"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// syncOrderSem 降级同步下单的信号量：最多 10 个并发同时写 MySQL，
// 防止 MQ 不可用时大量请求降级把 MySQL 打挂
var syncOrderSem = make(chan struct{}, 10)

// MakeOrder 下单
func MakeOrder(userID int64, address string) (*domain2.OrderAcceptance, error) {
	logger.Sugar.Infof("[Service] 开始下单 | 用户ID：%d", userID)

	items, err := GetCartItems(userID)
	if err != nil {
		return nil, err
	}
	err = CheckCart(items)
	if err != nil {
		return nil, err
	}

	// 先生成受理单号（MsgID）再扣库存：号先造，扣成功才随消息生效（对齐主流「号随扣减事件走」）。
	// 扣库存失败早返回时该号作废即可；正常/降级路径共用同一 msgID，保证幂等一致。
	msgID := uuid.New().String()

	// 入口拦截：在 API 层用 Redis Lua 扣库存，库存不足当场返回，不进 MQ
	// （秒杀核心：让 Redis 挡掉大部分请求，避免 100 个请求全进 MQ 才扣减）
	deducted := make([]domain2.OrderItemMsg, 0, len(items))
	for _, item := range items {
		if _, err := stockDeductService.DeductStock(item.ProductID, item.Quantity); err != nil {
			// 部分扣减补偿：扣 A 成功、扣 B 失败时，把 A 加回，避免少卖
			for _, d := range deducted {
				_ = stockDeductService.AddBackStock(d.ProductID, d.Quantity)
			}
			logger.Sugar.Errorf("[Service] 扣减库存失败 | 商品ID：%d | 错误：%v", item.ProductID, err)
			return nil, &domain2.BusinessError{
				Code: domain2.ErrCodeParamInvalid,
				Msg:  "商品库存不足",
			}
		}
		deducted = append(deducted, domain2.OrderItemMsg{ProductID: item.ProductID, Quantity: item.Quantity})
	}

	// 构造受理凭证（冻结的商品清单 + 受理单号），回给前端/压测拿 msg_id 去轮询结果。
	// 注意：这里不带金额——算金额要查 MySQL，会污染「受理只走 Redis+MQ」的压测主战场；
	// 权威金额以订单落库后的 order.total 为准，由 /order/result 的 success 态带回。
	acc := &domain2.OrderAcceptance{MsgID: msgID, Items: make([]domain2.OrderItemSnapshot, 0, len(deducted))}
	for _, d := range deducted {
		acc.Items = append(acc.Items, domain2.OrderItemSnapshot{ProductID: d.ProductID, Quantity: d.Quantity})
	}

	// 发 MQ（带上已扣减的商品清单，消费者据此建订单，不再读购物车）
	err = createOrder(msgID, userID, address, deducted)
	if err == nil {
		logger.Sugar.Infof("[Service] 下单成功 | 用户ID：%d", userID)
		return acc, nil
	}

	// ===== MQ 不可用，降级为同步下单 =====
	logger.Sugar.Warnf("[Service] MQ 发布失败，降级为同步下单 | 用户ID：%d | 错误：%v", userID, err)

	// 信号量限流：最多 10 个并发同步写 MySQL，防降级把 MySQL 打挂
	select {
	case syncOrderSem <- struct{}{}:
		defer func() { <-syncOrderSem }() // 用完释放车位
	default:
		logger.Sugar.Errorf("[Service] 降级并发超限，拒绝下单 | 用户ID：%d", userID)
		// 释放已扣的库存
		for _, d := range deducted {
			_ = stockDeductService.AddBackStock(d.ProductID, d.Quantity)
		}
		return nil, &domain2.BusinessError{
			Code: domain2.ErrCodeSystemError,
			Msg:  "系统繁忙，请稍后再试",
		}
	}

	// 同步建订单（复用 msgID，走幂等）
	_, syncErr := ExecuteOrderCreation(msgID, userID, address, deducted)
	if syncErr != nil {
		// 同步下单也失败：加回已扣的库存
		for _, d := range deducted {
			_ = stockDeductService.AddBackStock(d.ProductID, d.Quantity)
		}
		logger.Sugar.Errorf("[Service] 降级同步下单失败 | 用户ID：%d | 错误：%v", userID, syncErr)
		return nil, syncErr
	}

	logger.Sugar.Infof("[Service] 降级同步下单成功 | 用户ID：%d", userID)
	return acc, nil
}

// GetCartItems 获取用户购物车商品（Redis 优先，未命中读 MySQL 兜底）
func GetCartItems(userID int64) ([]domain2.CartItem, error) {
	// 1. 空购物车标记命中，直接返回空（防穿透）
	if isNull, err := cartService.IsNullCart(userID); err == nil && isNull {
		return []domain2.CartItem{}, nil
	}

	// 2. 读 Redis 购物车
	items, err := cartService.GetCartRedis(userID)
	if err != nil {
		logger.Sugar.Errorf("[Service] 获取购物车失败 | 用户ID：%d | 错误：%v", userID, err)
		return nil, &domain2.BusinessError{
			Code: domain2.ErrCodeDBError,
			Msg:  "获取购物车失败",
		}
	}
	if len(items) > 0 {
		return items, nil
	}

	// 3. Redis 为空，读 MySQL 兜底
	dbItems, err := dao2.ShowCart(userID)
	if err != nil {
		logger.Sugar.Errorf("[Service] 获取购物车(MySQL)失败 | 用户ID：%d | 错误：%v", userID, err)
		return nil, &domain2.BusinessError{
			Code: domain2.ErrCodeDBError,
			Msg:  "获取购物车失败",
		}
	}

	// 4. MySQL 有数据 → 回填 Redis
	if len(dbItems) > 0 {
		if err := cartService.CacheCartItems(userID, dbItems, myredis.DefaultSessionTTL); err != nil {
			logger.Sugar.Errorf("[Service] 回填购物车缓存失败 | 用户ID：%d | 错误：%v", userID, err)
		}
		return dbItems, nil
	}

	// 5. MySQL 也为空 → 缓存空标记（短 TTL）
	if err := cartService.CacheNullCart(userID, myredis.DefaultNullCacheTTL); err != nil {
		logger.Sugar.Errorf("[Service] 缓存空购物车失败 | 用户ID：%d | 错误：%v", userID, err)
	}
	return []domain2.CartItem{}, nil
}

// ExecuteOrderCreation 创建订单
func ExecuteOrderCreation(msgID string, userID int64, address string, items []domain2.OrderItemMsg) (orderID int64, err error) {
	// 开启事务
	tx := dao2.DB.Begin()
	if tx.Error != nil {
		logger.Sugar.Errorf("[Service] 开启事务失败: %v", tx.Error)
		return 0, &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "无法处理订单"}
	}

	// 确保事务会被处理；recover 接住 panic 后必须把 err 传出去，
	// 否则函数会返回 (0, nil)，调用方误以为下单成功。
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.Sugar.Errorf("[Service] 事务处理中发生 Panic: %v", r)
			err = &domain2.BusinessError{
				Code: domain2.ErrCodeDBError,
				Msg:  "订单处理异常，请稍后重试",
			}
		}
	}()

	// 创建订单（显式落 created = 成交终态；不用 DB 默认 pending，那个词暗示「还有待支付一步」，本系统没有）
	order := domain2.Order{
		MsgID:   msgID,
		UserID:  userID,
		Address: address,
		Total:   calculateTotalAmount(items),
		Status:  domain2.OrderStatusCreated,
	}

	if err := dao2.CreateOrder(tx, &order); err != nil {
		tx.Rollback()
		// 幂等命中：同一条消息已处理过（唯一索引挡住），回滚本次事务并当成功返回，
		// 否则消费者会 Nack 导致 MQ 无限重投这条「已处理」的消息
		if errors.Is(err, dao2.ErrDuplicateMsgID) {
			logger.Sugar.Infof("[Service] 消息重复（幂等命中），订单创建被跳过 | msg_id: %s", msgID)
			return 0, nil
		}
		logger.Sugar.Errorf("[Service] 创建订单失败 | 用户ID：%d | 错误：%v", userID, err)
		return 0, &domain2.BusinessError{
			Code: domain2.ErrCodeDBError,
			Msg:  "创建订单失败",
		}
	}

	// 创建订单商品并处理库存
	if err := createOrderItemsAndUpdateStock(tx, order.OrderID, items); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		logger.Sugar.Errorf("[Service] 提交事务失败: %v", err)
		tx.Rollback()
		return 0, &domain2.BusinessError{Code: domain2.ErrCodeDBError, Msg: "订单最终提交失败"}
	}

	logger.Sugar.Infof("[Service] 订单 %d 创建成功，事务已提交", order.OrderID)
	return order.OrderID, nil
}

// createOrder 创建订单（发 MQ，带上入口已扣减的商品清单）
func createOrder(msgID string, userID int64, address string, items []domain2.OrderItemMsg) error {
	logger.Sugar.Infof("[Service] 接收到下单请求 | 用户ID: %d, 地址: %s", userID, address)

	// 创建订单消息（MsgID 由调用方传入，正常/降级路径共用）
	orderMsg := &domain2.OrderMessage{
		MsgID:   msgID,
		UserID:  userID,
		Address: address,
		Items:   items,
	}

	// 发送消息到MQ
	err := SendMessage(orderMsg, "order")
	if err != nil {
		return err
	}

	logger.Sugar.Infof("[Service] 下单请求已成功发送到MQ | 用户ID: %d", userID)
	return nil // 立刻返回成功
}

// calculateTotalAmount 计算订单总金额
func calculateTotalAmount(items []domain2.OrderItemMsg) decimal.Decimal {
	total := decimal.NewFromInt(0)
	for _, item := range items {
		product, err := dao2.GetProductByID(item.ProductID)
		if err == nil {
			itemTotal := product.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))
			total = total.Add(itemTotal)
		}
	}
	return total
}

// createOrderItemsAndUpdateStock 创建订单商品并更新库存
func createOrderItemsAndUpdateStock(tx *gorm.DB, orderID int64, items []domain2.OrderItemMsg) error {
	for _, item := range items {
		if err := processCartItem(tx, orderID, item); err != nil {
			return err
		}
	}
	return nil
}

// getProductInfo 获取商品信息
func getProductInfo(productID int64) (*domain2.Product, error) {
	product, err := dao2.GetProductByID(productID)
	if err != nil {
		logger.Sugar.Errorf("[Service] 获取商品信息失败 | 商品ID：%d | 错误：%v", productID, err)
		return nil, &domain2.BusinessError{
			Code: domain2.ErrCodeDBError,
			Msg:  "获取商品信息失败",
		}
	}
	return product, nil
}

// createOrderItem 创建订单商品
func createOrderItem(tx *gorm.DB, orderID int64, item domain2.OrderItemMsg, product *domain2.Product) error {
	orderItem := domain2.OrderItem{
		OrderID:     orderID,
		ProductID:   item.ProductID,
		ProductName: product.ProductName,
		Quantity:    item.Quantity,
		Price:       product.Price,
	}

	if err := dao2.CreateOrderItem(tx, &orderItem); err != nil {
		logger.Sugar.Errorf("[Service] 创建订单商品失败 | 订单ID：%d | 商品ID：%d | 错误：%v", orderID, item.ProductID, err)
		return &domain2.BusinessError{
			Code: domain2.ErrCodeDBError,
			Msg:  "创建订单商品失败",
		}
	}
	return nil
}

// updateProductStock 更新商品库存（Redis 已在入口扣减，这里只扣 MySQL）
func updateProductStock(tx *gorm.DB, product *domain2.Product, quantity int) error {
	if err := dao2.DeductStock(tx, product.ProductID, quantity); err != nil {
		logger.Sugar.Errorf("数据库更新库存失败: %v", err)
		return &domain2.BusinessError{
			Code: domain2.ErrCodeDBError,
			Msg:  "更新商品库存失败",
		}
	}

	logger.Sugar.Infof("数据库库存减扣成功 | 商品ID: %d | 数量: %d", product.ProductID, quantity)
	return nil
}

// GetOrderList 获取用户订单列表
func GetOrderList(userID int64) ([]domain2.Order, error) {
	logger.Sugar.Infof("[Service] 获取用户订单列表 | 用户ID：%d", userID)

	orders, err := dao2.GetOrdersByUserID(userID)
	if err != nil {
		logger.Sugar.Errorf("[Service] 获取订单列表失败 | 用户ID：%d | 错误：%v", userID, err)
		return nil, &domain2.BusinessError{
			Code: domain2.ErrCodeDBError,
			Msg:  "获取订单列表失败",
		}
	}

	logger.Sugar.Infof("[Service] 获取订单列表成功 | 用户ID：%d | 订单数量：%d", userID, len(orders))
	return orders, nil
}

// GetOrderDetail 获取订单详情
func GetOrderDetail(orderID, userID int64) (*domain2.Order, []domain2.OrderItem, error) {
	logger.Sugar.Infof("[Service] 获取订单详情 | 订单ID：%d | 用户ID：%d", orderID, userID)

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

	logger.Sugar.Infof("[Service] 获取订单详情成功 | 订单ID：%d | 用户ID：%d", orderID, userID)
	return order, orderItems, nil
}

// getOrderByID 根据订单ID获取订单
func getOrderByID(orderID int64) (*domain2.Order, error) {
	order, err := dao2.GetOrderByID(orderID)
	if err != nil {
		logger.Sugar.Errorf("[Service] 获取订单失败 | 订单ID：%d | 错误：%v", orderID, err)
		return nil, &domain2.BusinessError{
			Code: domain2.ErrCodeDBError,
			Msg:  "获取订单失败",
		}
	}
	return order, nil
}

// checkOrderOwnership 检查订单归属
func checkOrderOwnership(order *domain2.Order, userID int64) error {
	if order.UserID != userID {
		logger.Sugar.Errorf("[Service] 订单归属错误 | 订单ID：%d | 用户ID：%d | 订单所属用户：%d", order.OrderID, userID, order.UserID)
		return &domain2.BusinessError{
			Code: domain2.ErrCodePermissionDenied,
			Msg:  "无权访问此订单",
		}
	}
	return nil
}

// getOrderItemsByOrderID 根据订单ID获取订单商品
func getOrderItemsByOrderID(orderID int64) ([]domain2.OrderItem, error) {
	orderItems, err := dao2.GetOrderItemsByOrderID(orderID)
	if err != nil {
		logger.Sugar.Errorf("[Service] 获取订单商品失败 | 订单ID：%d | 错误：%v", orderID, err)
		return nil, &domain2.BusinessError{
			Code: domain2.ErrCodeDBError,
			Msg:  "获取订单商品失败",
		}
	}
	return orderItems, nil
}

// QueryOrderResult 按 msg_id 查询下单处理结果。
// 返回三态之一：查无键 = processing（消息还在队列/未建单）；success / failed 为终态。
// 归属校验：防止他人拿 msg_id 窥探订单结果（查无态无法校验归属，但 uuid 难猜 + 需登录）。
func QueryOrderResult(msgID string, userID int64) (*domain2.OrderResult, error) {
	result, ok, err := Order.GetOrderResult(msgID)
	if err != nil {
		return nil, err
	}
	if !ok {
		// 键不存在：还没被消费者处理到 → 处理中
		return &domain2.OrderResult{Status: domain2.OrderResultProcessing, UserID: userID}, nil
	}
	if result.UserID != userID {
		return nil, &domain2.BusinessError{
			Code: domain2.ErrCodePermissionDenied,
			Msg:  "无权查看该订单处理结果",
		}
	}
	return &result, nil
}

// ReleaseDeadOrderStock 死信订单的库存释放：按死信单自己的冻结清单精确加回 Redis。
// 只释放「这一单」入口扣过的份额，不做整 key 刷新 MySQL 现值——刷新式会把仍在途订单
// 已预扣的份额误放回 Redis → 入口超发受理（见核心概念笔记 10.9 修订）。
// 死信 = 订单最终没建成 = 该单不该成交，份额归还可卖池。失败靠 Redis TTL 兜底。
func ReleaseDeadOrderStock(items []domain2.OrderItemMsg) error {
	for _, it := range items {
		if err := stockDeductService.AddBackStock(it.ProductID, it.Quantity); err != nil {
			logger.Sugar.Errorf("[Service] 死信库存释放失败 | 商品ID: %d | 数量: %d | 错误: %v", it.ProductID, it.Quantity, err)
			return err
		}
	}
	return nil
}
