package repository

import (
	"theater-ticket-system/internal/models/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Performances struct {
	db *gorm.DB
}

func NewPerformances(db *gorm.DB) *Performances {
	return &Performances{db: db}
}

func (r *Performances) GetAll(playID *uuid.UUID, dateFrom, dateTo *time.Time) ([]model.Performance, error) {
	var performances []model.Performance
	query := r.db.Preload("Play").Order("date ASC")

	if playID != nil {
		query = query.Where("play_id = ?", *playID)
	}
	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}

	err := query.Find(&performances).Error
	return performances, err
}

func (r *Performances) GetByID(id uuid.UUID) (*model.Performance, error) {
	var performance model.Performance
	err := r.db.Preload("Play").First(&performance, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &performance, nil
}

func (r *Performances) GetSeats(performanceID uuid.UUID) ([]model.PerformanceSeat, error) {
	var seats []model.PerformanceSeat
	err := r.db.Preload("Seat").
		Joins("JOIN seats ON seats.id = performance_seats.seat_id").
		Where("performance_seats.performance_id = ?", performanceID).
		Order("seats.row ASC, seats.number ASC").
		Find(&seats).Error
	return seats, err
}
