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

// MockPlaysRepository - мок репозитория для тестирования
type MockPlaysRepository struct {
	mock.Mock
}

func (m *MockPlaysRepository) GetAll() ([]model.Play, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Play), args.Error(1)
}

func (m *MockPlaysRepository) GetByID(id uuid.UUID) (*model.Play, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Play), args.Error(1)
}

func (m *MockPlaysRepository) Create(play *model.Play) error {
	args := m.Called(play)
	return args.Error(0)
}

func (m *MockPlaysRepository) Update(play *model.Play) error {
	args := m.Called(play)
	return args.Error(0)
}

func (m *MockPlaysRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestGetAllPlays тестирует получение всех спектаклей
func TestGetAllPlays(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		expectedPlays := []model.Play{
			{
				ID:          uuid.New(),
				Title:       "Гамлет",
				Author:      "Шекспир",
				Description: "Трагедия",
				Duration:    180,
				Genre:       "трагедия",
			},
			{
				ID:          uuid.New(),
				Title:       "Вишневый сад",
				Author:      "Чехов",
				Description: "Драма",
				Duration:    150,
				Genre:       "драма",
			},
		}

		mockRepo.On("GetAll").Return(expectedPlays, nil)

		plays, err := service.GetAllPlays()

		assert.NoError(t, err)
		assert.Equal(t, expectedPlays, plays)
		assert.Len(t, plays, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		mockRepo.On("GetAll").Return(nil, errors.New("database error"))

		plays, err := service.GetAllPlays()

		assert.Error(t, err)
		assert.Nil(t, plays)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		mockRepo.On("GetAll").Return([]model.Play{}, nil)

		plays, err := service.GetAllPlays()

		assert.NoError(t, err)
		assert.Empty(t, plays)
		mockRepo.AssertExpectations(t)
	})
}

// TestGetPlayByID тестирует получение спектакля по ID
func TestGetPlayByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		playID := uuid.New()
		expectedPlay := &model.Play{
			ID:          playID,
			Title:       "Гамлет",
			Author:      "Шекспир",
			Description: "Трагедия о принце датском",
			Duration:    180,
			Genre:       "трагедия",
			CreatedAt:   time.Now(),
		}

		mockRepo.On("GetByID", playID).Return(expectedPlay, nil)

		play, err := service.GetPlayByID(playID.String())

		assert.NoError(t, err)
		assert.Equal(t, expectedPlay, play)
		assert.Equal(t, "Гамлет", play.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid uuid format", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		play, err := service.GetPlayByID("invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, play)
		assert.EqualError(t, err, "invalid play ID format")
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("play not found", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		playID := uuid.New()
		mockRepo.On("GetByID", playID).Return(nil, errors.New("not found"))

		play, err := service.GetPlayByID(playID.String())

		assert.Error(t, err)
		assert.Nil(t, play)
		assert.EqualError(t, err, "play not found")
		mockRepo.AssertExpectations(t)
	})
}

