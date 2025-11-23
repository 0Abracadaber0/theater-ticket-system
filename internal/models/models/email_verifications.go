package model

import (
	"time"

	"github.com/google/uuid"
)

// EmailVerification - код подтверждения email
type EmailVerification struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Email     string    `gorm:"not null;index"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time
}

func (*EmailVerification) TableName() string {
	return "email_verifications"
}
