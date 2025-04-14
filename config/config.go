package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	RateLimitIP        int
	BlockDurationIP    int
	RateLimitToken     int
	BlockDurationToken int
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not loaded, using system variables.")
	}

	return &Config{
		Port:               getEnv("PORT", "8080"),
		RateLimitIP:        getEnvAsInt("RATE_LIMIT_IP", 5),
		BlockDurationIP:    getEnvAsInt("BLOCK_DURATION_IP", 60),
		RateLimitToken:     getEnvAsInt("RATE_LIMIT_TOKEN_DEFAULT", 10),
		BlockDurationToken: getEnvAsInt("BLOCK_DURATION_TOKEN", 60),
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            getEnvAsInt("REDIS_DB", 0),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if valStr := os.Getenv(key); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
