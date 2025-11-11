package model

import (
	response "theater-ticket-system/internal/models/responses"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (*Booking) TableName() string {
	return "bookings"
}

func (b *Booking) Response() response.Booking {
	seats := make([]response.PerformanceSeat, len(b.PerformanceSeats))
	for i := range b.PerformanceSeats {
		seats[i] = b.PerformanceSeats[i].Response()
	}

	return response.Booking{
		ID:            b.ID,
		UserID:        b.UserID,
		PerformanceID: b.PerformanceID,
		TotalPrice:    b.TotalPrice,
		Status:        b.Status,
		ExpiresAt:     b.ExpiresAt,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		Performance: func() *response.Performance {
			if b.Performance.ID != uuid.Nil {
				perf := b.Performance.Response()
				return &perf
			}
			return nil
		}(),
		Seats: seats,
	}
}
