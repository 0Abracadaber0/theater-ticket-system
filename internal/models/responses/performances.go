package response

import (
	"time"

	"github.com/google/uuid"
)

type Performance struct {
	ID uuid.UUID `json:"id" binding:"required"`

	Date      time.Time `json:"date" binding:"required"`
	Status    string    `json:"status" enums:"scheduled,completed,cancelled"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" binding:"required"`

	Play *Play `json:"play" binding:"omitempty"`
	// Hall Hall `json:"hall" binding:"required"`
}
