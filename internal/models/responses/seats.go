package response

import (
	"github.com/google/uuid"
)

type Seat struct {
	ID       uuid.UUID `json:"id" binding:"required"`
	HallID   uuid.UUID `json:"hall_id" binding:"required"`
	Row      int       `json:"row" binding:"required"`
	Number   int       `json:"number" binding:"required"`
	Category string    `json:"category" binding:"required"` // parterre, balcony, box
}

type PerformanceSeat struct {
	ID            uuid.UUID `json:"id" binding:"required"`
	PerformanceID uuid.UUID `json:"performance_id" binding:"required"`
	SeatID        uuid.UUID `json:"seat_id" binding:"required"`
	Price         int       `json:"price" binding:"required"`
	Status        string    `json:"status" binding:"required"` // available, reserved, sold
	Seat          *Seat     `json:"seat,omitempty"`
}
