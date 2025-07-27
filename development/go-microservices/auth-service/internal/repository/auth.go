package repository

import (
	"auth-service/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// User operations
func (r *AuthRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *AuthRepository) DeleteUser(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// Business operations
func (r *AuthRepository) CreateBusiness(business *models.Business) error {
	return r.db.Create(business).Error
}

func (r *AuthRepository) GetBusinessByID(id uuid.UUID) (*models.Business, error) {
	var business models.Business
	err := r.db.First(&business, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

func (r *AuthRepository) UpdateBusiness(business *models.Business) error {
	return r.db.Save(business).Error
}

func (r *AuthRepository) DeleteBusiness(id uuid.UUID) error {
	return r.db.Delete(&models.Business{}, "id = ?", id).Error
}

// Session operations
func (r *AuthRepository) CreateSession(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *AuthRepository) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	err := r.db.First(&session, "token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AuthRepository) GetSessionByRefreshToken(refreshToken string) (*models.Session, error) {
	var session models.Session
	err := r.db.First(&session, "refresh_token = ?", refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AuthRepository) DeleteSession(id uuid.UUID) error {
	return r.db.Delete(&models.Session{}, "id = ?", id).Error
}

// Login attempt operations
func (r *AuthRepository) CreateLoginAttempt(attempt *models.LoginAttempt) error {
	return r.db.Create(attempt).Error
}

func (r *AuthRepository) GetRecentLoginAttempts(userID uuid.UUID, limit int) ([]models.LoginAttempt, error) {
	var attempts []models.LoginAttempt
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&attempts).Error
	if err != nil {
		return nil, err
	}
	return attempts, nil
}

// Email verification operations
func (r *AuthRepository) CreateEmailVerification(verification *models.EmailVerification) error {
	return r.db.Create(verification).Error
}

func (r *AuthRepository) GetEmailVerificationByToken(token string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.First(&verification, "token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (r *AuthRepository) DeleteEmailVerification(id uuid.UUID) error {
	return r.db.Delete(&models.EmailVerification{}, "id = ?", id).Error
}

// Two-factor authentication operations
func (r *AuthRepository) CreateTwoFactorAuth(tfa *models.TwoFactorAuth) error {
	return r.db.Create(tfa).Error
}

func (r *AuthRepository) GetTwoFactorAuthByUserID(userID uuid.UUID) (*models.TwoFactorAuth, error) {
	var tfa models.TwoFactorAuth
	err := r.db.First(&tfa, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &tfa, nil
}

func (r *AuthRepository) UpdateTwoFactorAuth(tfa *models.TwoFactorAuth) error {
	return r.db.Save(tfa).Error
}

// Demo session operations
func (r *AuthRepository) CreateDemoSession(session *models.DemoSession) error {
	return r.db.Create(session).Error
}

func (r *AuthRepository) GetDemoSessionByUserID(userID uuid.UUID) (*models.DemoSession, error) {
	var session models.DemoSession
	err := r.db.First(&session, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AuthRepository) DeleteDemoSession(id uuid.UUID) error {
	return r.db.Delete(&models.DemoSession{}, "id = ?", id).Error
}

// Transaction operations
func (r *AuthRepository) CreateUserWithBusiness(user *models.User, business *models.Business) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create business first
		if err := tx.Create(business).Error; err != nil {
			return err
		}

		// Set business ID for user
		user.BusinessID = &business.ID

		// Create user (password will be hashed by GORM hook)
		if err := tx.Model(&models.User{}).Create(user).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *AuthRepository) CreateSessionWithUser(session *models.Session, user *models.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(session).Error; err != nil {
			return err
		}
		now := time.Now()
		user.LastLoginAt = &now
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		return nil
	})
}

// Cleanup operations
func (r *AuthRepository) CleanupExpiredSessions() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
}

func (r *AuthRepository) CleanupExpiredEmailVerifications() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.EmailVerification{}).Error
}

func (r *AuthRepository) CleanupExpiredDemoSessions() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.DemoSession{}).Error
}

func (r *AuthRepository) CleanupOldLoginAttempts(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.Where("created_at < ?", cutoff).Delete(&models.LoginAttempt{}).Error
}

// Password reset operations
func (r *AuthRepository) CreatePasswordReset(reset *models.PasswordReset) error {
	return r.db.Create(reset).Error
}

func (r *AuthRepository) GetPasswordResetByToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.First(&reset, "token = ? AND used = ?", token, false).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

func (r *AuthRepository) MarkPasswordResetAsUsed(id uuid.UUID) error {
	return r.db.Model(&models.PasswordReset{}).
		Where("id = ?", id).
		Update("used", true).Error
}

func (r *AuthRepository) CleanupExpiredPasswordResets() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.PasswordReset{}).Error
}
