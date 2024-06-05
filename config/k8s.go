//go:build k8s

package config

var Config = config{
	DB: dBConfig{
		DSN: "root:root@tcp(webook-mysql:13309)/webook?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: redisConfig{
		Addr: "webook-redis:16379",
	},
}
