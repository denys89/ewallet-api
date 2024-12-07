package routes

import (
	"net/http"

	"github.com/denys89/ewallet-api/config"
	"github.com/denys89/ewallet-api/middleware"
	"github.com/denys89/ewallet-api/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Address   string `json:"address" binding:"required"`
}

// GetProfile handles retrieving user profile information
func GetProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	userRepo := repositories.NewUserRepository(config.DB)
	user, err := userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"user_id":      user.ID,
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"phone_number": user.PhoneNumber,
			"address":      user.Address,
		},
	})
}

// UpdateProfile handles updating user profile information
func UpdateProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repositories.NewUserRepository(config.DB)
	user, err := userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Address != "" {
		user.Address = req.Address
	}

	if err := userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"user_id":      user.ID,
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"phone_number": user.PhoneNumber,
			"address":      user.Address,
			"updated_at":   user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// GetBalance handles retrieving user balance
func GetBalance(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	userRepo := repositories.NewUserRepository(config.DB)
	user, err := userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"balance": user.Balance,
		},
	})
}
