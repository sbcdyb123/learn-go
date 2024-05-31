package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}
func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 登录验证逻辑
		fmt.Println("登录验证逻辑")
	}
}
