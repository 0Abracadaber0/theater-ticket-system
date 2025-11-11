package repository

import (
	"theater-ticket-system/internal/models/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Seats struct {
	db *gorm.DB
}

func NewSeats(db *gorm.DB) *Seats {
	return &Seats{db: db}
}

func (r *Seats) GetByHallID(hallID uuid.UUID) ([]model.Seat, error) {
	var seats []model.Seat
	err := r.db.Where("hall_id = ?", hallID).
		Order("row ASC, number ASC").
		Find(&seats).Error
	return seats, err
}
