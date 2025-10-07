package service

import (
	"errors"
	"github.com/google/uuid"
	"theater-ticket-system/internal/models/models"
	"theater-ticket-system/internal/repository"
)

type PlaysService interface {
	GetAllPlays() ([]model.Play, error)
	GetPlayByID(id string) (*model.Play, error)
	CreatePlay(play *model.Play) error
	UpdatePlay(id string, play *model.Play) error
	DeletePlay(id string) error
}

type playsService struct {
	repo repository.PlaysRepository
}

func NewPlaysService(repo repository.PlaysRepository) PlaysService {
	return &playsService{repo: repo}
}

func (s *playsService) GetAllPlays() ([]model.Play, error) {
	return s.repo.GetAll()
}

func (s *playsService) GetPlayByID(id string) (*model.Play, error) {
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

func (s *playsService) CreatePlay(play *model.Play) error {
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

func (s *playsService) UpdatePlay(id string, play *model.Play) error {
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

func (s *playsService) DeletePlay(id string) error {
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
