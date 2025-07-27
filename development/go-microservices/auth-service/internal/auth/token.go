package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID     uuid.UUID `json:"user_id"`
	BusinessID uuid.UUID `json:"business_id"`
	Role       string    `json:"role"`
	TokenType  TokenType `json:"token_type"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewTokenManager(accessSecret, refreshSecret string) *TokenManager {
	return &TokenManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     15 * time.Minute,    // 15 minutes
		refreshTTL:    14 * 24 * time.Hour, // 14 days
	}
}

func (tm *TokenManager) GenerateTokenPair(userID, businessID uuid.UUID, role string) (*TokenPair, error) {
	// Generate Access Token
	accessToken, err := tm.generateToken(userID, businessID, role, AccessToken, tm.accessSecret, tm.accessTTL)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshToken, err := tm.generateToken(userID, businessID, role, RefreshToken, tm.refreshSecret, tm.refreshTTL)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (tm *TokenManager) generateToken(userID, businessID uuid.UUID, role string, tokenType TokenType, secret []byte, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID:     userID,
		BusinessID: businessID,
		Role:       role,
		TokenType:  tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (tm *TokenManager) ValidateToken(tokenString string, tokenType TokenType) (*TokenClaims, error) {
	var secret []byte
	if tokenType == AccessToken {
		secret = tm.accessSecret
	} else {
		secret = tm.refreshSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		if claims.TokenType != tokenType {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
