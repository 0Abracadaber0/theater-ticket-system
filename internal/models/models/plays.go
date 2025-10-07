package model

import (
	"github.com/google/uuid"
	response "theater-ticket-system/internal/models/responses"
	"time"

	"gorm.io/gorm"
)

// Play - спектакль
type Play struct {
	ID uuid.UUID `gorm:"primaryKey"`

	Title       string `gorm:"not null"`
	Author      string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Duration    int    `gorm:"not null"`
	PosterURL   string
	Genre       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Performances []Performance `gorm:"foreignKey:PlayID"`
}

func (*Play) TableName() string {
	return "plays"
}

func (p *Play) Response() response.Play {
	performances := make([]response.Performance, len(p.Performances))
	for i := range p.Performances {
		performances[i] = p.Performances[i].Response()
	}

	return response.Play{
		ID:          p.ID,
		Title:       p.Title,
		Author:      p.Author,
		Description: p.Description,
		Duration:    p.Duration,
		PosterURL:   p.PosterURL,
		Genre:       p.Genre,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,

		Performances: performances,
	}
}
