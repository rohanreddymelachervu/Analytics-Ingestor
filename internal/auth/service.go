package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	jwtSecret string
}

func NewService(db *gorm.DB, jwtSecret string) *Service {
	// Ensure users table exists
	AutoMigrate(db)
	return &Service{db: db, jwtSecret: jwtSecret}
}

func (s *Service) CreateUser(name, email, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &User{Name: name, Email: email, Password: string(hash), Role: role}
	return s.db.Create(user).Error
}

// Authenticate verifies credentials and issues a WRITE-scoped JWT
func (s *Service) Authenticate(email, password string) (string, error) {
	var user User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("invalid credentials")
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	// Determine scopes by role
	var scopes []string
	switch user.Role {
	case "writer":
		scopes = []string{"WRITE"}
	case "reader":
		scopes = []string{"READ"}
	default:
		return "", fmt.Errorf("unrecognized role")
	}
	// Issue token with appropriate scopes
	return GenerateJWT(user.ID, s.jwtSecret, scopes)
}
