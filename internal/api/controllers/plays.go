package controllers

import (
	"net/http"
	"theater-ticket-system/internal/models/requests"
	"theater-ticket-system/internal/models/responses"
	"theater-ticket-system/internal/services"

	"github.com/gin-gonic/gin"
)

type PlaysController struct {
	service service.PlaysService
}

func NewPlaysController(service service.PlaysService) *PlaysController {
	return &PlaysController{service: service}
}

// GetAllPlays godoc
// @Summary Get all plays
// @Description Get list of all plays
// @Tags plays
// @Produce json
// @Success 200 {array} response.Play
// @Router /api/plays [get]
func (pc *PlaysController) GetAllPlays(c *gin.Context) {
	plays, err := pc.service.GetAllPlays()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]response.Play, len(plays))
	for i := range plays {
		resp[i] = plays[i].Response()
	}

	c.JSON(http.StatusOK, resp)
}

// GetPlayByID godoc
// @Summary Get play by ID
// @Description Get detailed information about a play
// @Tags plays
// @Produce json
// @Param id path string true "Play ID"
// @Success 200 {object} response.Play
// @Router /api/plays/{id} [get]
func (pc *PlaysController) GetPlayByID(c *gin.Context) {
	id := c.Param("id")

	play, err := pc.service.GetPlayByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, play.Response())
}

// CreatePlay godoc
// @Summary Create new play
// @Description Create a new play
// @Tags plays
// @Accept json
// @Produce json
// @Param play body request.Play true "Play object"
// @Success 201 {object} response.Play
// @Router /api/plays [post]
func (pc *PlaysController) CreatePlay(c *gin.Context) {
	// todo отдавать на выход модель (чтобы был id и timestamps)
	var req request.Play
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.service.CreatePlay(req.Model()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req.Model().Response())
}

// UpdatePlay godoc
// @Summary Update play
// @Description Update an existing play
// @Tags plays
// @Accept json
// @Produce json
// @Param id path string true "Play ID"
// @Param play body request.Play true "Play object"
// @Success 200 {object} response.Play
// @Router /api/plays/{id} [put]
func (pc *PlaysController) UpdatePlay(c *gin.Context) {
	// todo отдавать на выход модель (чтобы был id и timestamps)
	id := c.Param("id")

	var req request.Play
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.service.UpdatePlay(id, req.Model()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req.Model().Response())
}

// DeletePlay godoc
// @Summary Delete play
// @Description Delete a play
// @Tags plays
// @Produce json
// @Param id path string true "Play ID"
// @Success 204
// @Router /api/plays/{id} [delete]
func (pc *PlaysController) DeletePlay(c *gin.Context) {
	id := c.Param("id")

	if err := pc.service.DeletePlay(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
