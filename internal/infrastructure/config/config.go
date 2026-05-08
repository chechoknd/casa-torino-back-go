package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	DatabaseURL         string
	FrontendURL         string
	JWTSecret           string
	JWTExpiresIn        time.Duration
	RefreshTokenExpires time.Duration
	BcryptCost          int
	Port                string
	Env                 string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		FrontendURL:         os.Getenv("FRONTEND_URL"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		JWTExpiresIn:        getDurationEnv("JWT_EXPIRES_IN", 15*time.Minute),
		RefreshTokenExpires: getDurationEnv("REFRESH_TOKEN_EXPIRES", 7*24*time.Hour),
		BcryptCost:          getIntEnv("BCRYPT_COST", bcrypt.DefaultCost),
		Port:                getEnv("PORT", "8080"),
		Env:                 getEnv("ENV", "development"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}
	if cfg.BcryptCost < bcrypt.MinCost || cfg.BcryptCost > bcrypt.MaxCost {
		return Config{}, errors.New("BCRYPT_COST is out of range")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
