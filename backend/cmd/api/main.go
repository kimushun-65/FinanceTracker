package main

import (
	"fmt"
	"log"

	"financetracker/pkg/config"
	"financetracker/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize logger
	log := logger.New()
	defer log.Sync()

	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Create Gin router
	router := gin.New()

	// Middlewares
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"app":    "FinanceTracker API",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome to FinanceTracker API",
				"version": "1.0.0",
			})
		})
	}

	// Start server
	port := fmt.Sprintf(":%s", cfg.AppPort)
	log.Info("Starting server on port", port)
	
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server", err)
	}
}