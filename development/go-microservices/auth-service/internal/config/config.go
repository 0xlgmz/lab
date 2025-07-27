package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server Configuration
	Port string

	// Database Configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT Configuration
	JWTSecret          string
	JWTRefreshSecret   string
	JWTExpirationHours int
	RefreshTokenHours  int

	// Email Configuration
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string

	// Redis Configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Security Configuration
	PasswordMinLength int
	MaxLoginAttempts  int
	LockoutDuration   time.Duration
	Enable2FA         bool
	SessionTimeout    time.Duration

	// Business Configuration
	DefaultBusinessRole string
	EnableDemoMode      bool
	DemoModeDuration    time.Duration
}

func New() *Config {
	return &Config{
		// Server Configuration
		Port: getEnvOrDefault("PORT", "8081"),

		// Database Configuration
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("DB_PORT", "5432"),
		DBUser:     getEnvOrDefault("DB_USER", "postgres"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:     getEnvOrDefault("DB_NAME", "0xlgmxlab"),

		// JWT Configuration
		JWTSecret:          getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		JWTRefreshSecret:   getEnvOrDefault("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
		JWTExpirationHours: getEnvAsIntOrDefault("JWT_EXPIRATION_HOURS", 24),
		RefreshTokenHours:  getEnvAsIntOrDefault("REFRESH_TOKEN_HOURS", 168), // 7 days

		// Email Configuration
		SMTPHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvOrDefault("SMTP_PORT", "587"),
		SMTPUsername: getEnvOrDefault("SMTP_USERNAME", ""),
		SMTPPassword: getEnvOrDefault("SMTP_PASSWORD", ""),
		FromEmail:    getEnvOrDefault("FROM_EMAIL", "noreply@0xlgmxlab.com"),
		FromName:     getEnvOrDefault("FROM_NAME", "0xlgmzLab"),

		// Redis Configuration
		RedisHost:     getEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:     getEnvOrDefault("REDIS_PORT", "6379"),
		RedisPassword: getEnvOrDefault("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsIntOrDefault("REDIS_DB", 0),

		// Security Configuration
		PasswordMinLength: getEnvAsIntOrDefault("PASSWORD_MIN_LENGTH", 8),
		MaxLoginAttempts:  getEnvAsIntOrDefault("MAX_LOGIN_ATTEMPTS", 5),
		LockoutDuration:   time.Duration(getEnvAsIntOrDefault("LOCKOUT_DURATION_MINUTES", 30)) * time.Minute,
		Enable2FA:         getEnvAsBoolOrDefault("ENABLE_2FA", false),
		SessionTimeout:    time.Duration(getEnvAsIntOrDefault("SESSION_TIMEOUT_MINUTES", 60)) * time.Minute,

		// Business Configuration
		DefaultBusinessRole: getEnvOrDefault("DEFAULT_BUSINESS_ROLE", "Admin"),
		EnableDemoMode:      getEnvAsBoolOrDefault("ENABLE_DEMO_MODE", true),
		DemoModeDuration:    time.Duration(getEnvAsIntOrDefault("DEMO_MODE_DURATION_MINUTES", 30)) * time.Minute,
	}
}

func (c *Config) GetDSN() string {
	return "host=" + c.DBHost + " port=" + c.DBPort + " user=" + c.DBUser + " password=" + c.DBPassword + " dbname=" + c.DBName + " sslmode=disable"
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
