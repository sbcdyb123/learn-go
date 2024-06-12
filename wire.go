//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sbcdyb123/learn-go/internal/repository"
	"github.com/sbcdyb123/learn-go/internal/repository/cache"
	"github.com/sbcdyb123/learn-go/internal/repository/dao"
	"github.com/sbcdyb123/learn-go/internal/service"
	"github.com/sbcdyb123/learn-go/internal/web"
	"github.com/sbcdyb123/learn-go/ioc"
)

func initWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		dao.NewUserDao,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		service.NewUserService,
		service.NewCodeService,
		ioc.InitSmsService,
		web.NewUserHandler,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
