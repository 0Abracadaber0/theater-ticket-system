package request

import (
	model "theater-ticket-system/internal/models/models"
	"time"

	"github.com/google/uuid"
)

type CreateBooking struct {
	Phone         string      `json:"phone" binding:"required"`
	Name          string      `json:"name" binding:"required"`
	PerformanceID uuid.UUID   `json:"performance_id" binding:"required"`
	SeatIDs       []uuid.UUID `json:"seat_ids" binding:"required,min=1"`
}

func (b *CreateBooking) Model(userID uuid.UUID) *model.Booking {
	return &model.Booking{
		UserID:        userID,
		PerformanceID: b.PerformanceID,
		Status:        "pending",
		ExpiresAt:     time.Now().Add(15 * time.Minute),
	}
}
