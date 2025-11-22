package controllers

import (
	"net/http"
	model "theater-ticket-system/internal/models/models"
	response "theater-ticket-system/internal/models/responses"

	"github.com/gin-gonic/gin"
)

type SeatsService interface {
	GetSeatsByHallID(hallID string) ([]model.Seat, error)
}

type SeatsController struct {
	service SeatsService
}

func NewSeatsController(service SeatsService) *SeatsController {
	return &SeatsController{service: service}
}

// GetHallSeats godoc
// @Summary Get hall seats
// @Description Get static hall layout (not tied to a specific performance)
// @Tags seats
// @Produce json
// @Param id path string true "Hall ID"
// @Success 200 {array} response.Seat
// @Router /api/halls/{id}/seats [get]
func (c *SeatsController) GetHallSeats(ctx *gin.Context) {
	id := ctx.Param("id")

	seats, err := c.service.GetSeatsByHallID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]response.Seat, len(seats))
	for i := range seats {
		resp[i] = seats[i].Response()
	}

	ctx.JSON(http.StatusOK, resp)
}
