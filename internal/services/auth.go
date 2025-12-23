package service

import (
	"errors"
	"theater-ticket-system/internal/models/models"
	"theater-ticket-system/internal/repository"
	"time"

	"github.com/google/uuid"
)

type AuthRepository interface {
	CreateVerification(verification *model.EmailVerification) error
	GetVerification(email, code string) (*model.EmailVerification, error)
	MarkVerificationUsed(id string) error
	DeleteExpiredVerifications() error
}

type Auth struct {
	repo         AuthRepository
	emailService *EmailService
}

func NewAuth(repo *repository.Auth, emailService *EmailService) *Auth {
	return &Auth{
		repo:         repo,
		emailService: emailService,
	}
}

// SendVerificationCode отправляет код подтверждения на email
func (s *Auth) SendVerificationCode(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	// Генерируем код
	code, err := s.emailService.GenerateCode()
	if err != nil {
		return errors.New("failed to generate code")
	}

	// Сохраняем код в БД
	verification := &model.EmailVerification{
		ID:        uuid.New(),
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Used:      false,
	}

	if err := s.repo.CreateVerification(verification); err != nil {
		return errors.New("failed to save verification code")
	}

	// Отправляем код на email
	if err := s.emailService.SendVerificationCode(email, code); err != nil {
		return errors.New("failed to send email")
	}

	return nil
}

// VerifyCode проверяет код подтверждения
func (s *Auth) VerifyCode(email, code string) (bool, error) {
	if email == "" || code == "" {
		return false, errors.New("email and code are required")
	}

	verification, err := s.repo.GetVerification(email, code)
	if err != nil {
		return false, errors.New("invalid or expired code")
	}

	// Помечаем код как использованный
	if err := s.repo.MarkVerificationUsed(verification.ID.String()); err != nil {
		return false, errors.New("failed to mark code as used")
	}

	// Очищаем старые коды
	go s.repo.DeleteExpiredVerifications()

	return true, nil
}
