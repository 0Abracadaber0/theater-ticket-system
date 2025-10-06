package main

import (
	"log"
	"theater-ticket-system/internal/config"
	"theater-ticket-system/internal/database/postgres"
)

func main() {
	cfg := config.Init()
	if err := postgres.Init(cfg); err != nil {
		log.Fatal(err)
	}

	if err := postgres.Migrate(); err != nil {
		log.Fatal(err)
	}

	if err := postgres.Seed(); err != nil {
		log.Fatal(err)
	}
}
