package controllers

import (
	"net/http"
	model "theater-ticket-system/internal/models/models"
	request "theater-ticket-system/internal/models/requests"
	response "theater-ticket-system/internal/models/responses"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingsService interface {
	CreateBooking(phone, name string, performanceID uuid.UUID, seatIDs []uuid.UUID) (*model.Booking, error)
	GetBookingByID(id string) (*model.Booking, error)
	GetUserBookings(phone string) ([]model.Booking, error)
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
// @Description Create a new booking for selected seats. User will be created or found by phone number
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body request.CreateBooking true "Booking object"
// @Success 201 {object} response.Booking
// @Router /api/bookings [post]
func (c *BookingsController) CreateBooking(ctx *gin.Context) {
	var req request.CreateBooking
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := c.service.CreateBooking(req.Phone, req.Name, req.PerformanceID, req.SeatIDs)
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
// @Description Get booking history for a user by phone number
// @Tags bookings
// @Produce json
// @Param phone query string true "User phone number"
// @Success 200 {array} response.Booking
// @Router /api/bookings [get]
func (c *BookingsController) GetUserBookings(ctx *gin.Context) {
	phone := ctx.Query("phone")
	if phone == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "phone is required"})
		return
	}

	bookings, err := c.service.GetUserBookings(phone)
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
