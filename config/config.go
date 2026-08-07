package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Port        string
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	JWTSecret   string
}

func Load() (*Config, error) {
	c := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) validate() error {
	var missing []string

	if c.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}

	if c.DatabaseURL == "" {
		if c.DBUser == "" {
			missing = append(missing, "DB_USER")
		}

		if c.DBName == "" {
			missing = append(missing, "DB_NAME")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf(
			"config: require environtment variable not set up yet: %s",
			strings.Join(missing, ", "),
		)
	}

	return nil
}

func (c *Config) DSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	q := url.Values{}
	q.Set("sslmode", c.DBSSLMode)

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.DBUser, c.DBPassword),
		Host:     net.JoinHostPort(c.DBHost, c.DBPort),
		Path:     c.DBName,
		RawQuery: q.Encode(),
	}

	return u.String()
}

func (c *Config) Addr() string {
	return ":" + c.Port
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
