package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strings"
	_ "theater-ticket-system/docs"
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
			plays.GET("", s.playsController.GetAllPlays)
			plays.GET("/:id", s.playsController.GetPlayByID)
			plays.POST("", s.playsController.CreatePlay)
			plays.PUT("/:id", s.playsController.UpdatePlay)
			plays.DELETE("/:id", s.playsController.DeletePlay)
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
