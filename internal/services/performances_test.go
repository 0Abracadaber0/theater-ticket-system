package service

import (
	"errors"
	"testing"
	"theater-ticket-system/internal/models/models"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPerformancesRepository struct {
	mock.Mock
}

func (m *MockPerformancesRepository) GetAll(playID *uuid.UUID, dateFrom, dateTo *time.Time) ([]model.Performance, error) {
	args := m.Called(playID, dateFrom, dateTo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Performance), args.Error(1)
}

func (m *MockPerformancesRepository) GetByID(id uuid.UUID) (*model.Performance, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Performance), args.Error(1)
}

func (m *MockPerformancesRepository) GetSeats(performanceID uuid.UUID) ([]model.PerformanceSeat, error) {
	args := m.Called(performanceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PerformanceSeat), args.Error(1)
}

func TestGetAllPerformances(t *testing.T) {
	t.Run("success without filters", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		expectedPerformances := []model.Performance{
			{
				ID:     uuid.New(),
				PlayID: uuid.New(),
				Date:   time.Now().AddDate(0, 0, 1),
				Status: "scheduled",
			},
		}

		mockRepo.On("GetAll", (*uuid.UUID)(nil), (*time.Time)(nil), (*time.Time)(nil)).
			Return(expectedPerformances, nil)

		performances, err := service.GetAllPerformances(nil, nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, expectedPerformances, performances)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success with play ID filter", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		playID := uuid.New()
		playIDStr := playID.String()

		expectedPerformances := []model.Performance{
			{
				ID:     uuid.New(),
				PlayID: playID,
				Date:   time.Now(),
				Status: "scheduled",
			},
		}

		mockRepo.On("GetAll", &playID, (*time.Time)(nil), (*time.Time)(nil)).
			Return(expectedPerformances, nil)

		performances, err := service.GetAllPerformances(&playIDStr, nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, expectedPerformances, performances)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success with date filters", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		dateFrom := time.Now()
		dateTo := time.Now().AddDate(0, 0, 7)

		expectedPerformances := []model.Performance{
			{
				ID:   uuid.New(),
				Date: time.Now().AddDate(0, 0, 3),
			},
		}

		mockRepo.On("GetAll", (*uuid.UUID)(nil), &dateFrom, &dateTo).
			Return(expectedPerformances, nil)

		performances, err := service.GetAllPerformances(nil, &dateFrom, &dateTo)

		assert.NoError(t, err)
		assert.Equal(t, expectedPerformances, performances)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid play ID format", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		invalidID := "invalid-uuid"

		performances, err := service.GetAllPerformances(&invalidID, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, performances)
		assert.EqualError(t, err, "invalid play ID format")
		mockRepo.AssertNotCalled(t, "GetAll")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		mockRepo.On("GetAll", (*uuid.UUID)(nil), (*time.Time)(nil), (*time.Time)(nil)).
			Return(nil, errors.New("database error"))

		performances, err := service.GetAllPerformances(nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, performances)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetPerformanceByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		performanceID := uuid.New()
		expectedPerformance := &model.Performance{
			ID:     performanceID,
			PlayID: uuid.New(),
			Date:   time.Now(),
			Status: "scheduled",
		}

		mockRepo.On("GetByID", performanceID).Return(expectedPerformance, nil)

		performance, err := service.GetPerformanceByID(performanceID.String())

		assert.NoError(t, err)
		assert.Equal(t, expectedPerformance, performance)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid uuid format", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		performance, err := service.GetPerformanceByID("invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, performance)
		assert.EqualError(t, err, "invalid performance ID format")
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("performance not found", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		performanceID := uuid.New()
		mockRepo.On("GetByID", performanceID).Return(nil, errors.New("not found"))

		performance, err := service.GetPerformanceByID(performanceID.String())

		assert.Error(t, err)
		assert.Nil(t, performance)
		assert.EqualError(t, err, "performance not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetPerformanceSeats(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		performanceID := uuid.New()
		expectedSeats := []model.PerformanceSeat{
			{
				ID:            uuid.New(),
				PerformanceID: performanceID,
				Price:         1500,
				Status:        "available",
			},
			{
				ID:            uuid.New(),
				PerformanceID: performanceID,
				Price:         2000,
				Status:        "reserved",
			},
		}

		mockRepo.On("GetSeats", performanceID).Return(expectedSeats, nil)

		seats, err := service.GetPerformanceSeats(performanceID.String())

		assert.NoError(t, err)
		assert.Equal(t, expectedSeats, seats)
		assert.Len(t, seats, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid uuid format", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		seats, err := service.GetPerformanceSeats("invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, seats)
		assert.EqualError(t, err, "invalid performance ID format")
		mockRepo.AssertNotCalled(t, "GetSeats")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		performanceID := uuid.New()
		mockRepo.On("GetSeats", performanceID).Return(nil, errors.New("database error"))

		seats, err := service.GetPerformanceSeats(performanceID.String())

		assert.Error(t, err)
		assert.Nil(t, seats)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty seats list", func(t *testing.T) {
		mockRepo := new(MockPerformancesRepository)
		service := NewPerformances(mockRepo)

		performanceID := uuid.New()
		mockRepo.On("GetSeats", performanceID).Return([]model.PerformanceSeat{}, nil)

		seats, err := service.GetPerformanceSeats(performanceID.String())

		assert.NoError(t, err)
		assert.Empty(t, seats)
		mockRepo.AssertExpectations(t)
	})
}
