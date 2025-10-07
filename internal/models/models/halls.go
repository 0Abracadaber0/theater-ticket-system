package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Hall struct {
	ID uuid.UUID `gorm:"primaryKey"`

	Name      string         `gorm:"not null"`
	Capacity  int            `gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Seats        []Seat        `gorm:"foreignKey:HallID"`
	Performances []Performance `gorm:"foreignKey:HallID"`
}

func (Hall) TableName() string {
	return "halls"
}
