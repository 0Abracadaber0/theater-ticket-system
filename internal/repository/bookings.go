package repository

import (
	"theater-ticket-system/internal/models/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bookings struct {
	db *gorm.DB
}

func NewBookings(db *gorm.DB) *Bookings {
	return &Bookings{db: db}
}

func (r *Bookings) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *Bookings) GetByID(id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("Performance.Play").
		Preload("PerformanceSeats.Seat").
		First(&booking, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *Bookings) GetByUserID(userID uuid.UUID) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Preload("Performance.Play").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookings).Error
	return bookings, err
}

func (r *Bookings) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *Bookings) UpdatePerformanceSeatStatus(seatID uuid.UUID, status string, bookingID *uuid.UUID) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if bookingID != nil {
		updates["booking_id"] = *bookingID
	}
	return r.db.Model(&model.PerformanceSeat{}).
		Where("id = ?", seatID).
		Updates(updates).Error
}

func (r *Bookings) GetPerformanceSeatsByIDs(seatIDs []uuid.UUID, performanceID uuid.UUID) ([]model.PerformanceSeat, error) {
	var seats []model.PerformanceSeat
	err := r.db.Where("id IN ? AND performance_id = ? AND status = ?", seatIDs, performanceID, "available").
		Find(&seats).Error
	return seats, err
}
