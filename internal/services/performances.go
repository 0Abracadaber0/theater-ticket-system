package service

import (
	"errors"
	"theater-ticket-system/internal/models/models"
	"time"

	"github.com/google/uuid"
)

type PerformancesRepository interface {
	GetAll(playID *uuid.UUID, dateFrom, dateTo *time.Time) ([]model.Performance, error)
	GetByID(id uuid.UUID) (*model.Performance, error)
	GetSeats(performanceID uuid.UUID) ([]model.PerformanceSeat, error)
}

type Performances struct {
	repo PerformancesRepository
}

func NewPerformances(repo PerformancesRepository) *Performances {
	return &Performances{repo: repo}
}

func (s *Performances) GetAllPerformances(playID *string, dateFrom, dateTo *time.Time) ([]model.Performance, error) {
	var playUUID *uuid.UUID
	if playID != nil && *playID != "" {
		parsed, err := uuid.Parse(*playID)
		if err != nil {
			return nil, errors.New("invalid play ID format")
		}
		playUUID = &parsed
	}

	return s.repo.GetAll(playUUID, dateFrom, dateTo)
}

func (s *Performances) GetPerformanceByID(id string) (*model.Performance, error) {
	performanceID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid performance ID format")
	}

	performance, err := s.repo.GetByID(performanceID)
	if err != nil {
		return nil, errors.New("performance not found")
	}

	return performance, nil
}

func (s *Performances) GetPerformanceSeats(id string) ([]model.PerformanceSeat, error) {
	performanceID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid performance ID format")
	}

	seats, err := s.repo.GetSeats(performanceID)
	if err != nil {
		return nil, err
	}

	return seats, nil
}
