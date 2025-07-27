package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RoleManager      UserRole = "manager"
	RoleStaff        UserRole = "staff"
	RoleCustomer     UserRole = "customer"
	RoleClerk        UserRole = "clerk"
	RoleTeamLead     UserRole = "teamlead"
	RoleFinanceClerk UserRole = "financeclerk"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email           string         `gorm:"uniqueIndex;not null" json:"email"`
	Password        string         `gorm:"not null" json:"password,omitempty"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Role            UserRole       `gorm:"type:varchar(20);not null" json:"role"`
	BusinessID      *uuid.UUID     `gorm:"type:uuid" json:"business_id,omitempty"`
	BranchID        *uuid.UUID     `gorm:"type:uuid" json:"branch_id,omitempty"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	IsEmailVerified bool           `gorm:"default:false" json:"is_email_verified"`
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type Business struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Address   string         `json:"address"`
	Phone     string         `json:"phone"`
	Email     string         `json:"email"`
	Logo      string         `json:"logo"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Session struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token        string    `gorm:"uniqueIndex;not null" json:"token"`
	RefreshToken string    `gorm:"uniqueIndex;not null" json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LoginAttempt struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	IPAddress string    `gorm:"not null" json:"ip_address"`
	Success   bool      `json:"success"`
	CreatedAt time.Time `json:"created_at"`
}

type EmailVerification struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type TwoFactorAuth struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Secret    string    `gorm:"not null" json:"secret"`
	IsEnabled bool      `gorm:"default:false" json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DemoSession struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	BusinessID uuid.UUID `gorm:"type:uuid;not null" json:"business_id"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type PasswordReset struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate is a GORM hook that runs before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	log.Printf("BeforeCreate hook called for user: %s", u.Email)
	log.Printf("Original password length: %d", len(u.Password))
	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	log.Printf("Hashed password length: %d", len(hashedPassword))
	u.Password = string(hashedPassword)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// If the password is being updated, hash it
	if tx.Statement.Changed("Password") {
		log.Printf("BeforeUpdate hook called for user: %s", u.Email)
		log.Printf("Original password length: %d", len(u.Password))
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		log.Printf("Hashed password length: %d", len(hashedPassword))
		u.Password = string(hashedPassword)
	}
	u.UpdatedAt = time.Now()
	return nil
}

// ComparePassword compares the provided password with the hashed password
func (u *User) ComparePassword(password string) error {
	log.Printf("Comparing password for user: %s", u.Email)
	log.Printf("Stored hashed password: %s", u.Password)
	log.Printf("Provided password: %s", password)
	log.Printf("Stored hashed password length: %d", len(u.Password))
	log.Printf("Provided password length: %d", len(password))

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison failed: %v", err)
		return err
	}

	log.Printf("Password comparison successful for user: %s", u.Email)
	return nil
}

// HasPermission checks if the user has the required role
func (u *User) HasPermission(requiredRole UserRole) bool {
	roleHierarchy := map[UserRole]int{
		RoleAdmin:        7,
		RoleManager:      6,
		RoleTeamLead:     5,
		RoleFinanceClerk: 4,
		RoleClerk:        3,
		RoleStaff:        2,
		RoleCustomer:     1,
	}
	return roleHierarchy[u.Role] >= roleHierarchy[requiredRole]
}

func (b *Business) BeforeCreate(tx *gorm.DB) error {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	return nil
}

func (b *Business) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Session) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}

func (la *LoginAttempt) BeforeCreate(tx *gorm.DB) error {
	la.CreatedAt = time.Now()
	return nil
}

func (ev *EmailVerification) BeforeCreate(tx *gorm.DB) error {
	ev.CreatedAt = time.Now()
	return nil
}

func (tfa *TwoFactorAuth) BeforeCreate(tx *gorm.DB) error {
	tfa.CreatedAt = time.Now()
	tfa.UpdatedAt = time.Now()
	return nil
}

func (tfa *TwoFactorAuth) BeforeUpdate(tx *gorm.DB) error {
	tfa.UpdatedAt = time.Now()
	return nil
}

func (ds *DemoSession) BeforeCreate(tx *gorm.DB) error {
	ds.CreatedAt = time.Now()
	return nil
}

func (pr *PasswordReset) BeforeCreate(tx *gorm.DB) error {
	pr.CreatedAt = time.Now()
	pr.UpdatedAt = time.Now()
	return nil
}

func (pr *PasswordReset) BeforeUpdate(tx *gorm.DB) error {
	pr.UpdatedAt = time.Now()
	return nil
}
