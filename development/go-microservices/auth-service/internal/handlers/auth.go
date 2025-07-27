package handlers

import (
	"auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	repo         *repository.AuthRepository
	config       *config.Config
	tokenManager *auth.TokenManager
}

func NewAuthHandler(repo *repository.AuthRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		repo:         repo,
		config:       cfg,
		tokenManager: auth.NewTokenManager(cfg.JWTSecret, cfg.JWTSecret), // Using same secret for both access and refresh for now
	}
}

// User handlers
func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create business first
	business := &models.Business{
		Name:     user.FirstName + "'s Business", // Default name
		IsActive: true,
	}

	if err := h.repo.CreateUserWithBusiness(&user, business); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":     user,
		"business": business,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetUserByEmail(credentials.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := user.ComparePassword(credentials.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token pair
	tokenPair, err := h.tokenManager.GenerateTokenPair(user.ID, *user.BusinessID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Create session
	session := &models.Session{
		UserID:       user.ID,
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(h.config.RefreshTokenHours) * time.Hour),
	}

	if err := h.repo.CreateSessionWithUser(session, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"session": session,
	})
}

func (h *AuthHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id
	if err := h.repo.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Business handlers
func (h *AuthHandler) CreateBusiness(c *gin.Context) {
	var business models.Business
	if err := c.ShouldBindJSON(&business); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateBusiness(&business); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, business)
}

func (h *AuthHandler) GetBusiness(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	business, err := h.repo.GetBusinessByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Business not found"})
		return
	}

	c.JSON(http.StatusOK, business)
}

func (h *AuthHandler) UpdateBusiness(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var business models.Business
	if err := c.ShouldBindJSON(&business); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	business.ID = id
	if err := h.repo.UpdateBusiness(&business); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, business)
}

// Session handlers
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	claims, err := h.tokenManager.ValidateToken(request.RefreshToken, auth.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Get session
	session, err := h.repo.GetSessionByRefreshToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if time.Now().After(session.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	// Generate new token pair
	tokenPair, err := h.tokenManager.GenerateTokenPair(claims.UserID, claims.BusinessID, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Update session
	session.Token = tokenPair.AccessToken
	session.RefreshToken = tokenPair.RefreshToken
	session.ExpiresAt = time.Now().Add(time.Duration(h.config.RefreshTokenHours) * time.Hour)

	if err := h.repo.CreateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	// Validate access token
	_, err := h.tokenManager.ValidateToken(token, auth.AccessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	session, err := h.repo.GetSessionByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if err := h.repo.DeleteSession(session.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Email verification handlers
func (h *AuthHandler) SendVerificationEmail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	verification := &models.EmailVerification{
		UserID:    user.ID,
		Token:     "dummy-token", // TODO: Generate secure token
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := h.repo.CreateEmailVerification(verification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: Send verification email

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	verification, err := h.repo.GetEmailVerificationByToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid verification token"})
		return
	}

	if time.Now().After(verification.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token expired"})
		return
	}

	user, err := h.repo.GetUserByID(verification.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.IsEmailVerified = true
	if err := h.repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.DeleteEmailVerification(verification.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// Two-factor authentication handlers
func (h *AuthHandler) Enable2FA(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	tfa := &models.TwoFactorAuth{
		UserID:    id,
		Secret:    "dummy-secret", // TODO: Generate secure secret
		IsEnabled: true,
	}

	if err := h.repo.CreateTwoFactorAuth(tfa); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tfa)
}

func (h *AuthHandler) Disable2FA(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	tfa, err := h.repo.GetTwoFactorAuthByUserID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "2FA not found"})
		return
	}

	tfa.IsEnabled = false
	if err := h.repo.UpdateTwoFactorAuth(tfa); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

// Demo session handlers
func (h *AuthHandler) StartDemoSession(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.BusinessID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has no associated business"})
		return
	}

	// Create demo session
	session := &models.DemoSession{
		UserID:     user.ID,
		BusinessID: *user.BusinessID,
		ExpiresAt:  time.Now().Add(h.config.DemoModeDuration),
	}

	if err := h.repo.CreateDemoSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *AuthHandler) EndDemoSession(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.repo.DeleteDemoSession(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Demo session ended successfully"})
}

// RequestPasswordReset handles password reset requests
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.repo.GetUserByEmail(request.Email)
	if err != nil {
		// Don't reveal if email exists
		c.JSON(http.StatusOK, gin.H{"message": "if the email exists, a reset link will be sent"})
		return
	}

	// Generate reset token
	token := uuid.New().String()
	reset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24), // 24 hours validity
	}

	// Create password reset
	if err := h.repo.CreatePasswordReset(reset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reset token"})
		return
	}

	// TODO: Send reset email with token
	// For now, just return the token in the response
	c.JSON(http.StatusOK, gin.H{
		"message": "reset token generated",
		"token":   token,
	})
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var request struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get password reset
	reset, err := h.repo.GetPasswordResetByToken(request.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired reset token"})
		return
	}

	// Get user
	user, err := h.repo.GetUserByID(reset.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	// Update password
	user.Password = request.Password
	if err := h.repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	// Mark reset as used
	if err := h.repo.MarkPasswordResetAsUsed(reset.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark reset as used"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

// generateToken generates a JWT token for the user
func (h *AuthHandler) generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(h.config.JWTExpirationHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.config.JWTSecret))
}
