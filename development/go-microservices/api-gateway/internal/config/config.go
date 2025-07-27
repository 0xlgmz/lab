package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port                  string
	JWTSecret             string
	RedisURL              string
	AuthServiceURL        string
	BusinessServiceURL    string
	InventoryServiceURL   string
	TransactionServiceURL string
	FileServiceURL        string
	MenuServiceURL        string
	OrderServiceURL       string
	TableServiceURL       string
}

func New() (*Config, error) {
	cfg := &Config{
		Port:                  getEnvOrDefault("PORT", "8080"),
		JWTSecret:             getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		RedisURL:              getEnvOrDefault("REDIS_URL", "redis://localhost:6379"),
		AuthServiceURL:        getEnvOrDefault("AUTH_SERVICE_URL", "http://localhost:8081"),
		BusinessServiceURL:    getEnvOrDefault("BUSINESS_SERVICE_URL", "http://localhost:8082"),
		InventoryServiceURL:   getEnvOrDefault("INVENTORY_SERVICE_URL", "http://localhost:8083"),
		TransactionServiceURL: getEnvOrDefault("TRANSACTION_SERVICE_URL", "http://localhost:8084"),
		FileServiceURL:        getEnvOrDefault("FILE_SERVICE_URL", "http://localhost:8085"),
		MenuServiceURL:        getEnvOrDefault("MENU_SERVICE_URL", "http://localhost:8086"),
		OrderServiceURL:       getEnvOrDefault("ORDER_SERVICE_URL", "http://localhost:8087"),
		TableServiceURL:       getEnvOrDefault("TABLE_SERVICE_URL", "http://localhost:8088"),
	}

	// Validate required environment variables
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	required := []string{
		"JWT_SECRET",
		"REDIS_URL",
	}

	for _, env := range required {
		if os.Getenv(env) == "" {
			return fmt.Errorf("required environment variable %s is not set", env)
		}
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
