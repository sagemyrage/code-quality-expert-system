package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port string
}

type PostgresConfig struct {
	DB       string
	User     string
	Password string
	Host     string
	Port     string
	SSLMode  string
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.DB,
		p.SSLMode,
	)
}

type RedisConfig struct {
	DB       int
	Password string
	Host     string
	Port     string
}

type SessionConfig struct {
	Secret string
}

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Session  SessionConfig
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func mustGetEnv(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}

	return value, nil
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	cfg := &Config{}

	cfg.App.Port = getEnv("APP_PORT", "8080")

	cfg.Postgres.DB, err = mustGetEnv("POSTGRES_DB")
	if err != nil {
		return nil, err
	}

	cfg.Postgres.User, err = mustGetEnv("POSTGRES_USER")
	if err != nil {
		return nil, err
	}

	cfg.Postgres.Password, err = mustGetEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, err
	}

	cfg.Postgres.Host = getEnv("POSTGRES_HOST", "localhost")
	cfg.Postgres.Port = getEnv("POSTGRES_PORT", "5432")
	cfg.Postgres.SSLMode = getEnv("POSTGRES_SSLMODE", "disable")

	redisDBstr := getEnv("REDIS_DB", "0")
	redisDBint, err := strconv.Atoi(redisDBstr)
	if err != nil {
		return nil, fmt.Errorf("REDIS_DB must be an integer: %w", err)
	}

	cfg.Redis.DB = redisDBint
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnv("REDIS_PORT", "6379")

	cfg.Session.Secret, err = mustGetEnv("SESSION_SECRET")
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
