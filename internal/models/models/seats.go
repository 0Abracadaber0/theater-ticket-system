package models

import "github.com/google/uuid"

type Seat struct {
	ID     uuid.UUID `gorm:"primaryKey"`
	HallID uuid.UUID `gorm:"not null;index"`

	Row      int    `gorm:"not null"`
	Number   int    `gorm:"not null"`
	Category string `json:"category"` // parterre, balcony, box

	Hall             Hall              `gorm:"foreignKey:HallID" json:"hall,omitempty"`
	PerformanceSeats []PerformanceSeat `gorm:"foreignKey:SeatID" json:"performance_seats,omitempty"`
}

func (Seat) TableName() string {
	return "seats"
}
