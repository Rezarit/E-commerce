package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// 连接Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// 设置测试数据
	ctx := context.Background()
	rdb.Set(ctx, "seckill:stock:1", 2000, 0)

	// 从文件加载Lua脚本
	scriptPath := filepath.Join("..", "lua", "deduct_stock.lua")
	scriptBytes, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		log.Fatalf("读取Lua脚本文件失败: %v", err)
	}
	script := string(scriptBytes)

	log.Printf("成功加载Lua脚本: %s", scriptPath)

	// 并发测试
	var wg sync.WaitGroup
	successCount := 0
	failCount := 0
	mu := sync.Mutex{}

	start := time.Now()

	// 启动10000个并发goroutine
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			result, err := rdb.Eval(ctx, script, []string{"seckill:stock:1"}, 1).Result()
			if err != nil {
				log.Printf("Goroutine %d: Error: %v", id, err)
				return
			}

			mu.Lock()
			if result.(int64) >= 0 {
				successCount++
			} else {
				failCount++
			}
			mu.Unlock()

			log.Printf("Goroutine %d: Result: %d", id, result)
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	// 检查最终库存
	finalStock, _ := rdb.Get(ctx, "seckill:stock:1").Int()

	fmt.Printf("\n=== 并发测试结果 ===\n")
	fmt.Printf("测试时长: %v\n", duration)
	fmt.Printf("成功扣减: %d 次\n", successCount)
	fmt.Printf("失败次数: %d 次\n", failCount)
	fmt.Printf("最终库存: %d\n", finalStock)
	fmt.Printf("理论最大成功: 10 次\n")
	fmt.Printf("实际成功: %d 次\n", successCount)
	fmt.Printf("是否超卖: %v\n", finalStock >= 0 && successCount <= 10)
}
