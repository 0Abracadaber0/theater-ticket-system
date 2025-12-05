package service

import (
	"errors"
	"testing"
	"theater-ticket-system/internal/models/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSeatsRepository struct {
	mock.Mock
}

var _ SeatsRepository = (*MockSeatsRepository)(nil)

func (m *MockSeatsRepository) GetByHallID(hallID uuid.UUID) ([]model.Seat, error) {
	args := m.Called(hallID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Seat), args.Error(1)
}

func TestGetSeatsByHallID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockSeatsRepository)
		service := NewSeats(mockRepo)

		hallID := uuid.New()
		expectedSeats := []model.Seat{
			{
				ID:       uuid.New(),
				HallID:   hallID,
				Row:      1,
				Number:   1,
				Category: "parterre",
			},
			{
				ID:       uuid.New(),
				HallID:   hallID,
				Row:      1,
				Number:   2,
				Category: "parterre",
			},
			{
				ID:       uuid.New(),
				HallID:   hallID,
				Row:      5,
				Number:   10,
				Category: "balcony",
			},
		}

		mockRepo.On("GetByHallID", hallID).Return(expectedSeats, nil)

		seats, err := service.GetSeatsByHallID(hallID.String())

		assert.NoError(t, err)
		assert.Equal(t, expectedSeats, seats)
		assert.Len(t, seats, 3)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid hall ID format", func(t *testing.T) {
		mockRepo := new(MockSeatsRepository)
		service := NewSeats(mockRepo)

		seats, err := service.GetSeatsByHallID("invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, seats)
		assert.EqualError(t, err, "invalid hall ID format")
		mockRepo.AssertNotCalled(t, "GetByHallID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockSeatsRepository)
		service := NewSeats(mockRepo)

		hallID := uuid.New()
		mockRepo.On("GetByHallID", hallID).Return(nil, errors.New("database error"))

		seats, err := service.GetSeatsByHallID(hallID.String())

		assert.Error(t, err)
		assert.Nil(t, seats)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty seats list", func(t *testing.T) {
		mockRepo := new(MockSeatsRepository)
		service := NewSeats(mockRepo)

		hallID := uuid.New()
		mockRepo.On("GetByHallID", hallID).Return([]model.Seat{}, nil)

		seats, err := service.GetSeatsByHallID(hallID.String())

		assert.NoError(t, err)
		assert.Empty(t, seats)
		mockRepo.AssertExpectations(t)
	})

	t.Run("seats from different categories", func(t *testing.T) {
		mockRepo := new(MockSeatsRepository)
		service := NewSeats(mockRepo)

		hallID := uuid.New()
		expectedSeats := []model.Seat{
			{ID: uuid.New(), HallID: hallID, Row: 1, Number: 1, Category: "parterre"},
			{ID: uuid.New(), HallID: hallID, Row: 2, Number: 5, Category: "balcony"},
			{ID: uuid.New(), HallID: hallID, Row: 3, Number: 3, Category: "box"},
		}

		mockRepo.On("GetByHallID", hallID).Return(expectedSeats, nil)

		seats, err := service.GetSeatsByHallID(hallID.String())

		assert.NoError(t, err)
		assert.Len(t, seats, 3)

		categories := make(map[string]bool)
		for _, seat := range seats {
			categories[seat.Category] = true
		}
		assert.True(t, categories["parterre"])
		assert.True(t, categories["balcony"])
		assert.True(t, categories["box"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("large number of seats", func(t *testing.T) {
		mockRepo := new(MockSeatsRepository)
		service := NewSeats(mockRepo)

		hallID := uuid.New()
		expectedSeats := make([]model.Seat, 200)
		for i := 0; i < 200; i++ {
			expectedSeats[i] = model.Seat{
				ID:       uuid.New(),
				HallID:   hallID,
				Row:      (i / 20) + 1,
				Number:   (i % 20) + 1,
				Category: "parterre",
			}
		}

		mockRepo.On("GetByHallID", hallID).Return(expectedSeats, nil)

		seats, err := service.GetSeatsByHallID(hallID.String())

		assert.NoError(t, err)
		assert.Len(t, seats, 200)
		mockRepo.AssertExpectations(t)
	})
}
