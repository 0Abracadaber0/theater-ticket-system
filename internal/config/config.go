package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port int
}

func Init() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
	return &Config{
		Port: port,
	}
}
