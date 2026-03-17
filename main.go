package main

import (
	"github.com/Rezarit/go-seckill-system/pkg/bootstrap"
	"github.com/Rezarit/go-seckill-system/route"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
)

func main() {
	log.SetOutput(io.Discard)

	if err := bootstrap.Init(); err != nil {
		log.Printf("bootstrap err: %v", err)
	}

	r := route.InitRoute()

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
