package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"theater-ticket-system/internal/config"
	"theater-ticket-system/internal/database/postgres"
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
	api := s.router.Group("/api")
	{
		api.GET("/health-check", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		_ = postgres.Init(s.cfg)

		api.GET("/plays", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"plays": []map[string]interface{}{
					{"id": "1", "title": "Вишневый сад", "author": "А.П. Чехов"},
					{"id": "2", "title": "Ревизор", "author": "Н.В. Гоголь"},
					{"id": "3", "title": "Гамлет", "author": "У. Шекспир"},
				},
			})
		})

		api.GET("/plays/:id", func(c *gin.Context) {
			id := c.Param("id")
			// TODO: загрузка из БД
			c.JSON(200, gin.H{
				"id":     id,
				"title":  "Вишневый сад",
				"author": "А.П. Чехов",
				"date":   "2025-10-15",
				"time":   "19:00",
				"price":  1500,
			})
		})
	}

	s.router.Static("/css", "./frontend/public/css")
	s.router.Static("/js", "./frontend/public/js")
	s.router.Static("/images", "./frontend/public/images")
	s.router.StaticFile("/favicon.ico", "./frontend/public/favicon.ico")

	s.router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		filePath := "./frontend/public" + path + "/index.html"
		if path == "/" {
			filePath = "./frontend/public/index.html"
		}

		c.File(filePath)
	})
}

func (s *Server) Run() error {
	return s.router.Run("0.0.0.0:" + strconv.Itoa(s.cfg.Port))
}
