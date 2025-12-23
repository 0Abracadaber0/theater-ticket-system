package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  int
	DB    DBConfig
	Email EmailConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type EmailConfig struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort string
}

func Init() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		log.Fatal("Invalid PORT:", err)
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatal("Invalid DB_PORT:", err)
	}

	return &Config{
		Port: port,
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "theater_tickets"),
		},
		Email: EmailConfig{
			From:     getEnv("EMAIL_FROM", ""),
			Password: getEnv("EMAIL_PASSWORD", ""),
			SMTPHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort: getEnv("SMTP_PORT", "587"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
