package controllers

import (
	"net/http"
	model "theater-ticket-system/internal/models/models"
	response "theater-ticket-system/internal/models/responses"

	"github.com/gin-gonic/gin"
)

type PerformancesService interface {
	GetAllPerformances(ctx *gin.Context) ([]model.Performance, error)
}

type PerformancesController struct {
	service PerformancesService
}

func NewPerformancesController(service PerformancesService) *PerformancesController {
	return &PerformancesController{service: service}
}

// GetAllPerformances godoc
// @Summary Get all performances
// @Description Get list of all performances
// @Tags performances
// @Produce json
// @Success 200 {array} response.Performance
// @Router /api/performances [get]
func (c *PerformancesController) GetAllPerformances(ctx *gin.Context) {
	performances, err := c.service.GetAllPerformances(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []response.Performance
	for _, p := range performances {
		resp = append(resp, p.Response())
	}

	ctx.JSON(http.StatusOK, resp)
}
