package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type BybitConfig struct {
	PublicKey  string `env:"BYBIT_PUBLIC_KEY"`
	PrivateKey string `env:"BYBIT_PRIVATE_KEY"`
	BaseURL    string `env:"BYBIT_BASE_URL"`
}

func LoadBybitConfig() (*BybitConfig, error) {
	var cfg BybitConfig

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	err = env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error parse to cfg: %w", err)
	}

	return &cfg, nil
}
