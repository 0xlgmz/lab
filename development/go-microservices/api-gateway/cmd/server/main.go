package main

import (
	"log"
	"os"
	"strings"
	"time"

	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/proxy"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	return logger
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize logger
	logger := initLogger()
	defer logger.Sync()

	// Initialize configuration
	cfg, err := config.New()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware(logger))

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	router.Use(authMiddleware.Authenticate())

	// Initialize rate limiter
	rateLimiter, err := middleware.NewRateLimiter(cfg.RedisURL)
	if err != nil {
		logger.Fatal("Failed to initialize rate limiter", zap.Error(err))
	}

	// Apply rate limiting to all routes except public endpoints
	router.Use(func(c *gin.Context) {
		// List of public endpoints that don't require rate limiting
		publicEndpoints := []string{
			"/api/v1/auth/register",
			"/api/v1/auth/login",
			"/api/v1/auth/refresh",
			"/api/v1/auth/password/reset/request",
			"/api/v1/auth/password/reset",
			"/api/v1/auth/verify",
			"/health",
			"/metrics",
		}

		// Check if the current path is a public endpoint
		currentPath := c.Request.URL.Path
		for _, endpoint := range publicEndpoints {
			if strings.HasPrefix(currentPath, endpoint) {
				c.Next()
				return
			}
		}

		// For all other endpoints, apply rate limiting
		rateLimiter.RateLimit(100, 1*time.Minute)(c)
	})

	// Add metrics middleware
	router.Use(middleware.MetricsMiddleware("api-gateway"))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Initialize proxy
	proxy := proxy.New(
		cfg.AuthServiceURL,
		cfg.BusinessServiceURL,
		cfg.InventoryServiceURL,
		cfg.TransactionServiceURL,
		cfg.FileServiceURL,
		cfg.MenuServiceURL,
		cfg.OrderServiceURL,
		cfg.TableServiceURL,
		logger,
	)

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api/v1")
	{
		// Public endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/register", proxy.HandleRequest)
			auth.POST("/login", proxy.HandleRequest)
			auth.POST("/refresh", proxy.HandleRequest)
			auth.POST("/password/reset/request", proxy.HandleRequest)
			auth.POST("/password/reset", proxy.HandleRequest)
			auth.POST("/verify", proxy.HandleRequest)
		}

		// Protected endpoints
		protected := api.Group("")
		protected.Use(rateLimiter.UserRateLimit(1000, 1*time.Minute)) // 1000 requests per minute per user
		protected.GET("/:service/*path", proxy.HandleRequest)
		protected.POST("/:service/*path", proxy.HandleRequest)
		protected.PUT("/:service/*path", proxy.HandleRequest)
		protected.PATCH("/:service/*path", proxy.HandleRequest)
		protected.DELETE("/:service/*path", proxy.HandleRequest)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting API Gateway", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
