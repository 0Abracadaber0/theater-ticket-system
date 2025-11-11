package api

import (
	"strings"
	_ "theater-ticket-system/docs"
	"theater-ticket-system/internal/api/controllers"
	"theater-ticket-system/internal/database/postgres"
	"theater-ticket-system/internal/repository"
	service "theater-ticket-system/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) setupRoutes() {
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := s.router.Group("/api")
	{
		api.GET("/health-check", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		plays := api.Group("/plays")
		{
			playsRepo := repository.NewPlays(postgres.DB)
			playsService := service.NewPlays(playsRepo)
			playsController := controllers.NewPlays(playsService)

			plays.GET("", playsController.GetAllPlays)
			plays.GET("/:id", playsController.GetPlayByID)
			plays.POST("", playsController.CreatePlay)
			plays.PUT("/:id", playsController.UpdatePlay)
			plays.DELETE("/:id", playsController.DeletePlay)
		}
	}

	s.router.Static("/css", "./frontend/public/css")
	s.router.Static("/js", "./frontend/public/js")
	s.router.Static("/images", "./frontend/public/images")
	s.router.StaticFile("/favicon.ico", "./frontend/public/favicon.ico")

	s.router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasPrefix(path, "/api") {
			c.JSON(404, gin.H{"error": "api endpoint not found"})
			return
		}

		filePath := "./frontend/public" + path + "/index.html"
		if path == "/" {
			filePath = "./frontend/public/index.html"
		}

		c.File(filePath)
	})
}
