package repository

import (
	"auth-service/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// User operations
func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUsersByBusiness(businessID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("business_id = ?", businessID).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetUsersByBranch(branchID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("branch_id = ?", branchID).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *UserRepository) UpdateLastLogin(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login", now).Error
}

// Session operations
func (r *UserRepository) CreateSession(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *UserRepository) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	err := r.db.First(&session, "token = ? AND expires_at > ?", token, time.Now()).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *UserRepository) DeleteSession(id uuid.UUID) error {
	return r.db.Delete(&models.Session{}, "id = ?", id).Error
}

func (r *UserRepository) DeleteExpiredSessions() error {
	return r.db.Where("expires_at <= ?", time.Now()).Delete(&models.Session{}).Error
}

// Password Reset operations
func (r *UserRepository) CreatePasswordReset(reset *models.PasswordReset) error {
	return r.db.Create(reset).Error
}

func (r *UserRepository) GetPasswordResetByToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.First(&reset, "token = ? AND expires_at > ? AND used = ?", token, time.Now(), false).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

func (r *UserRepository) MarkPasswordResetAsUsed(id uuid.UUID) error {
	return r.db.Model(&models.PasswordReset{}).Where("id = ?", id).Update("used", true).Error
}

func (r *UserRepository) DeleteExpiredPasswordResets() error {
	return r.db.Where("expires_at <= ?", time.Now()).Delete(&models.PasswordReset{}).Error
}
