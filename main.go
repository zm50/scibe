package main

import (
	"scibe/global"
	"scibe/handler"
	"scibe/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	global.InitDB()

	r := gin.Default()
	v1 := r.Group("/api/v1")

	user := v1.Group("/user")
	user.POST("/login", handler.Login)
	user.POST("/register", handler.Register)

	file := v1.Group("/file").Use(middleware.Auth)
	file.POST("/upload", handler.UploadFile)
	file.GET("/list/:uid", handler.Files)

	chat := v1.Group("/chat").Use(middleware.Auth)
	chat.GET("/completion/:prompt", handler.ChatCompletion)

	if err := r.Run(":80"); err != nil {
		panic(err)
	}
}
