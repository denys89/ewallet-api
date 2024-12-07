package routes

import (
	"github.com/denys89/ewallet-api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth routes (no authentication required)
		v1.POST("/auth/register", Register)
		v1.POST("/auth/login", Login)
		v1.POST("/auth/refresh-token", RefreshToken)

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			protected.GET("/user/profile", GetProfile)
			protected.PUT("/user/profile", UpdateProfile)
			protected.GET("/user/balance", GetBalance)

			// Transaction routes
			protected.GET("/transactions", GetTransactionHistory)
			protected.POST("/transactions/topup", TopUp)
			protected.POST("/transactions/transfer", Transfer)
			protected.POST("/transactions/payment", Payment)
		}
	}
}
