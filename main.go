package main

import (
	"github.com/Rezarit/go-seckill-system/pkg/bootstrap"
	"github.com/Rezarit/go-seckill-system/pkg/rabbitmq"
	"github.com/Rezarit/go-seckill-system/route"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
)

func main() {
	log.SetOutput(io.Discard)

	if err := bootstrap.Init(); err != nil {
		log.Fatalf("bootstrap err: %v", err)
	}
	defer rabbitmq.Close()

	r := route.InitRoute()

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