// TestCreatePlay тестирует создание нового спектакля
func TestCreatePlay(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		newPlay := &model.Play{
			Title:       "Ревизор",
			Author:      "Гоголь",
			Description: "Комедия",
			Duration:    120,
			Genre:       "комедия",
			PosterURL:   "https://example.com/poster.jpg",
		}

		mockRepo.On("Create", mock.MatchedBy(func(p *model.Play) bool {
			return p.Title == "Ревизор" && p.Author == "Гоголь" && p.Duration == 120
		})).Return(nil)

		err := service.CreatePlay(newPlay)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, newPlay.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("missing title", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		newPlay := &model.Play{
			Title:    "",
			Author:   "Гоголь",
			Duration: 120,
		}

		err := service.CreatePlay(newPlay)

		assert.Error(t, err)
		assert.EqualError(t, err, "play title is required")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("missing author", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		newPlay := &model.Play{
			Title:    "Ревизор",
			Author:   "",
			Duration: 120,
		}

		err := service.CreatePlay(newPlay)

		assert.Error(t, err)
		assert.EqualError(t, err, "play author is required")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("invalid duration zero", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		newPlay := &model.Play{
			Title:    "Ревизор",
			Author:   "Гоголь",
			Duration: 0,
		}

		err := service.CreatePlay(newPlay)

		assert.Error(t, err)
		assert.EqualError(t, err, "play duration must be positive")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("invalid duration negative", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		newPlay := &model.Play{
			Title:    "Ревизор",
			Author:   "Гоголь",
			Duration: -10,
		}

		err := service.CreatePlay(newPlay)

		assert.Error(t, err)
		assert.EqualError(t, err, "play duration must be positive")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		newPlay := &model.Play{
			Title:    "Ревизор",
			Author:   "Гоголь",
			Duration: 120,
		}

		mockRepo.On("Create", mock.Anything).Return(errors.New("database error"))

		err := service.CreatePlay(newPlay)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdatePlay тестирует обновление спектакля
func TestUpdatePlay(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		playID := uuid.New()
		createdAt := time.Now().Add(-24 * time.Hour)

		existingPlay := &model.Play{
			ID:        playID,
			Title:     "Старое название",
			CreatedAt: createdAt,
		}

		updatedPlay := &model.Play{
			Title:       "Новое название",
			Author:      "Новый автор",
			Description: "Новое описание",
			Duration:    150,
			Genre:       "драма",
		}

		mockRepo.On("GetByID", playID).Return(existingPlay, nil)
		mockRepo.On("Update", mock.MatchedBy(func(p *model.Play) bool {
			return p.ID == playID &&
				p.Title == "Новое название" &&
				p.CreatedAt.Equal(createdAt)
		})).Return(nil)

		err := service.UpdatePlay(playID.String(), updatedPlay)

		assert.NoError(t, err)
		assert.Equal(t, playID, updatedPlay.ID)
		assert.Equal(t, createdAt, updatedPlay.CreatedAt)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid uuid format", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		updatedPlay := &model.Play{
			Title: "Новое название",
		}

		err := service.UpdatePlay("invalid-uuid", updatedPlay)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid play ID format")
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("play not found", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		playID := uuid.New()
		updatedPlay := &model.Play{
			Title: "Новое название",
		}

		mockRepo.On("GetByID", playID).Return(nil, errors.New("not found"))

		err := service.UpdatePlay(playID.String(), updatedPlay)

		assert.Error(t, err)
		assert.EqualError(t, err, "play not found")
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("repository update error", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		playID := uuid.New()
		existingPlay := &model.Play{
			ID:    playID,
			Title: "Старое название",
		}

		updatedPlay := &model.Play{
			Title: "Новое название",
		}

		mockRepo.On("GetByID", playID).Return(existingPlay, nil)
		mockRepo.On("Update", mock.Anything).Return(errors.New("database error"))

		err := service.UpdatePlay(playID.String(), updatedPlay)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})
}

// TestDeletePlay тестирует удаление спектакля
func TestDeletePlay(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		playID := uuid.New()
		existingPlay := &model.Play{
			ID:    playID,
			Title: "Спектакль для удаления",
		}

		mockRepo.On("GetByID", playID).Return(existingPlay, nil)
		mockRepo.On("Delete", playID).Return(nil)

		err := service.DeletePlay(playID.String())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid uuid format", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		err := service.DeletePlay("invalid-uuid")

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid play ID format")
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("play not found", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		playID := uuid.New()
		mockRepo.On("GetByID", playID).Return(nil, errors.New("not found"))

		err := service.DeletePlay(playID.String())

		assert.Error(t, err)
		assert.EqualError(t, err, "play not found")
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("repository delete error", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		playID := uuid.New()
		existingPlay := &model.Play{
			ID:    playID,
			Title: "Спектакль",
		}

		mockRepo.On("GetByID", playID).Return(existingPlay, nil)
		mockRepo.On("Delete", playID).Return(errors.New("database error"))

		err := service.DeletePlay(playID.String())

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})
}

// TestServiceEdgeCases дополнительные граничные случаи
func TestServiceEdgeCases(t *testing.T) {
	t.Run("create play with minimum valid values", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		newPlay := &model.Play{
			Title:    "A",
			Author:   "B",
			Duration: 1,
		}

		mockRepo.On("Create", mock.Anything).Return(nil)

		err := service.CreatePlay(newPlay)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create play with very long strings", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		longString := string(make([]byte, 10000))
		newPlay := &model.Play{
			Title:       longString,
			Author:      longString,
			Description: longString,
			Duration:    999999,
		}

		mockRepo.On("Create", mock.Anything).Return(nil)

		err := service.CreatePlay(newPlay)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update preserves created_at timestamp", func(t *testing.T) {
		mockRepo := new(MockPlaysRepository)
		service := NewPlays(mockRepo)

		playID := uuid.New()
		originalCreatedAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

		existingPlay := &model.Play{
			ID:        playID,
			CreatedAt: originalCreatedAt,
		}

		updatedPlay := &model.Play{
			Title:     "Updated",
			CreatedAt: time.Now(), // Это значение должно быть перезаписано
		}

		mockRepo.On("GetByID", playID).Return(existingPlay, nil)
		mockRepo.On("Update", mock.MatchedBy(func(p *model.Play) bool {
			return p.CreatedAt.Equal(originalCreatedAt)
		})).Return(nil)

		err := service.UpdatePlay(playID.String(), updatedPlay)

		assert.NoError(t, err)
		assert.Equal(t, originalCreatedAt, updatedPlay.CreatedAt)
		mockRepo.AssertExpectations(t)
	})
}
