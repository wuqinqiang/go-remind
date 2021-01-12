package main

import (
	"github.com/gin-gonic/gin"
	"go-remind/config"
	"go-remind/db"
	"go-remind/handlers"
)

func init() {
	config.LoadConfig()
	db.InitDb(config.All.Db)
}

func main() {
	r := gin.Default()
	r.GET("/msg", handlers.Message)
	_ = r.Run()
}
