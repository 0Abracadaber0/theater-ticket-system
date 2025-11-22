package response

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID            uuid.UUID         `json:"id" binding:"required"`
	UserID        uuid.UUID         `json:"user_id" binding:"required"`
	PerformanceID uuid.UUID         `json:"performance_id" binding:"required"`
	TotalPrice    int               `json:"total_price" binding:"required"`
	Status        string            `json:"status" binding:"required"` // pending, confirmed, cancelled
	ExpiresAt     time.Time         `json:"expires_at" binding:"required"`
	CreatedAt     time.Time         `json:"created_at" binding:"required"`
	UpdatedAt     time.Time         `json:"updated_at" binding:"required"`
	Performance   *Performance      `json:"performance,omitempty"`
	Seats         []PerformanceSeat `json:"seats,omitempty"`
}
