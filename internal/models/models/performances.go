package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	response "theater-ticket-system/internal/models/responses"
	"time"
)

// Performance - показ спектакля
type Performance struct {
	ID     uuid.UUID `gorm:"primaryKey"`
	PlayID uuid.UUID `gorm:"not null;index"`
	HallID uuid.UUID `gorm:"not null;index"`

	Date      time.Time `gorm:"not null;index"`
	Status    string    `gorm:"default:'scheduled'"` // scheduled, completed, cancelled
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Play             *Play             `gorm:"foreignKey:PlayID"`
	Hall             Hall              `gorm:"foreignKey:HallID"`
	PerformanceSeats []PerformanceSeat `gorm:"foreignKey:PerformanceID"`
	Bookings         []Booking         `gorm:"foreignKey:PerformanceID"`
}

func (*Performance) TableName() string {
	return "performances"
}

func (p *Performance) Response() response.Performance {
	return response.Performance{
		ID:        p.ID,
		Date:      p.Date,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,

		Play: func() *response.Play {
			if p.Play != nil {
				play := p.Play.Response()
				return &play
			}
			return nil
		}(),
	}
}
