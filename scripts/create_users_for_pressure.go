package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// UserRegisterRequest 用户注册请求结构
type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserRegisterResponse 用户注册响应结构
type UserRegisterResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		UserID   int64  `json:"user_id"`
		Username string `json:"username"`
	} `json:"data"`
}

func main() {
	const (
		baseURL     = "http://82.157.177.158:8080"
		userCount   = 500        // 创建500个用户
		workerCount = 10         // 10个并发worker
		password    = "12345678" // 统一密码
	)

	fmt.Println("🚀 开始创建压测用户...")
	fmt.Printf("📊 目标用户数: %d\n", userCount)
	fmt.Printf("👥 并发worker数: %d\n", workerCount)
	fmt.Println("=====================================")

	// 创建用户通道
	userChan := make(chan int, userCount)

	// 填充用户编号
	for i := 1; i <= userCount; i++ {
		userChan <- i
	}
	close(userChan)

	// 统计结果
	var (
		successCount int64
		failureCount int64
		mu           sync.Mutex
		wg           sync.WaitGroup
	)

	// 启动worker并发创建用户
	startTime := time.Now()

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			client := &http.Client{
				Timeout: 30 * time.Second,
			}

			for userNum := range userChan {
				username := fmt.Sprintf("test%03d", userNum) // test001, test002, ..., test500

				if err := registerUser(client, baseURL, username, password); err != nil {
					mu.Lock()
					failureCount++
					mu.Unlock()

					if failureCount <= 10 { // 只打印前10个错误
						fmt.Printf("❌ Worker%d 创建用户%s失败: %v\n", workerID, username, err)
					}
				} else {
					mu.Lock()
					successCount++
					mu.Unlock()

					if successCount <= 10 { // 只打印前10个成功
						fmt.Printf("✅ Worker%d 创建用户%s成功\n", workerID, username)
					}
				}

				// 短暂延迟，避免对系统造成过大压力
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	totalDuration := time.Since(startTime)

	// 输出统计结果
	fmt.Println("=====================================")
	fmt.Printf("🎯 用户创建完成统计:\n")
	fmt.Printf("⏱️  总耗时: %v\n", totalDuration)
	fmt.Printf("📊 成功用户: %d\n", successCount)
	fmt.Printf("❌ 失败用户: %d\n", failureCount)
	fmt.Printf("📈 成功率: %.2f%%\n", float64(successCount)/float64(userCount)*100)
	fmt.Printf("🚀 创建速度: %.2f 用户/秒\n", float64(successCount)/totalDuration.Seconds())

	// 生成压测脚本可用的用户列表
	fmt.Println("\n📋 压测用户列表（前20个）:")
	for i := 1; i <= 20 && i <= userCount; i++ {
		fmt.Printf("test%03d/12345678\n", i)
	}

	if userCount > 20 {
		fmt.Printf("... 还有 %d 个用户\n", userCount-20)
	}

	fmt.Println("\n💡 提示: 这些用户可以直接用于压测脚本！")
}

// registerUser 调用注册接口创建用户
func registerUser(client *http.Client, baseURL, username, password string) error {
	registerReq := UserRegisterRequest{
		Username: username,
		Password: password,
	}

	reqBody, err := json.Marshal(registerReq)
	if err != nil {
		return fmt.Errorf("请求体序列化失败: %v", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/user/register", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("网络请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	var registerResp UserRegisterResponse
	if err := json.Unmarshal(body, &registerResp); err != nil {
		return fmt.Errorf("响应解析失败: %v, 原始响应: %s", err, string(body))
	}

	// 检查注册是否成功
	if registerResp.Code == 10000 {
		return nil
	}

	return fmt.Errorf("业务逻辑失败: Code=%d, Msg=%s", registerResp.Code, registerResp.Msg)
}
