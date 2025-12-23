package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	SendVerificationCode(email string) error
	VerifyCode(email, code string) (bool, error)
}

type AuthController struct {
	service AuthService
}

func NewAuthController(service AuthService) *AuthController {
	return &AuthController{service: service}
}

// SendCode godoc
// @Summary Send verification code
// @Description Send verification code to email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{email=string} true "Email"
// @Success 200 {object} object{message=string}
// @Router /api/auth/send-code [post]
func (c *AuthController) SendCode(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	if err := c.service.SendVerificationCode(req.Email); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Verification code sent to email"})
}

// VerifyCode godoc
// @Summary Verify code
// @Description Verify email with code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{email=string,code=string} true "Email and code"
// @Success 200 {object} object{verified=boolean}
// @Router /api/auth/verify-code [post]
func (c *AuthController) VerifyCode(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	verified, err := c.service.VerifyCode(req.Email, req.Code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"verified": verified})
}
