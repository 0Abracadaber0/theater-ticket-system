package main

import (
	"log"
	"theater-ticket-system/internal/api"
	"theater-ticket-system/internal/config"
)

func main() {
	cfg := config.Init()

	server := api.NewServer(cfg)

	log.Fatal(server.Run())
}
