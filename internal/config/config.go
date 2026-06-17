package config

import "os"

const (
	defaultPort        = "8080"
	defaultDatabaseURL = "postgres://postgres:dengke258567@localhost:5432/nextvideo?sslmode=disable"
)

type Config struct {
	Port        string
	DatabaseURL string
}

// Load 读取服务配置；没有环境变量时使用本地开发默认值。
func Load() Config {
	return Config{
		Port:        envOrDefault("PORT", defaultPort),
		DatabaseURL: envOrDefault("DATABASE_URL", defaultDatabaseURL),
	}
}

// envOrDefault 读取指定环境变量；key 为变量名，fallback 为变量为空时使用的值。
func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
