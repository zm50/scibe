package main

import (
	"flag"
	"scibe/global"
	"scibe/handler"
	"scibe/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	f := flag.String("f", "", "config file")

	flag.Parse()

	if f == nil {
		panic("config file is required")
	}

	global.InitConfig(*f)
	global.InitDB()
	global.InitLLM()

	r := gin.Default()
	v1 := r.Group("/api/v1")

	user := v1.Group("/user")
	user.POST("/login", handler.Login)
	user.POST("/register", handler.Register)

	file := v1.Group("/file").Use(middleware.Auth)
	file.POST("/upload", handler.UploadFile)
	file.GET("/list", handler.Files)
	file.GET("/text/extract", handler.FileExtractText)

	chat := v1.Group("/chat").Use(middleware.Auth)
	chat.POST("/text/extract", handler.ChatExtractText)
	chat.POST("/summary", handler.ChatPaperSummary)
	chat.POST("/completion", handler.ChatCompletion)

	if err := r.Run(":8000"); err != nil {
		panic(err)
	}
}
