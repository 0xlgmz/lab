package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// List of public endpoints that don't require authentication
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

		// For all other endpoints, require authentication
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := m.validateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Add claims to context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", claims["user_id"])
			c.Set("businessID", claims["business_id"])
			c.Set("userRole", claims["role"])
		}

		c.Next()
	}
}

func (m *AuthMiddleware) validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.jwtSecret), nil
	})
}
