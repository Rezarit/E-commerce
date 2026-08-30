package main

import (
	"github.com/Rezarit/go-seckill-system/internal/route"
	"github.com/Rezarit/go-seckill-system/pkg/bootstrap"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"github.com/Rezarit/go-seckill-system/pkg/rabbitmq"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//log.SetOutput(io.Discard)

	if err := bootstrap.Init(); err != nil {
		panic(err)
	}
	defer logger.Sync()
	defer rabbitmq.Close()

	r := route.InitRoute()

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
