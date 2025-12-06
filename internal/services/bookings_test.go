package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
	"theater-ticket-system/internal/models/models"
)

type MockBookingsRepository struct {
	mock.Mock
}

func (m *MockBookingsRepository) Create(booking *model.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *MockBookingsRepository) GetByID(id uuid.UUID) (*model.Booking, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Booking), args.Error(1)
}

func (m *MockBookingsRepository) GetByUserID(userID uuid.UUID) ([]model.Booking, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Booking), args.Error(1)
}

func (m *MockBookingsRepository) Update(booking *model.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *MockBookingsRepository) UpdatePerformanceSeatStatus(seatID uuid.UUID, status string, bookingID *uuid.UUID) error {
	args := m.Called(seatID, status, bookingID)
	return args.Error(0)
}

func (m *MockBookingsRepository) GetPerformanceSeatsByIDs(seatIDs []uuid.UUID, performanceID uuid.UUID) ([]model.PerformanceSeat, error) {
	args := m.Called(seatIDs, performanceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PerformanceSeat), args.Error(1)
}

type MockUsersRepository struct {
	mock.Mock
}

func (m *MockUsersRepository) FindByPhone(phone string) (*model.User, error) {
	args := m.Called(phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUsersRepository) Create(user *model.User) error {
	args := m.Called(user)
	if args.Error(0) == nil && user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUsersRepository) GetByID(id uuid.UUID) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func TestCreateBooking(t *testing.T) {
	t.Run("success with existing user", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		userID := uuid.New()
		performanceID := uuid.New()
		seatIDs := []uuid.UUID{uuid.New(), uuid.New()}

		existingUser := &model.User{
			ID:    userID,
			Phone: "+1234567890",
			Name:  "John Doe",
		}

		availableSeats := []model.PerformanceSeat{
			{ID: seatIDs[0], Price: 1500, Status: "available"},
			{ID: seatIDs[1], Price: 2000, Status: "available"},
		}

		expectedBooking := &model.Booking{
			ID:            uuid.New(),
			UserID:        userID,
			PerformanceID: performanceID,
			TotalPrice:    3500,
			Status:        "pending",
		}

		mockUsersRepo.On("FindByPhone", "+1234567890").Return(existingUser, nil)
		mockBookingsRepo.On("GetPerformanceSeatsByIDs", seatIDs, performanceID).
			Return(availableSeats, nil)
		mockBookingsRepo.On("Create", mock.AnythingOfType("*model.Booking")).Return(nil)
		mockBookingsRepo.On("UpdatePerformanceSeatStatus", seatIDs[0], "reserved", mock.AnythingOfType("*uuid.UUID")).
			Return(nil)
		mockBookingsRepo.On("UpdatePerformanceSeatStatus", seatIDs[1], "reserved", mock.AnythingOfType("*uuid.UUID")).
			Return(nil)
		mockBookingsRepo.On("GetByID", mock.AnythingOfType("uuid.UUID")).
			Return(expectedBooking, nil)

		booking, err := service.CreateBooking("+1234567890", "John Doe", performanceID, seatIDs)

		assert.NoError(t, err)
		assert.NotNil(t, booking)
		assert.Equal(t, expectedBooking.TotalPrice, booking.TotalPrice)
		mockUsersRepo.AssertExpectations(t)
		mockBookingsRepo.AssertExpectations(t)
	})

	t.Run("success with new user", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		performanceID := uuid.New()
		seatIDs := []uuid.UUID{uuid.New()}

		availableSeats := []model.PerformanceSeat{
			{ID: seatIDs[0], Price: 1500, Status: "available"},
		}

		mockUsersRepo.On("FindByPhone", "+9876543210").Return(nil, gorm.ErrRecordNotFound)
		mockUsersRepo.On("Create", mock.MatchedBy(func(u *model.User) bool {
			return u.Phone == "+9876543210" && u.Name == "Jane Doe"
		})).Return(nil)
		mockBookingsRepo.On("GetPerformanceSeatsByIDs", seatIDs, performanceID).
			Return(availableSeats, nil)
		mockBookingsRepo.On("Create", mock.AnythingOfType("*model.Booking")).Return(nil)
		mockBookingsRepo.On("UpdatePerformanceSeatStatus", seatIDs[0], "reserved", mock.AnythingOfType("*uuid.UUID")).
			Return(nil)
		mockBookingsRepo.On("GetByID", mock.AnythingOfType("uuid.UUID")).
			Return(&model.Booking{ID: uuid.New(), TotalPrice: 1500}, nil)

		booking, err := service.CreateBooking("+9876543210", "Jane Doe", performanceID, seatIDs)

		assert.NoError(t, err)
		assert.NotNil(t, booking)
		mockUsersRepo.AssertExpectations(t)
		mockBookingsRepo.AssertExpectations(t)
	})

	t.Run("no seats selected", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		booking, err := service.CreateBooking("+1234567890", "John", uuid.New(), []uuid.UUID{})

		assert.Error(t, err)
		assert.Nil(t, booking)
		assert.EqualError(t, err, "at least one seat must be selected")
		mockUsersRepo.AssertNotCalled(t, "FindByPhone")
		mockBookingsRepo.AssertNotCalled(t, "Create")
	})

	t.Run("some seats not available", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		userID := uuid.New()
		performanceID := uuid.New()
		seatIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

		existingUser := &model.User{ID: userID, Phone: "+1234567890"}

		availableSeats := []model.PerformanceSeat{
			{ID: seatIDs[0], Price: 1500, Status: "available"},
		}

		mockUsersRepo.On("FindByPhone", "+1234567890").Return(existingUser, nil)
		mockBookingsRepo.On("GetPerformanceSeatsByIDs", seatIDs, performanceID).
			Return(availableSeats, nil)

		booking, err := service.CreateBooking("+1234567890", "John", performanceID, seatIDs)

		assert.Error(t, err)
		assert.Nil(t, booking)
		assert.EqualError(t, err, "some seats are not available")
		mockBookingsRepo.AssertNotCalled(t, "Create")
	})

	t.Run("user creation fails", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		mockUsersRepo.On("FindByPhone", "+1234567890").Return(nil, gorm.ErrRecordNotFound)
		mockUsersRepo.On("Create", mock.AnythingOfType("*model.User")).
			Return(errors.New("database error"))

		booking, err := service.CreateBooking("+1234567890", "John", uuid.New(), []uuid.UUID{uuid.New()})

		assert.Error(t, err)
		assert.Nil(t, booking)
		assert.EqualError(t, err, "failed to create user")
		mockBookingsRepo.AssertNotCalled(t, "Create")
	})
}

func TestGetBookingByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		bookingID := uuid.New()
		expectedBooking := &model.Booking{
			ID:         bookingID,
			TotalPrice: 3500,
			Status:     "confirmed",
		}

		mockBookingsRepo.On("GetByID", bookingID).Return(expectedBooking, nil)

		booking, err := service.GetBookingByID(bookingID.String())

		assert.NoError(t, err)
		assert.Equal(t, expectedBooking, booking)
		mockBookingsRepo.AssertExpectations(t)
	})

	t.Run("invalid uuid format", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		booking, err := service.GetBookingByID("invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, booking)
		assert.EqualError(t, err, "invalid booking ID format")
		mockBookingsRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("booking not found", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		bookingID := uuid.New()
		mockBookingsRepo.On("GetByID", bookingID).Return(nil, errors.New("not found"))

		booking, err := service.GetBookingByID(bookingID.String())

		assert.Error(t, err)
		assert.Nil(t, booking)
		assert.EqualError(t, err, "booking not found")
		mockBookingsRepo.AssertExpectations(t)
	})
}

func TestGetUserBookings(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		userID := uuid.New()
		user := &model.User{ID: userID, Phone: "+1234567890"}

		expectedBookings := []model.Booking{
			{ID: uuid.New(), UserID: userID, TotalPrice: 1500},
			{ID: uuid.New(), UserID: userID, TotalPrice: 2000},
		}

		mockUsersRepo.On("FindByPhone", "+1234567890").Return(user, nil)
		mockBookingsRepo.On("GetByUserID", userID).Return(expectedBookings, nil)

		bookings, err := service.GetUserBookings("+1234567890")

		assert.NoError(t, err)
		assert.Equal(t, expectedBookings, bookings)
		assert.Len(t, bookings, 2)
		mockUsersRepo.AssertExpectations(t)
		mockBookingsRepo.AssertExpectations(t)
	})

	t.Run("empty phone", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		bookings, err := service.GetUserBookings("")

		assert.Error(t, err)
		assert.Nil(t, bookings)
		assert.EqualError(t, err, "phone is required")
		mockUsersRepo.AssertNotCalled(t, "FindByPhone")
	})

	t.Run("user not found", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		mockUsersRepo.On("FindByPhone", "+1234567890").Return(nil, gorm.ErrRecordNotFound)

		bookings, err := service.GetUserBookings("+1234567890")

		assert.NoError(t, err)
		assert.Empty(t, bookings)
		mockUsersRepo.AssertExpectations(t)
		mockBookingsRepo.AssertNotCalled(t, "GetByUserID")
	})
}

