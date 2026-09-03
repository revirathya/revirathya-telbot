package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv              string
	AppPort             string
	TelegramBotToken    string
	TelegramWebhookURL  string
	FinanceCoreGRPCHost string
}

func getEnv(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Loader
func Load() (*Config, error) {
	// Ignore error - .env not exists on 'prod'
	_ = godotenv.Load()

	// Setup config
	cfg := &Config{}
	cfg.AppEnv = getEnv("APP_ENV", "dev")
	cfg.AppPort = getEnv("APP_PORT", "8080")
	cfg.TelegramBotToken = getEnv("TELEGRAM_BOT_TOKEN", "")
	cfg.TelegramWebhookURL = getEnv("TELEGRAM_WEBHOOK_URL", "")
	cfg.FinanceCoreGRPCHost = getEnv("FINANCE_CORE_GRPC_HOST", "")

	// Required config
	if cfg.TelegramBotToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	if cfg.TelegramWebhookURL == "" {
		return nil, fmt.Errorf("TELEGRAM_WEBHOOK_URL is required")
	}

	if cfg.FinanceCoreGRPCHost == "" {
		return nil, fmt.Errorf("FINANCE_CORE_GRPC_HOST is required")
	}

	return cfg, nil
}
