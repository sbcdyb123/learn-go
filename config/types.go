package config

type config struct {
	DB    dBConfig
	Redis redisConfig
}

type dBConfig struct {
	DSN string
}

type redisConfig struct {
	Addr string
}
