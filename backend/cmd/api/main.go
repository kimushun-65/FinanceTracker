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

	"financetracker/internal/application/service"
	"financetracker/internal/infrastructure/auth0"
	gormRepo "financetracker/internal/infrastructure/gorm/repository"
	"financetracker/internal/interface/handler"
	"financetracker/internal/interface/router"
	"financetracker/pkg/config"
	"financetracker/pkg/database"
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
	log := logger.New()
	defer func() {
		if err := log.Sync(); err != nil {
			// Use standard log since our logger might not be available
			stdlog.Printf("Failed to sync logger: %v", err)
		}
	}()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.New(cfg, log)
	if err != nil {
		log.Error("Failed to connect to database: " + err.Error())
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Failed to close database: " + err.Error())
		}
	}()

	// Initialize Auth0
	auth0Client := auth0.NewClient(cfg.Auth0Domain, cfg.Auth0ClientID, cfg.Auth0Audience)
	authMiddleware := auth0.NewAuthMiddleware(auth0Client)

	// Initialize repositories
	userRepo := gormRepo.NewUserRepository(db.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, auth0Client, authMiddleware, cfg.Auth0ClientID, cfg.Auth0CallbackURL)

	// Create router with handlers
	r := router.NewWithHandlers(cfg, log, &router.Handlers{
		AuthHandler: authHandler,
	})

	// Start server
	if err := r.Run(); err != nil {
		log.Error("Failed to start server: " + err.Error())
		return
	}
}
