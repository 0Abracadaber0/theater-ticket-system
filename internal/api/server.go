package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"theater-ticket-system/internal/api/controllers"
	"theater-ticket-system/internal/config"
	"theater-ticket-system/internal/database/postgres"
	"theater-ticket-system/internal/repository"
	"theater-ticket-system/internal/services"
)

type Server struct {
	router          *gin.Engine
	cfg             *config.Config
	playsController *controllers.PlaysController
}

func NewServer(cfg *config.Config) *Server {
	if err := postgres.Init(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	playsRepo := repository.NewPlaysRepository(postgres.DB)
	playsService := service.NewPlaysService(playsRepo)
	playsController := controllers.NewPlaysController(playsService)

	server := &Server{
		router:          gin.Default(),
		cfg:             cfg,
		playsController: playsController,
	}

	server.setupRoutes()

	return server
}

func (s *Server) Run() error {
	return s.router.Run("0.0.0.0:" + strconv.Itoa(s.cfg.Port))
}
