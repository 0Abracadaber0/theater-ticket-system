package service

import (
	"errors"
	"theater-ticket-system/internal/models/models"
	"theater-ticket-system/internal/repository"

	"github.com/google/uuid"
)

type PlaysRepository interface {
	GetAll() ([]model.Play, error)
	GetByID(id uuid.UUID) (*model.Play, error)
	Create(play *model.Play) error
	Update(play *model.Play) error
	Delete(id uuid.UUID) error
}

type PlaysService struct {
	repo repository.PlaysRepository
}

func NewPlaysService(repo repository.PlaysRepository) *PlaysService {
	return &PlaysService{repo: repo}
}

func (s *PlaysService) GetAllPlays() ([]model.Play, error) {
	return s.repo.GetAll()
}

func (s *PlaysService) GetPlayByID(id string) (*model.Play, error) {
	playID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid play ID format")
	}

	play, err := s.repo.GetByID(playID)
	if err != nil {
		return nil, errors.New("play not found")
	}

	return play, nil
}

func (s *PlaysService) CreatePlay(play *model.Play) error {
	play.ID = uuid.New()
	if play.Title == "" {
		return errors.New("play title is required")
	}
	if play.Author == "" {
		return errors.New("play author is required")
	}
	if play.Duration <= 0 {
		return errors.New("play duration must be positive")
	}

	return s.repo.Create(play)
}

func (s *PlaysService) UpdatePlay(id string, play *model.Play) error {
	playID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid play ID format")
	}

	existing, err := s.repo.GetByID(playID)
	if err != nil {
		return errors.New("play not found")
	}

	play.ID = existing.ID
	play.CreatedAt = existing.CreatedAt

	return s.repo.Update(play)
}

func (s *PlaysService) DeletePlay(id string) error {
	playID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid play ID format")
	}

	_, err = s.repo.GetByID(playID)
	if err != nil {
		return errors.New("play not found")
	}

	return s.repo.Delete(playID)
}
