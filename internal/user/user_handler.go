package user

import (
	"net/http"
	"otas/models"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name       string            `json:"name"        binding:"required"`
	Email      string            `json:"email"       binding:"required,email"`
	Phone      string            `json:"phone"       binding:"required"`
	Password   string            `json:"password"    binding:"required,min=6"`
	SavingType models.SavingType `json:"saving_type" binding:"required"`
	DailyLimit models.DailyLimit `json:"daily_limit" binding:"required"`
}

type Handler struct {
	service *userService
}

func NewUserHandler(service *userService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(&models.User{
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   req.Password,
		SavingType: req.SavingType,
		DailyLimit: req.DailyLimit,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "account created successfully",
		"user":    user,
	})
}
