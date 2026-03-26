package auth

import (
	"time"

	"gorm.io/gorm"
)

type Store struct {
	DB *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) CreateUser(username string) (*User, error) {
	u := &User{
		Username:  username,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.DB.Create(u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Store) CreateToken(userID uint, role, tokenHash string, expiresAt *time.Time) (*Token, error) {
	t := &Token{
		UserID:    userID,
		Role:      role,
		TokenHash: tokenHash,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: expiresAt,
	}

	if err := s.DB.Create(t).Error; err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Store) CreateRegistrationCode(codeHash string, usesRemaining int, expiresAt *time.Time) (*RegistrationCode, error) {
	c := &RegistrationCode{
		CodeHash:      codeHash,
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     expiresAt,
		UsesRemaining: usesRemaining,
	}

	if err := s.DB.Create(c).Error; err != nil {
		return nil, err
	}

	return c, nil
}

// BootstrapExists checks if the bootstrap indicator exists in the database
func (s *Store) BootstrapExists() (bool, error) {
	var count int64

	if err := s.DB.Model(&BootstrapState{}).Where("id = ?", 1).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// MarkBootstrap creates a record in the database to signal
// database initialization has alrady occurred.
func (s *Store) MarkBootstrapped() error {
	return s.DB.Create(&BootstrapState{
		ID:            1,
		InitializedAt: time.Now().UTC(),
	}).Error
}
