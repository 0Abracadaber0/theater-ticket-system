package model

import (
	response "theater-ticket-system/internal/models/responses"

	"github.com/google/uuid"
)

type Seat struct {
	ID     uuid.UUID `gorm:"primaryKey"`
	HallID uuid.UUID `gorm:"not null;index"`

	Row      int    `gorm:"not null"`
	Number   int    `gorm:"not null"`
	Category string // parterre, balcony, box

	Hall             Hall              `gorm:"foreignKey:HallID" json:"hall,omitempty"`
	PerformanceSeats []PerformanceSeat `gorm:"foreignKey:SeatID" json:"performance_seats,omitempty"`
}

func (*Seat) TableName() string {
	return "seats"
}

func (s *Seat) Response() response.Seat {
	return response.Seat{
		ID:       s.ID,
		HallID:   s.HallID,
		Row:      s.Row,
		Number:   s.Number,
		Category: s.Category,
	}
}
