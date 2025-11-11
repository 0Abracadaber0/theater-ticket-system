package repository

import (
	"theater-ticket-system/internal/models/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *Users {
	return &Users{db: db}
}

func (r *Users) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Users) Create(user *model.User) error {
	user.ID = uuid.New()
	return r.db.Create(user).Error
}

func (r *Users) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
