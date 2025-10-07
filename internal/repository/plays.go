package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"theater-ticket-system/internal/models/models"
)

type PlaysRepository interface {
	GetAll() ([]model.Play, error)
	GetByID(id uuid.UUID) (*model.Play, error)
	Create(play *model.Play) error
	Update(play *model.Play) error
	Delete(id uuid.UUID) error
}

type playsRepository struct {
	db *gorm.DB
}

func NewPlaysRepository(db *gorm.DB) PlaysRepository {
	return &playsRepository{db: db}
}

func (r *playsRepository) GetAll() ([]model.Play, error) {
	var plays []model.Play
	err := r.db.Preload("Performances").
		Order("created_at DESC").Find(&plays).Error
	return plays, err
}

func (r *playsRepository) GetByID(id uuid.UUID) (*model.Play, error) {
	var play model.Play
	err := r.db.First(&play, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &play, nil
}

func (r *playsRepository) Create(play *model.Play) error {
	play.ID = uuid.New()
	return r.db.Create(play).Error
}

func (r *playsRepository) Update(play *model.Play) error {
	return r.db.Save(play).Error
}

func (r *playsRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Play{}, "id = ?", id).Error
}
