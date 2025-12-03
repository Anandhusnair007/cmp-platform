package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerPort   int
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	RedisHost    string
	RedisPort    string
	VaultAddress string
	VaultToken   string
	LogLevel     string
}

func Load() *Config {
	return &Config{
		ServerPort:   getIntEnv("SERVER_PORT", 8080),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "cmp_user"),
		DBPassword:   getEnv("DB_PASSWORD", "cmp_pass"),
		DBName:       getEnv("DB_NAME", "cmp_db"),
		DBSSLMode:    getEnv("DB_SSLMODE", "disable"),
		RedisHost:    getEnv("REDIS_HOST", "localhost"),
		RedisPort:    getEnv("REDIS_PORT", "6379"),
		VaultAddress: getEnv("VAULT_ADDR", "http://localhost:8200"),
		VaultToken:   getEnv("VAULT_TOKEN", "dev-only-token"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Warning: invalid integer value for %s, using default", key)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}
