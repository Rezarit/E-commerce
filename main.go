package main

import (
	"github.com/Rezarit/go-seckill-system/pkg/config"
	"github.com/Rezarit/go-seckill-system/route"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	r := route.InitRoute()

	config.InitConfig()
	log.Println("配置初始化成功，JWT密钥已加载")

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
