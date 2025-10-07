package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID uuid.UUID `gorm:"primaryKey"`

	Email        string `gorm:"uniqueIndex;not null"`
	Name         string `gorm:"not null"`
	Phone        string `gorm:"uniqueIndex"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	Bookings []Booking `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
