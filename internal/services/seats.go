package service

import (
	"errors"
	"theater-ticket-system/internal/models/models"
	"theater-ticket-system/internal/repository"

	"github.com/google/uuid"
)

type SeatsRepository interface {
	GetByHallID(hallID uuid.UUID) ([]model.Seat, error)
}

type Seats struct {
	repo SeatsRepository
}

func NewSeats(repo *repository.Seats) *Seats {
	return &Seats{repo: repo}
}

func (s *Seats) GetSeatsByHallID(hallID string) ([]model.Seat, error) {
	id, err := uuid.Parse(hallID)
	if err != nil {
		return nil, errors.New("invalid hall ID format")
	}

	seats, err := s.repo.GetByHallID(id)
	if err != nil {
		return nil, err
	}

	return seats, nil
}
