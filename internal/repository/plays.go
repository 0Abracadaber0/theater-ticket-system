package repository

import (
	"theater-ticket-system/internal/models/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Plays struct {
	db *gorm.DB
}

func NewPlays(db *gorm.DB) *Plays {
	return &Plays{db: db}
}

func (r *Plays) GetAll() ([]model.Play, error) {
	var plays []model.Play
	err := r.db.Preload("Performances").
		Order("created_at DESC").Find(&plays).Error
	return plays, err
}

func (r *Plays) GetByID(id uuid.UUID) (*model.Play, error) {
	var play model.Play
	err := r.db.Preload("Performances").First(&play, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &play, nil
}

func (r *Plays) Create(play *model.Play) error {
	if play.ID == uuid.Nil {
		play.ID = uuid.New()
	}
	return r.db.Create(play).Error
}

func (r *Plays) Update(play *model.Play) error {
	return r.db.Save(play).Error
}

func (r *Plays) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Play{}, "id = ?", id).Error
}
