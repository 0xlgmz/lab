package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/middleware"
	"auth-service/internal/models"
	"auth-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database connection
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database models
	db.AutoMigrate(&models.User{}, &models.Session{}, &models.Business{}, &models.LoginAttempt{}, &models.EmailVerification{}, &models.TwoFactorAuth{}, &models.DemoSession{}, &models.PasswordReset{})

	// Initialize repository
	repo := repository.NewAuthRepository(db)

	// Initialize handler
	handler := handlers.NewAuthHandler(repo, cfg)

	// Initialize token manager
	tokenManager := auth.NewTokenManager(cfg.JWTSecret, cfg.JWTRefreshSecret)

	// Initialize Gin router
	router := gin.Default()

	// Add metrics middleware
	router.Use(middleware.MetricsMiddleware("auth-service"))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register routes
	api := router.Group("/auth")
	{
		// Public routes
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
		api.POST("/refresh", handler.RefreshToken)
		api.POST("/verify/send", handler.SendVerificationEmail)
		api.POST("/verify/:token", handler.VerifyEmail)
		api.POST("/password/reset/request", handler.RequestPasswordReset)
		api.POST("/password/reset", handler.ResetPassword)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(tokenManager))
		{
			protected.GET("/user/:id", handler.GetUser)
			protected.PUT("/user/:id", handler.UpdateUser)
			protected.POST("/business", handler.CreateBusiness)
			protected.GET("/business/:id", handler.GetBusiness)
			protected.PUT("/business/:id", handler.UpdateBusiness)
			protected.POST("/logout", handler.Logout)
			protected.POST("/2fa/enable", handler.Enable2FA)
			protected.POST("/2fa/disable", handler.Disable2FA)
			protected.POST("/demo/start", handler.StartDemoSession)
			protected.POST("/demo/end", handler.EndDemoSession)
		}
	}

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
