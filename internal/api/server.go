package api

import (
	"log"
	"strconv"
	"theater-ticket-system/internal/config"
	"theater-ticket-system/internal/database/postgres"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
}

func NewServer(cfg *config.Config) *Server {
	if err := postgres.Init(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	server := &Server{
		router: gin.Default(),
		cfg:    cfg,
	}

	server.setupRoutes()

	return server
}

func (s *Server) Run() error {
	return s.router.Run("0.0.0.0:" + strconv.Itoa(s.cfg.Port))
}
