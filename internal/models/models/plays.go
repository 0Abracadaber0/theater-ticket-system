package models

import (
	"github.com/google/uuid"
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

func (Play) TableName() string {
	return "plays"
}
