package main

import (
	"fmt"
	"log"

	"github.com/denys89/ewallet-api/config"
	"github.com/denys89/ewallet-api/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Setup database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	config.DB = db

	// Setup Gin router
	router := gin.Default()

	// Configure trusted proxies
	trustedProxies := []string{
		"127.0.0.1/8", // Localhost
		"10.0.0.0/8",  // Private network
	}

	err = router.SetTrustedProxies(trustedProxies)
	if err != nil {
		log.Printf("Warning: Failed to set trusted proxies: %v", err)
	}

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
