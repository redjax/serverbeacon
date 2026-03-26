package auth

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"no null"`
	UpdatedAt time.Time
	Tokens    []Token
}

type Token struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	User      User      `gorm:"constraint:OnDelete:CASCADE"`
	Role      string    `gorm:"not null;index"`
	TokenHash string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"not null"`
	ExpiresAt *time.Time
	RevokedAt *time.Time
}

type RegistrationCode struct {
	ID            uint      `gorm:"primaryKey"`
	CodeHash      string    `gorm:"uniqueIndex;not null"`
	CreatedAt     time.Time `gorm:"not null"`
	ExpiresAt     *time.Time
	UsesRemaining int `gorm:"not null;default:1"`
	RevokedAt     *time.Time
}

type BootstrapState struct {
	ID            uint      `gorm:"primaryKey"`
	InitializedAt time.Time `gorm:"not null"`
}
