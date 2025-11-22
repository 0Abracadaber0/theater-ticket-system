package controllers

import (
	"net/http"
	model "theater-ticket-system/internal/models/models"
	response "theater-ticket-system/internal/models/responses"
	"time"

	"github.com/gin-gonic/gin"
)

type PerformancesService interface {
	GetAllPerformances(playID *string, dateFrom, dateTo *time.Time) ([]model.Performance, error)
	GetPerformanceByID(id string) (*model.Performance, error)
	GetPerformanceSeats(id string) ([]model.PerformanceSeat, error)
}

type PerformancesController struct {
	service PerformancesService
}

func NewPerformancesController(service PerformancesService) *PerformancesController {
	return &PerformancesController{service: service}
}

// GetAllPerformances godoc
// @Summary Get all performances
// @Description Get list of all performances with optional filters
// @Tags performances
// @Produce json
// @Param play_id query string false "Filter by play ID"
// @Param date_from query string false "Filter by date from (RFC3339)"
// @Param date_to query string false "Filter by date to (RFC3339)"
// @Success 200 {array} response.Performance
// @Router /api/performances [get]
func (c *PerformancesController) GetAllPerformances(ctx *gin.Context) {
	playID := ctx.Query("play_id")
	dateFromStr := ctx.Query("date_from")
	dateToStr := ctx.Query("date_to")

	var playIDPtr *string
	if playID != "" {
		playIDPtr = &playID
	}

	var dateFrom, dateTo *time.Time
	if dateFromStr != "" {
		parsed, err := time.Parse(time.RFC3339, dateFromStr)
		if err == nil {
			dateFrom = &parsed
		}
	}
	if dateToStr != "" {
		parsed, err := time.Parse(time.RFC3339, dateToStr)
		if err == nil {
			dateTo = &parsed
		}
	}

	performances, err := c.service.GetAllPerformances(playIDPtr, dateFrom, dateTo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]response.Performance, len(performances))
	for i := range performances {
		resp[i] = performances[i].Response()
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetPerformanceByID godoc
// @Summary Get performance by ID
// @Description Get detailed information about a performance
// @Tags performances
// @Produce json
// @Param id path string true "Performance ID"
// @Success 200 {object} response.Performance
// @Router /api/performances/{id} [get]
func (c *PerformancesController) GetPerformanceByID(ctx *gin.Context) {
	id := ctx.Param("id")

	performance, err := c.service.GetPerformanceByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, performance.Response())
}

// GetPerformanceSeats godoc
// @Summary Get performance seats
// @Description Get hall layout with seat availability for a performance
// @Tags performances
// @Produce json
// @Param id path string true "Performance ID"
// @Success 200 {array} response.PerformanceSeat
// @Router /api/performances/{id}/seats [get]
func (c *PerformancesController) GetPerformanceSeats(ctx *gin.Context) {
	id := ctx.Param("id")

	seats, err := c.service.GetPerformanceSeats(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]response.PerformanceSeat, len(seats))
	for i := range seats {
		resp[i] = seats[i].Response()
	}

	ctx.JSON(http.StatusOK, resp)
}
