package response

import (
	"time"

	"github.com/google/uuid"
)

// Play - спектакль
type Play struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Duration    int       `json:"duration" binding:"required"`
	PosterURL   string    `json:"poster_url" binding:"required"`
	Genre       string    `json:"genre" binding:"required"`
	CreatedAt   time.Time `json:"created_at" binding:"required"`
	UpdatedAt   time.Time `json:"updated_at" binding:"required"`

	Performances []Performance `json:"performances" binding:"omitempty"`
}
