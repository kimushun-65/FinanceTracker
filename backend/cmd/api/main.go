// Package main is the entry point for the FinanceTracker API server.
// It sets up the HTTP server with all necessary middleware and routes.
package main

// @title FinanceTracker API
// @version 1.0
// @description Finance management REST API with Auth0 authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@financetracker.local

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	stdlog "log"

	"financetracker/internal/di"
	"financetracker/internal/interface/router"
	"financetracker/pkg/config"
	"financetracker/pkg/logger"

	_ "financetracker/docs" // Swagger生成ドキュメント

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		stdlog.Println("No .env file found")
	}

	// Initialize logger
	log := logger.NewLogger("api")
	defer func() {
		if err := log.Sync(); err != nil {
			// Use standard log since our logger might not be available
			stdlog.Printf("Failed to sync logger: %v", err)
		}
	}()

	// Initialize DI container
	container, err := di.NewContainer()
	if err != nil {
		log.Error("Failed to initialize DI container: " + err.Error())
		return
	}
	defer func() {
		if err := container.Close(); err != nil {
			log.Error("Failed to close DI container: " + err.Error())
		}
	}()

	// Load configuration for router (router expects pkg/config.Config)
	cfg := config.Load()

	// Create router with handlers from DI container
	r := router.NewWithHandlers(cfg, log, &router.Handlers{
		AuthHandler:        container.AuthHandler,
		UserHandler:        container.UserHandler,
		AccountHandler:     container.AccountHandler,
		TransactionHandler: container.TransactionHandler,
		CategoryHandler:    container.CategoryHandler,
	})

	// Start server
	if err := r.Run(); err != nil {
		log.Error("Failed to start server: " + err.Error())
		return
	}
}
