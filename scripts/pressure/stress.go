package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// 使用你的真实domain结构体
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		UserID       int64  `json:"user_id"`
		Username     string `json:"username"`
		AccessToken  string `json:"access_token"` // ✅ 修正字段名
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}

type OrderCreateRequest struct {
	Address string `json:"address"` // ✅ 只需要address字段
}

type OrderResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		OrderID int64 `json:"order_id"`
	} `json:"data"`
}

// 购物车相关结构体
type CartAddRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type CartResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CartID int64 `json:"cart_id"`
	} `json:"data"`
}

func main() {
	const (
		baseURL           = "http://localhost:8080"
		duration          = 60 * time.Second
		totalStock        = 10000000 // 增加到1000万库存，避免库存不足
		operationsPerUser = 5000     // 每个用户执行5000次操作（极限！）
		workerCount       = 500      // 增加到500个并发worker（极限！）
	)

	fmt.Println("🚀 压测模式：服务日志输出已关闭")
	fmt.Println("💡 提示：如需查看详细日志，请设置 STRESS_TEST=false")

	// 使用压测专用用户（500个）
	testUsers := make([]struct {
		username string
		password string
	}, 500)

	// 生成test01到test500的用户列表
	for i := 0; i < 500; i++ {
		testUsers[i] = struct {
			username string
			password string
		}{
			username: fmt.Sprintf("test%03d", i+1),
			password: "12345678",
		}
	}

	fmt.Printf("🚀 开始基于真实用户流程的压测\n")
	fmt.Printf("⏱️  测试时长: %v\n", duration)
	fmt.Printf("👥 测试用户数: %d\n", len(testUsers))
	fmt.Println("=====================================")

	// 为每个用户获取access_token和用户ID
	type UserInfo struct {
		Token  string
		UserID int64
	}
	userInfos := make(map[string]UserInfo)

	// 并发用户登录（分批控制，避免服务器过载）
	var loginWg sync.WaitGroup
	var loginMu sync.Mutex

	fmt.Println("🚀 开始分批并发用户登录...")

	const loginBatchSize = 50                 // 每批50个用户
	const loginDelay = 100 * time.Millisecond // 批次间延迟

	for batch := 0; batch < len(testUsers); batch += loginBatchSize {
		end := batch + loginBatchSize
		if end > len(testUsers) {
			end = len(testUsers)
		}

		batchUsers := testUsers[batch:end]
		fmt.Printf("🔑 登录批次 %d: 处理用户 %d-%d\n", batch/loginBatchSize+1, batch+1, end)

		// 批次内并发登录
		for _, user := range batchUsers {
			loginWg.Add(1)
			go func(u struct {
				username string
				password string
			}) {
				defer loginWg.Done()

				token, userID, err := loginWithID(baseURL, u.username, u.password)
				if err != nil {
					log.Printf("❌ 用户%s登录失败: %v", u.username, err)
					return
				}

				loginMu.Lock()
				userInfos[u.username] = UserInfo{Token: token, UserID: userID}
				loginMu.Unlock()

				fmt.Printf("✅ 用户%s登录成功! 用户ID: %d\n", u.username, userID)
			}(user)
		}

		loginWg.Wait() // 等待当前批次完成

		// 批次间延迟，避免服务器过载
		if end < len(testUsers) {
			time.Sleep(loginDelay)
		}
	}

	fmt.Printf("🎯 用户登录完成 | 成功登录: %d/%d\n", len(userInfos), len(testUsers))

	if len(userInfos) == 0 {
		log.Fatal("❌ 所有用户登录失败，无法进行压测")
	}

	// 创建测试商品（如果还没有的话）
	productID, err := createTestProduct(baseURL, userInfos["test01"].Token)
	if err != nil {
		log.Printf("⚠️  创建测试商品失败: %v，使用默认商品ID: 1", err)
		productID = 1 // 使用默认商品ID
	}
	fmt.Printf("🛒 使用商品ID: %d 进行压测\n", productID)

	var (
		wg            sync.WaitGroup
		successCount  int64
		failureCount  int64
		totalRequests int64
		startTime     time.Time
	)

	// 2. 创建请求通道和停止通道（增加通道容量）
	requestChan := make(chan int, 10000) // 增加到10000，避免阻塞
	stopChan := make(chan struct{})

	// 3. 启动worker持续发送请求（高并发压测）
	fmt.Printf("👥 启动 %d 个worker持续发送请求...\n", workerCount)

	// 准备用户列表
	usernames := make([]string, 0, len(userInfos))
	for username := range userInfos {
		usernames = append(usernames, username)
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			client := &http.Client{
				Timeout: 30 * time.Second, // 增加到30秒，避免高并发下超时
			}

			// 每个worker固定使用一个用户（避免并发竞争）
			currentUser := usernames[workerID%len(usernames)]
			userInfo := userInfos[currentUser]
			accessToken := userInfo.Token

			// 交替执行：加入购物车 → 下单 → 加入购物车 → 下单
			for i := 0; i < operationsPerUser; i++ {
				select {
				case <-stopChan:
					return
				case <-requestChan:
					requestID := atomic.AddInt64(&totalRequests, 1)

					// 第一步：加入购物车
					err := addToCart(baseURL, accessToken, productID)
					if err != nil {
						if requestID <= 5 {
							fmt.Printf("❌ 用户%s添加购物车失败[%d]: %v\n", currentUser, requestID, err)
						}
						atomic.AddInt64(&failureCount, 1)
						continue
					}

					if requestID <= 5 {
						fmt.Printf("✅ 用户%s添加购物车成功[%d]\n", currentUser, requestID)
					}

					// 短暂等待，确保购物车数据写入完成
					time.Sleep(10 * time.Millisecond)

					// 第二步：下单
					orderReq := OrderCreateRequest{
						Address: fmt.Sprintf("测试地址%d", requestID),
					}

					reqBody, _ := json.Marshal(orderReq)
					req, _ := http.NewRequest("POST", baseURL+"/order/create", bytes.NewBuffer(reqBody))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+accessToken)

					// 发送请求
					resp, err := client.Do(req)
					if err != nil {
						if requestID <= 5 {
							fmt.Printf("❌ 用户%s下单网络错误[%d]: %v\n", currentUser, requestID, err)
						}
						atomic.AddInt64(&failureCount, 1)
						continue
					}

					// 解析响应
					body, _ := io.ReadAll(resp.Body)
					var orderResp OrderResponse
					json.Unmarshal(body, &orderResp)
					resp.Body.Close()

					if resp.StatusCode == http.StatusOK && orderResp.Code == 10000 {
						atomic.AddInt64(&successCount, 1)
						if requestID <= 5 {
							fmt.Printf("✅ 用户%s下单成功[%d]: 订单ID=%d\n", currentUser, requestID, orderResp.Data.OrderID)
						}
					} else {
						atomic.AddInt64(&failureCount, 1)
						if requestID <= 5 {
							fmt.Printf("❌ 用户%s下单失败[%d]: HTTP=%d, Code=%d, Msg=%s\n",
								currentUser, requestID, resp.StatusCode, orderResp.Code, orderResp.Msg)
						}
					}
				}
			}
		}(i)
	}

	// 4. 开始压测计时
	fmt.Printf("⏰ 开始 %v 压测...\n", duration)
	startTime = time.Now()

	// 持续向通道发送请求（极限压测）
	go func() {
		ticker := time.NewTicker(10 * time.Microsecond) // 减少到10微秒，极限压测！
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// 一次性发送大量请求，极限压测
				for i := 0; i < 200; i++ { // 增加到200个请求
					select {
					case requestChan <- 1:
						// 成功发送请求
					default:
						// 通道满，跳过
						break
					}
				}
			case <-stopChan:
				return
			}
		}
	}()

	// 5. 等待压测时间结束
	time.Sleep(duration)
	close(stopChan)
	wg.Wait()

	// 6. 统计结果
	totalDuration := time.Since(startTime)
	totalProcessed := successCount + failureCount
	qps := float64(totalProcessed) / totalDuration.Seconds()

	fmt.Println("=====================================")
	fmt.Printf("🎯 压测结果统计:\n")
	fmt.Printf("⏱️  实际测试时长: %v\n", totalDuration)
	fmt.Printf("📈 总QPS: %.2f\n", qps)
	fmt.Printf("📊 总请求数: %d\n", totalProcessed)
	fmt.Printf("✅ 成功订单: %d\n", successCount)
	fmt.Printf("❌ 失败订单: %d\n", failureCount)

	if totalProcessed > 0 {
		fmt.Printf("📊 成功率: %.2f%%\n", float64(successCount)/float64(totalProcessed)*100)
	}

	// 7. 性能分级
	fmt.Printf("📊 性能分级:\n")
	switch {
	case qps < 100:
		fmt.Printf("   ⚠️  性能一般 (QPS < 100)\n")
	case qps < 500:
		fmt.Printf("   ✅ 性能良好 (100 ≤ QPS < 500)\n")
	case qps < 1000:
		fmt.Printf("   🎉 性能优秀 (500 ≤ QPS < 1000)\n")
	default:
		fmt.Printf("   🚀 性能卓越 (QPS ≥ 1000)\n")
	}
}

