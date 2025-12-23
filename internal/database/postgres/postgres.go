package postgres

import (
	"fmt"
	"log"
	"theater-ticket-system/internal/config"
	"theater-ticket-system/internal/models/models"
	"time"

	"github.com/google/uuid"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Minsk",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected")
	return nil
}

func Migrate() error {
	log.Println("Running migrations...")

	err := DB.AutoMigrate(
		&model.User{},
		&model.Play{},
		&model.Hall{},
		&model.Seat{},
		&model.Performance{},
		&model.PerformanceSeat{},
		&model.Booking{},
		&model.EmailVerification{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	log.Println("Migrations completed")
	return nil
}

func Seed() error {
	log.Println("Seeding database...")

	var count int64
	DB.Model(&model.Play{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	hall := model.Hall{
		ID:       uuid.New(),
		Name:     "Большой зал",
		Capacity: 200,
	}
	if err := DB.Create(&hall).Error; err != nil {
		return err
	}

	for row := 1; row <= 10; row++ {
		for number := 1; number <= 20; number++ {
			category := "parterre"
			if row > 7 {
				category = "balcony"
			}

			seat := model.Seat{
				ID:       uuid.New(),
				HallID:   hall.ID,
				Row:      row,
				Number:   number,
				Category: category,
			}
			if err := DB.Create(&seat).Error; err != nil {
				return err
			}
		}
	}

	plays := []model.Play{
		{
			ID:          uuid.New(),
			Title:       "Вишневый сад",
			Author:      "А.П. Чехов",
			Description: "Классическая пьеса о судьбе дворянской усадьбы",
			Duration:    180,
			Genre:       "драма",
			PosterURL:   "https://shorturl.at/c2sCS",
		},
		{
			ID:          uuid.New(),
			Title:       "Ревизор",
			Author:      "Н.В. Гоголь",
			Description: "Сатирическая комедия о коррупции в провинции",
			Duration:    150,
			Genre:       "комедия",
			PosterURL:   "https://shorturl.at/c2sCS",
		},
		{
			ID:          uuid.New(),
			Title:       "Гамлет",
			Author:      "У. Шекспир",
			Description: "Трагедия о принце датском",
			Duration:    210,
			Genre:       "трагедия",
			PosterURL:   "https://shorturl.at/c2sCS",
		},
	}

	for _, play := range plays {
		if err := DB.Create(&play).Error; err != nil {
			return err
		}

		for i := 1; i <= 5; i++ {
			performance := model.Performance{
				ID:     uuid.New(),
				PlayID: play.ID,
				HallID: hall.ID,
				Date:   time.Now().AddDate(0, 0, i*7),
				Status: "scheduled",
			}
			if err := DB.Create(&performance).Error; err != nil {
				return err
			}

			var seats []model.Seat
			DB.Where("hall_id = ?", hall.ID).Find(&seats)

			for _, seat := range seats {
				price := 1500
				if seat.Category == "parterre" && seat.Row <= 5 {
					price = 3500
				} else if seat.Category == "balcony" {
					price = 1000
				}

				perfSeat := model.PerformanceSeat{
					ID:            uuid.New(),
					PerformanceID: performance.ID,
					SeatID:        seat.ID,
					Price:         price,
					Status:        "available",
				}
				if err := DB.Create(&perfSeat).Error; err != nil {
					return err
				}
			}
		}
	}

	log.Println("Database seeded successfully")
	return nil
}
