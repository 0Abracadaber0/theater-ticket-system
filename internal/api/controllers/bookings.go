package controllers

import (
	"net/http"
	model "theater-ticket-system/internal/models/models"
	response "theater-ticket-system/internal/models/responses"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingsService interface {
	CreateBooking(email, name string, performanceID uuid.UUID, seatIDs []uuid.UUID) (*model.Booking, error)
	GetBookingByID(id string) (*model.Booking, error)
	GetUserBookings(email string) ([]model.Booking, error)
	CancelBooking(id string) error
}

type BookingsController struct {
	service BookingsService
}

func NewBookingsController(service BookingsService) *BookingsController {
	return &BookingsController{service: service}
}

// CreateBooking godoc
// @Summary Create booking
// @Description Create a new booking for selected seats. User will be created or found by email
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body object{email=string,name=string,performance_id=string,seat_ids=[]string} true "Booking object"
// @Success 201 {object} response.Booking
// @Router /api/bookings [post]
func (c *BookingsController) CreateBooking(ctx *gin.Context) {
	var req struct {
		Email         string      `json:"email" binding:"required,email"`
		Name          string      `json:"name" binding:"required"`
		PerformanceID uuid.UUID   `json:"performance_id" binding:"required"`
		SeatIDs       []uuid.UUID `json:"seat_ids" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := c.service.CreateBooking(req.Email, req.Name, req.PerformanceID, req.SeatIDs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, booking.Response())
}

// GetBookingByID godoc
// @Summary Get booking by ID
// @Description Get detailed information about a booking
// @Tags bookings
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} response.Booking
// @Router /api/bookings/{id} [get]
func (c *BookingsController) GetBookingByID(ctx *gin.Context) {
	id := ctx.Param("id")

	booking, err := c.service.GetBookingByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, booking.Response())
}

// GetUserBookings godoc
// @Summary Get user bookings
// @Description Get booking history for a user by email
// @Tags bookings
// @Produce json
// @Param email query string true "User email"
// @Success 200 {array} response.Booking
// @Router /api/bookings [get]
func (c *BookingsController) GetUserBookings(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	bookings, err := c.service.GetUserBookings(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]response.Booking, len(bookings))
	for i := range bookings {
		resp[i] = bookings[i].Response()
	}

	ctx.JSON(http.StatusOK, resp)
}

// CancelBooking godoc
// @Summary Cancel booking
// @Description Cancel an existing booking (only pending bookings can be cancelled)
// @Tags bookings
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} response.Booking
// @Router /api/bookings/{id}/cancel [patch]
func (c *BookingsController) CancelBooking(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.CancelBooking(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := c.service.GetBookingByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, booking.Response())
}
