// Package app 应用信息
package app

import (
	"fmt"
	"gocms/pkg/config"
)

func IsLocal() bool {
	return config.Get("app.env") == "local"
}

func IsProduction() bool {
	return config.Get("app.env") == "production"
}

func IsTesting() bool {
	fmt.Println("123123")
	return config.Get("app.env") == "testing"
}
