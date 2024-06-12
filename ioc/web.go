package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sbcdyb123/learn-go/internal/web"
	"github.com/sbcdyb123/learn-go/internal/web/middleware"
	"github.com/sbcdyb123/learn-go/pkg/ginx/middlewares/ratelimit"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	return server
}
func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		middleware.NewLoginMiddlewareBuilder().IgnorePaths("/user/signup").IgnorePaths("/user/login").IgnorePaths("/user/login_sms/code/send").
			IgnorePaths("/user/login_sms").Build(),
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
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
	})
}
