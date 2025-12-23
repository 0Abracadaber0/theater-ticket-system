package service

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"theater-ticket-system/internal/models/models"
)

type BookingsRepository interface {
	Create(booking *model.Booking) error
	GetByID(id uuid.UUID) (*model.Booking, error)
	GetByUserID(userID uuid.UUID) ([]model.Booking, error)
	Update(booking *model.Booking) error
	UpdatePerformanceSeatStatus(seatID uuid.UUID, status string, bookingID *uuid.UUID) error
	GetPerformanceSeatsByIDs(seatIDs []uuid.UUID, performanceID uuid.UUID) ([]model.PerformanceSeat, error)
}

type UsersRepository interface {
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	GetByID(id uuid.UUID) (*model.User, error)
}

type Bookings struct {
	repo      BookingsRepository
	usersRepo UsersRepository
}

func NewBookings(repo BookingsRepository, usersRepo UsersRepository) *Bookings {
	return &Bookings{
		repo:      repo,
		usersRepo: usersRepo,
	}
}

func (s *Bookings) CreateBooking(email, name string, performanceID uuid.UUID, seatIDs []uuid.UUID) (*model.Booking, error) {
	if len(seatIDs) == 0 {
		return nil, errors.New("at least one seat must be selected")
	}

	// Найти или создать пользователя по email
	user, err := s.usersRepo.FindByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Создаем нового пользователя
			user = &model.User{
				Email:        email,
				Name:         name,
				PasswordHash: "", // Для гостевых бронирований
			}
			if err := s.usersRepo.Create(user); err != nil {
				return nil, errors.New("failed to create user")
			}
		} else {
			return nil, errors.New("failed to find user")
		}
	}

	// Проверяем доступность мест
	seats, err := s.repo.GetPerformanceSeatsByIDs(seatIDs, performanceID)
	if err != nil {
		return nil, err
	}

	if len(seats) != len(seatIDs) {
		return nil, errors.New("some seats are not available")
	}

	// Рассчитываем общую стоимость
	totalPrice := 0
	for _, seat := range seats {
		totalPrice += seat.Price
	}

	// Создаем бронирование
	booking := &model.Booking{
		ID:            uuid.New(),
		UserID:        user.ID,
		PerformanceID: performanceID,
		TotalPrice:    totalPrice,
		Status:        "pending",
	}

	if err := s.repo.Create(booking); err != nil {
		return nil, err
	}

	// Резервируем места
	for _, seat := range seats {
		if err := s.repo.UpdatePerformanceSeatStatus(seat.ID, "reserved", &booking.ID); err != nil {
			return nil, err
		}
	}

	// Получаем полное бронирование с данными
	fullBooking, err := s.repo.GetByID(booking.ID)
	if err != nil {
		return nil, err
	}

	return fullBooking, nil
}

func (s *Bookings) GetBookingByID(id string) (*model.Booking, error) {
	bookingID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid booking ID format")
	}

	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	return booking, nil
}

func (s *Bookings) GetUserBookings(email string) ([]model.Booking, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := s.usersRepo.FindByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []model.Booking{}, nil
		}
		return nil, errors.New("failed to find user")
	}

	return s.repo.GetByUserID(user.ID)
}

func (s *Bookings) CancelBooking(id string) error {
	bookingID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid booking ID format")
	}

	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	if booking.Status == "cancelled" {
		return errors.New("booking already cancelled")
	}

	if booking.Status == "confirmed" {
		return errors.New("cannot cancel confirmed booking")
	}

	booking.Status = "cancelled"
	if err := s.repo.Update(booking); err != nil {
		return err
	}

	// Освобождаем места
	for _, seat := range booking.PerformanceSeats {
		if err := s.repo.UpdatePerformanceSeatStatus(seat.ID, "available", nil); err != nil {
			return err
		}
	}

	return nil
}