// loginWithID 用户登录并返回用户ID
func loginWithID(baseURL, username, password string) (string, int64, error) {
	loginReq := UserLoginRequest{
		Username: username,
		Password: password,
	}

	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		return "", 0, fmt.Errorf("请求体序列化失败: %v", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/user/login", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", 0, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("网络请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("读取响应失败: %v", err)
	}

	var loginResp UserLoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", 0, fmt.Errorf("响应解析失败: %v, 原始响应: %s", err, string(body))
	}

	// ✅ 正确的成功判断：Code == 10000 表示成功！
	if loginResp.Code == 10000 {
		// ✅ 检查access_token是否为空
		if loginResp.Data.AccessToken == "" {
			return "", 0, fmt.Errorf("access_token为空")
		}
		fmt.Printf("✅ 登录成功: Code=%d, Msg=%s\n", loginResp.Code, loginResp.Msg)
		return loginResp.Data.AccessToken, loginResp.Data.UserID, nil
	}

	// ❌ Code != 10000 表示业务逻辑失败
	return "", 0, fmt.Errorf("业务逻辑失败: Code=%d, Msg=%s", loginResp.Code, loginResp.Msg)
}

// 添加商品到购物车
func addToCart(baseURL, accessToken string, productID int64) error {
	// 你的路由是 /cart/add/:product_id，需要请求体包含quantity参数
	cartReq := CartAddRequest{
		Quantity: 1, // 默认添加1个商品
	}

	reqBody, _ := json.Marshal(cartReq)
	url := fmt.Sprintf("%s/cart/add/%d", baseURL, productID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("网络请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP状态码错误: %d, 响应: %s", resp.StatusCode, string(body))
	}

	body, _ := io.ReadAll(resp.Body)
	var cartResp CartResponse
	if err := json.Unmarshal(body, &cartResp); err != nil {
		return fmt.Errorf("响应解析失败: %v, 原始响应: %s", err, string(body))
	}

	if cartResp.Code != 10000 {
		return fmt.Errorf("添加购物车失败: Code=%d, Msg=%s", cartResp.Code, cartResp.Msg)
	}

	return nil
}

// 创建测试商品（如果还没有的话）
func createTestProduct(baseURL, accessToken string) (int64, error) {
	// 这里简化处理，返回默认商品ID
	// 实际项目中可以调用商品创建接口
	return 1, nil // 假设商品ID为1
}
