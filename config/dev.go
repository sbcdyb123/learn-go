//go:build !k8s

package config

var Config = config{
	DB: dBConfig{
		DSN: "root:root@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: redisConfig{
		Addr: "localhost:6379",
	},
}
