package repository

import (
	"theater-ticket-system/internal/models/models"
	"time"

	"gorm.io/gorm"
)

type Auth struct {
	db *gorm.DB
}

func NewAuth(db *gorm.DB) *Auth {
	return &Auth{db: db}
}

func (r *Auth) CreateVerification(verification *model.EmailVerification) error {
	return r.db.Create(verification).Error
}

func (r *Auth) GetVerification(email, code string) (*model.EmailVerification, error) {
	var verification model.EmailVerification
	err := r.db.Where("email = ? AND code = ? AND used = ? AND expires_at > ?",
		email, code, false, time.Now()).
		First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (r *Auth) MarkVerificationUsed(id string) error {
	return r.db.Model(&model.EmailVerification{}).
		Where("id = ?", id).
		Update("used", true).Error
}

func (r *Auth) DeleteExpiredVerifications() error {
	return r.db.Where("expires_at < ? OR used = ?", time.Now(), true).
		Delete(&model.EmailVerification{}).Error
}
