package model

import (
	"time"

	"github.com/google/uuid"
)

// PerformanceSeat - место на конкретном показе
type PerformanceSeat struct {
	ID            uuid.UUID  `gorm:"primaryKey"`
	PerformanceID uuid.UUID  `gorm:"not null;index"`
	SeatID        uuid.UUID  `gorm:"not null;index"`
	BookingID     *uuid.UUID `gorm:"index"`

	Price         int    `gorm:"not null"`
	Status        string `gorm:"default:'available'"` // available, reserved, sold
	ReservedUntil time.Time

	Performance Performance `gorm:"foreignKey:PerformanceID"`
	Seat        Seat        `gorm:"foreignKey:SeatID"`
	Booking     *Booking    `gorm:"foreignKey:BookingID"`
}

func (PerformanceSeat) TableName() string {
	return "performance_seats"
}
