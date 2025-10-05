package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"theater-ticket-system/internal/config"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
}

func NewServer(cfg *config.Config) *Server {
	server := &Server{
		router: gin.Default(),
		cfg:    cfg,
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	s.router.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func (s *Server) Run() error {
	return s.router.Run("0.0.0.0:" + strconv.Itoa(s.cfg.Port))
}
