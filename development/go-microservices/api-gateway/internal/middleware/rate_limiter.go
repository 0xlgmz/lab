package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	redisClient *redis.Client
}

func NewRateLimiter(redisURL string) (*RateLimiter, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RateLimiter{
		redisClient: client,
	}, nil
}

func (rl *RateLimiter) RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())
		count, err := rl.redisClient.Incr(context.Background(), key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process rate limit"})
			c.Abort()
			return
		}

		if count == 1 {
			rl.redisClient.Expire(context.Background(), key, window)
		}

		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) UserRateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		key := fmt.Sprintf("user_rate_limit:%v", userID)
		count, err := rl.redisClient.Incr(context.Background(), key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process rate limit"})
			c.Abort()
			return
		}

		if count == 1 {
			rl.redisClient.Expire(context.Background(), key, window)
		}

		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
