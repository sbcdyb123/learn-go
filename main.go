package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sbcdyb123/learn-go/internal/repository"
	"github.com/sbcdyb123/learn-go/internal/repository/dao"
	"github.com/sbcdyb123/learn-go/internal/service"
	"github.com/sbcdyb123/learn-go/internal/web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	// 初始化数据库
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	// 初始化gin
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "127.0.0.1")
		},
		MaxAge: 12 * time.Hour,
	}))
	u.RegisterRoutes(server)
	server.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