func TestCancelBooking(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		bookingID := uuid.New()
		seatID1 := uuid.New()
		seatID2 := uuid.New()

		booking := &model.Booking{
			ID:     bookingID,
			Status: "pending",
			PerformanceSeats: []model.PerformanceSeat{
				{ID: seatID1},
				{ID: seatID2},
			},
		}

		mockBookingsRepo.On("GetByID", bookingID).Return(booking, nil)
		mockBookingsRepo.On("Update", mock.MatchedBy(func(b *model.Booking) bool {
			return b.ID == bookingID && b.Status == "cancelled"
		})).Return(nil)
		mockBookingsRepo.On("UpdatePerformanceSeatStatus", seatID1, "available", (*uuid.UUID)(nil)).
			Return(nil)
		mockBookingsRepo.On("UpdatePerformanceSeatStatus", seatID2, "available", (*uuid.UUID)(nil)).
			Return(nil)

		err := service.CancelBooking(bookingID.String())

		assert.NoError(t, err)
		mockBookingsRepo.AssertExpectations(t)
	})

	t.Run("invalid uuid format", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		err := service.CancelBooking("invalid-uuid")

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid booking ID format")
		mockBookingsRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("booking already cancelled", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		bookingID := uuid.New()
		booking := &model.Booking{
			ID:     bookingID,
			Status: "cancelled",
		}

		mockBookingsRepo.On("GetByID", bookingID).Return(booking, nil)

		err := service.CancelBooking(bookingID.String())

		assert.Error(t, err)
		assert.EqualError(t, err, "booking already cancelled")
		mockBookingsRepo.AssertNotCalled(t, "Update")
	})

	t.Run("cannot cancel confirmed booking", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		bookingID := uuid.New()
		booking := &model.Booking{
			ID:     bookingID,
			Status: "confirmed",
		}

		mockBookingsRepo.On("GetByID", bookingID).Return(booking, nil)

		err := service.CancelBooking(bookingID.String())

		assert.Error(t, err)
		assert.EqualError(t, err, "cannot cancel confirmed booking")
		mockBookingsRepo.AssertNotCalled(t, "Update")
	})

	t.Run("booking not found", func(t *testing.T) {
		mockBookingsRepo := new(MockBookingsRepository)
		mockUsersRepo := new(MockUsersRepository)
		service := NewBookings(mockBookingsRepo, mockUsersRepo)

		bookingID := uuid.New()
		mockBookingsRepo.On("GetByID", bookingID).Return(nil, errors.New("not found"))

		err := service.CancelBooking(bookingID.String())

		assert.Error(t, err)
		assert.EqualError(t, err, "booking not found")
		mockBookingsRepo.AssertNotCalled(t, "Update")
	})
}
