package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Booking struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	UserID        uuid.UUID `gorm:"not null;index"`
	PerformanceID uuid.UUID `gorm:"not null;index"`

	TotalPrice int    `gorm:"not null"`
	Status     string `gorm:"default:'pending'"` // pending, confirmed, cancelled
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	User             User              `gorm:"foreignKey:UserID"`
	Performance      Performance       `gorm:"foreignKey:PerformanceID"`
	PerformanceSeats []PerformanceSeat `gorm:"foreignKey:BookingID"`
}

func (Booking) TableName() string {
	return "bookings"
}
