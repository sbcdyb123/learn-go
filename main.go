package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sbcdyb123/learn-go/internal/web"
)

func main() {
	server := gin.Default()
	u := web.NewUserHandler()
	u.RegisterRoutes(server)
	server.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
