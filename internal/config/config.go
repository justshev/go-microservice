package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServiceName string
	HTTPPort    int
	LogLevel    string
	DBURL	  string
	RedisAddr string 
}

func Load(serviceName string) (Config, error) {
	port := getEnvInt("HTTP_PORT", 8080)
	level := getEnv("LOG_LEVEL", "info")
	dbURL := getEnv("DB_URL", "postgres://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable")
	redisAddr := getEnv("REDDIS_ADDR","localhost:6379")

	cfg := Config{
		ServiceName: serviceName,
		HTTPPort:    port,
		LogLevel:    level,
		DBURL:       dbURL,
		RedisAddr: redisAddr,
	}

	if cfg.HTTPPort <= 0 || cfg.HTTPPort > 65535 {
		return Config{}, fmt.Errorf("invalid HTTP_PORT: %d", cfg.HTTPPort)
	}
	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
