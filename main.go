package main

import (
	"github.com/gin-gonic/gin"
	. "go-remind/config"
	"go-remind/db"
	"go-remind/handlers"
	"log"
)

func init() {
	err := LoadConfig()
	if err != nil {
		log.Fatal("初始化错误:", err)
	}

	if err = db.InitDb(ConfAll.Db); err != nil {
		log.Fatal("初始化错误:", err)
	}
}

func main() {
	r := gin.Default()
	r.GET("/msg", handlers.Message)
	_ = r.Run()
}
