package main

import "github.com/gin-gonic/gin"

func main() {
	//u.RegisterRoutes(server)

	server := initWebServer()

	server.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	server.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
