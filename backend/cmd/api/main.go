// Package main is the entry point for the FinanceTracker API server.
// It sets up the HTTP server with all necessary middleware and routes.
package main

import (
	"fmt"
	stdlog "log"

	"financetracker/pkg/config"
	"financetracker/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		stdlog.Println("No .env file found")
	}

	// Initialize logger
	log := logger.New()
	defer func() {
		if err := log.Sync(); err != nil {
			// Use standard log since our logger might not be available
			stdlog.Printf("Failed to sync logger: %v", err)
		}
	}()

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
	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to FinanceTracker API",
			"version": "1.0.0",
		})
	})

	// Start server
	port := fmt.Sprintf(":%s", cfg.AppPort)
	log.Info("Starting server on port", port)
	if err := router.Run(port); err != nil {
		log.Error("Failed to start server", err)
		return
	}
}
