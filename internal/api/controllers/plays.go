package controllers

import (
	"net/http"
	model "theater-ticket-system/internal/models/models"
	"theater-ticket-system/internal/models/requests"
	"theater-ticket-system/internal/models/responses"

	"github.com/gin-gonic/gin"
)

type PlaysService interface {
	GetAllPlays() ([]model.Play, error)
	GetPlayByID(id string) (*model.Play, error)
	CreatePlay(play *model.Play) error
	UpdatePlay(id string, play *model.Play) error
	DeletePlay(id string) error
}

type Plays struct {
	service PlaysService
}

func NewPlays(service PlaysService) *Plays {
	return &Plays{service: service}
}

// GetAllPlays godoc
// @Summary Get all plays
// @Description Get list of all plays
// @Tags plays
// @Produce json
// @Success 200 {array} response.Play
// @Router /api/plays [get]
func (c *Plays) GetAllPlays(ctx *gin.Context) {
	plays, err := c.service.GetAllPlays()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]response.Play, len(plays))
	for i := range plays {
		resp[i] = plays[i].Response()
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetPlayByID godoc
// @Summary Get play by ID
// @Description Get detailed information about a play
// @Tags plays
// @Produce json
// @Param id path string true "Play ID"
// @Success 200 {object} response.Play
// @Router /api/plays/{id} [get]
func (c *Plays) GetPlayByID(ctx *gin.Context) {
	id := ctx.Param("id")

	play, err := c.service.GetPlayByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, play.Response())
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
func (c *Plays) CreatePlay(ctx *gin.Context) {
	// todo отдавать на выход модель (чтобы был id и timestamps)
	var req request.Play
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreatePlay(req.Model()); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, req.Model().Response())
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
func (c *Plays) UpdatePlay(ctx *gin.Context) {
	// todo отдавать на выход модель (чтобы был id и timestamps)
	id := ctx.Param("id")

	var req request.Play
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdatePlay(id, req.Model()); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, req.Model().Response())
}

// DeletePlay godoc
// @Summary Delete play
// @Description Delete a play
// @Tags plays
// @Produce json
// @Param id path string true "Play ID"
// @Success 204
// @Router /api/plays/{id} [delete]
func (c *Plays) DeletePlay(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeletePlay(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
