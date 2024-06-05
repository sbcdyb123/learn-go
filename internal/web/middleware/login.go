package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sbcdyb123/learn-go/internal/web"
	"net/http"
	"strings"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}
func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}
func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {

	return func(c *gin.Context) {
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		// 登录验证逻辑
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的Token"})
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		//ParseWithClaims会修改claims的值，所以需要使用指针
		// 验证Token是否有效
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的Token"})
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的Token"})
			return
		}
		now := time.Now()
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			newToken, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
			if err != nil {
				fmt.Println("生成token失败", err)
			}
			c.Header("x-jwt-token", newToken)
		}

		//c.Header("Authorization", "Bearer "+newToken)
		c.Set("claims", claims)
	}
}
