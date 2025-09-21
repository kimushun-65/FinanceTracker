// Package router provides HTTP routing configuration for the FinanceTracker API.
// It sets up all API endpoints and applies necessary middleware.
package router

import (
	"financetracker/internal/interface/middleware"
	"financetracker/pkg/config"
	"financetracker/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Router represents the HTTP router with its dependencies.
type Router struct {
	engine *gin.Engine
	cfg    *config.Config
	logger *logger.Logger
}

// New creates a new router instance with all routes configured.
func New(cfg *config.Config, logger *logger.Logger) *Router {
	// Set Gin mode based on environment
	gin.SetMode(cfg.GinMode)

	// Create Gin engine with default middleware
	engine := gin.New()

	// Apply global middleware
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logger(logger))
	engine.Use(middleware.ErrorHandler(logger))
	engine.Use(middleware.RequestID())
	engine.Use(middleware.CORS(cfg))

	router := &Router{
		engine: engine,
		cfg:    cfg,
		logger: logger,
	}

	// Setup routes
	router.setupRoutes()

	return router
}

// setupRoutes configures all API routes.
func (r *Router) setupRoutes() {
	// Health check endpoint (no auth required)
	r.engine.GET("/health", r.healthCheck)

	// API v1 routes
	v1 := r.engine.Group("/api/v1")
	{
		// Public routes (no auth required)
		public := v1.Group("")
		public.GET("/", r.apiInfo)
		// TODO: Add auth callback endpoints
		// public.POST("/auth/callback", authHandler.Callback)

		// Protected routes (auth required)
		protected := v1.Group("")
		protected.Use(middleware.Auth(r.cfg))
		{
			// User routes
			users := protected.Group("/users")
			users.GET("/me", r.notImplemented) // TODO: userHandler.GetCurrentUser
			users.PUT("/me", r.notImplemented) // TODO: userHandler.UpdateCurrentUser

			// Account routes
			accounts := protected.Group("/accounts")
			accounts.GET("", r.notImplemented)                // TODO: accountHandler.List
			accounts.POST("", r.notImplemented)               // TODO: accountHandler.Create
			accounts.GET("/:id", r.notImplemented)            // TODO: accountHandler.Get
			accounts.PUT("/:id", r.notImplemented)            // TODO: accountHandler.Update
			accounts.DELETE("/:id", r.notImplemented)         // TODO: accountHandler.Delete
			accounts.POST("/:id/movements", r.notImplemented) // TODO: accountHandler.CreateMovement

			// Transaction routes
			transactions := protected.Group("/transactions")
			transactions.GET("", r.notImplemented)                 // TODO: transactionHandler.List
			transactions.POST("", r.notImplemented)                // TODO: transactionHandler.Create
			transactions.GET("/:id", r.notImplemented)             // TODO: transactionHandler.Get
			transactions.PUT("/:id", r.notImplemented)             // TODO: transactionHandler.Update
			transactions.DELETE("/:id", r.notImplemented)          // TODO: transactionHandler.Delete
			transactions.GET("/summary/monthly", r.notImplemented) // TODO: transactionHandler.MonthlySummary

			// Category routes
			categories := protected.Group("/categories")
			categories.GET("", r.notImplemented)        // TODO: categoryHandler.List
			categories.PUT("/:id", r.notImplemented)    // TODO: categoryHandler.Update
			categories.DELETE("/:id", r.notImplemented) // TODO: categoryHandler.Delete
			categories.GET("/master", r.notImplemented) // TODO: categoryHandler.ListMaster

			// Budget routes
			budgets := protected.Group("/budgets")
			budgets.GET("", r.notImplemented)         // TODO: budgetHandler.List
			budgets.POST("", r.notImplemented)        // TODO: budgetHandler.Create
			budgets.GET("/:id", r.notImplemented)     // TODO: budgetHandler.Get
			budgets.PUT("/:id", r.notImplemented)     // TODO: budgetHandler.Update
			budgets.DELETE("/:id", r.notImplemented)  // TODO: budgetHandler.Delete
			budgets.GET("/current", r.notImplemented) // TODO: budgetHandler.GetCurrent

			// Budget suggestion routes
			suggestions := protected.Group("/budget-suggestions")
			suggestions.GET("", r.notImplemented)           // TODO: suggestionHandler.List
			suggestions.POST("/generate", r.notImplemented) // TODO: suggestionHandler.Generate

			// Report routes
			reports := protected.Group("/reports")
			reports.GET("/assets/snapshots", r.notImplemented)        // TODO: reportHandler.AssetSnapshots
			reports.GET("/assets/forecasts/latest", r.notImplemented) // TODO: reportHandler.LatestForecast

			// Summary routes
			summary := protected.Group("/summary")
			summary.GET("/monthly", r.notImplemented) // TODO: summaryHandler.Monthly

			// Notification settings
			notifications := protected.Group("/notifications")
			notifications.GET("/settings", r.notImplemented) // TODO: notificationHandler.GetSettings
			notifications.PUT("/settings", r.notImplemented) // TODO: notificationHandler.UpdateSettings
		}
	}

	// 404 handler
	r.engine.NoRoute(r.notFound)
}

// Engine returns the underlying Gin engine.
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// Run starts the HTTP server.
func (r *Router) Run() error {
	addr := ":" + r.cfg.AppPort
	r.logger.Info("Starting server on " + addr)
	return r.engine.Run(addr)
}

// Handler functions

func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
		"app":    "FinanceTracker API",
	})
}

func (r *Router) apiInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to FinanceTracker API",
		"version": "1.0.0",
		"docs":    "/api/v1/docs", // TODO: Add Swagger docs
	})
}

func (r *Router) notImplemented(c *gin.Context) {
	c.JSON(501, gin.H{
		"code":    "NOT_IMPLEMENTED",
		"message": "This endpoint is not yet implemented",
	})
}

func (r *Router) notFound(c *gin.Context) {
	c.JSON(404, gin.H{
		"code":    "NOT_FOUND",
		"message": "The requested resource was not found",
		"path":    c.Request.URL.Path,
	})
}
