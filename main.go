package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sbcdyb123/learn-go/config"
	"github.com/sbcdyb123/learn-go/internal/repository"
	"github.com/sbcdyb123/learn-go/internal/repository/cache"
	"github.com/sbcdyb123/learn-go/internal/repository/dao"
	"github.com/sbcdyb123/learn-go/internal/service"
	"github.com/sbcdyb123/learn-go/internal/service/sms/memory"
	"github.com/sbcdyb123/learn-go/internal/web"
	"github.com/sbcdyb123/learn-go/internal/web/middleware"
	"github.com/sbcdyb123/learn-go/pkg/ginx/middlewares/ratelimit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initServer()
	rdb := initRedis()
	u := initUser(db, rdb)
	u.RegisterRoutes(server)
	server.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	server.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}

func initServer() *gin.Engine {
	server := gin.Default()
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "127.0.0.1")
		},
		MaxAge: 12 * time.Hour,
	}))
	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/user/signup").IgnorePaths("/user/login").IgnorePaths("/user/login_sms/code/send").
		IgnorePaths("/user/login_sms").Build())
	return server
}
func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	for err := redisClient.Ping(context.Background()).Err(); err != nil; {
		panic(err)
	}
	return redisClient

}
func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	// 初始化数据库
	ud := dao.NewUserDao(db)
	uc := cache.NewUserCache(rdb)
	cc := cache.NewCodeCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	codeRepo := repository.NewCodeRepository(cc)
	svc := service.NewUserService(repo)
	//client := tencent.NewSmsClient()
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
