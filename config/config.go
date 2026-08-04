package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		DBHost: 	getEnv("DB_HOST", "localhost"),
		DBPort: 	getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "vrimz"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "postgres"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func getEnv(key, fallback string) string{
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}