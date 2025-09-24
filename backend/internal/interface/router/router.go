// Package router provides HTTP routing configuration for the FinanceTracker API.
// It sets up all API endpoints and applies necessary middleware.
package router

import (
	"financetracker/internal/interface/middleware"
	"financetracker/pkg/config"
	"financetracker/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	AuthHandler interface {
		Login(*gin.Context)
		Callback(*gin.Context)
		Logout(*gin.Context)
		GetCurrentUser(*gin.Context)
		CheckAuth(*gin.Context)
		SetToken(*gin.Context)
		RemoveToken(*gin.Context)
	}
	UserHandler interface {
		GetCurrentUser(*gin.Context)
		UpdateCurrentUser(*gin.Context)
	}
	AccountHandler interface {
		List(*gin.Context)
		Create(*gin.Context)
		Get(*gin.Context)
		Update(*gin.Context)
		Delete(*gin.Context)
	}
	TransactionHandler interface {
		List(*gin.Context)
		Create(*gin.Context)
		Get(*gin.Context)
		Update(*gin.Context)
		Delete(*gin.Context)
		MonthlySummary(*gin.Context)
	}
	CategoryHandler interface {
		List(*gin.Context)
		Create(*gin.Context)
		Get(*gin.Context)
		Update(*gin.Context)
		Delete(*gin.Context)
		ListMaster(*gin.Context)
	}
}

// Router represents the HTTP router with its dependencies.
type Router struct {
	engine   *gin.Engine
	cfg      *config.Config
	logger   *logger.Logger
	handlers *Handlers
}

// New creates a new router instance with all routes configured.
func New(cfg *config.Config, logger *logger.Logger) *Router {
	return NewWithHandlers(cfg, logger, nil)
}

// NewWithHandlers creates a new router instance with handlers.
func NewWithHandlers(cfg *config.Config, logger *logger.Logger, handlers *Handlers) *Router {
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
		engine:   engine,
		cfg:      cfg,
		logger:   logger,
		handlers: handlers,
	}

	// Setup routes
	router.setupRoutes()

	return router
}

// setupRoutes configures all API routes.
func (r *Router) setupRoutes() {
	// ヘルスチェックエンドポイント（認証不要）
	r.engine.GET("/health", r.healthCheck)

	// Swagger UIエンドポイント
	r.engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 ルート
	v1 := r.engine.Group("/api/v1")
	{
		// 公開ルート（認証不要）
		public := v1.Group("")
		public.GET("/", r.apiInfo)

		// ドメイン別ルーターを初期化
		authRouter := NewAuthRouter(r.handlers)

		// 公開ルートを登録
		authRouter.RegisterRoutes(public)

		// 認証が必要なルート
		protected := v1.Group("")
		protected.Use(middleware.Auth(r.cfg))
		{
			// 認証が必要なルートを登録
			authRouter.RegisterProtectedRoutes(protected)

			// その他のルート（TODO: 各ドメインのルーターに分離）
			r.registerUserRoutes(protected)
			r.registerAccountRoutes(protected)
			r.registerTransactionRoutes(protected)
			r.registerCategoryRoutes(protected)
			r.registerBudgetRoutes(protected)
			r.registerReportRoutes(protected)
			r.registerNotificationRoutes(protected)
		}
	}

	// 404 handler
	r.engine.NoRoute(r.notFound)
}

// ユーザールート（TODO: 専用ファイルに分離）
func (r *Router) registerUserRoutes(group *gin.RouterGroup) {
	users := group.Group("/users")
	if r.handlers != nil && r.handlers.UserHandler != nil {
		users.GET("/me", r.handlers.UserHandler.GetCurrentUser)
		users.PUT("/me", r.handlers.UserHandler.UpdateCurrentUser)
	} else {
		users.GET("/me", r.notImplemented)
		users.PUT("/me", r.notImplemented)
	}
}

// アカウントルート（TODO: 専用ファイルに分離）
func (r *Router) registerAccountRoutes(group *gin.RouterGroup) {
	accounts := group.Group("/accounts")
	if r.handlers != nil && r.handlers.AccountHandler != nil {
		accounts.GET("", r.handlers.AccountHandler.List)
		accounts.POST("", r.handlers.AccountHandler.Create)
		accounts.GET("/:id", r.handlers.AccountHandler.Get)
		accounts.PUT("/:id", r.handlers.AccountHandler.Update)
		accounts.DELETE("/:id", r.handlers.AccountHandler.Delete)
		accounts.POST("/:id/movements", r.notImplemented) // TODO: accountHandler.CreateMovement
	} else {
		accounts.GET("", r.notImplemented)
		accounts.POST("", r.notImplemented)
		accounts.GET("/:id", r.notImplemented)
		accounts.PUT("/:id", r.notImplemented)
		accounts.DELETE("/:id", r.notImplemented)
		accounts.POST("/:id/movements", r.notImplemented)
	}
}

// トランザクションルート（TODO: 専用ファイルに分離）
func (r *Router) registerTransactionRoutes(group *gin.RouterGroup) {
	transactions := group.Group("/transactions")
	if r.handlers != nil && r.handlers.TransactionHandler != nil {
		transactions.GET("", r.handlers.TransactionHandler.List)
		transactions.POST("", r.handlers.TransactionHandler.Create)
		transactions.GET("/:id", r.handlers.TransactionHandler.Get)
		transactions.PUT("/:id", r.handlers.TransactionHandler.Update)
		transactions.DELETE("/:id", r.handlers.TransactionHandler.Delete)
		transactions.GET("/summary/monthly", r.handlers.TransactionHandler.MonthlySummary)
	} else {
		transactions.GET("", r.notImplemented)
		transactions.POST("", r.notImplemented)
		transactions.GET("/:id", r.notImplemented)
		transactions.PUT("/:id", r.notImplemented)
		transactions.DELETE("/:id", r.notImplemented)
		transactions.GET("/summary/monthly", r.notImplemented)
	}
}

// カテゴリールート（TODO: 専用ファイルに分離）
func (r *Router) registerCategoryRoutes(group *gin.RouterGroup) {
	categories := group.Group("/categories")
	if r.handlers != nil && r.handlers.CategoryHandler != nil {
		categories.GET("", r.handlers.CategoryHandler.List)
		categories.POST("", r.handlers.CategoryHandler.Create)
		categories.GET("/:id", r.handlers.CategoryHandler.Get)
		categories.PUT("/:id", r.handlers.CategoryHandler.Update)
		categories.DELETE("/:id", r.handlers.CategoryHandler.Delete)
		categories.GET("/master", r.handlers.CategoryHandler.ListMaster)
	} else {
		categories.GET("", r.notImplemented)
		categories.POST("", r.notImplemented)
		categories.GET("/:id", r.notImplemented)
		categories.PUT("/:id", r.notImplemented)
		categories.DELETE("/:id", r.notImplemented)
		categories.GET("/master", r.notImplemented)
	}
}

// 予算ルート（TODO: 専用ファイルに分離）
func (r *Router) registerBudgetRoutes(group *gin.RouterGroup) {
	budgets := group.Group("/budgets")
	budgets.GET("", r.notImplemented)         // TODO: budgetHandler.List
	budgets.POST("", r.notImplemented)        // TODO: budgetHandler.Create
	budgets.GET("/:id", r.notImplemented)     // TODO: budgetHandler.Get
	budgets.PUT("/:id", r.notImplemented)     // TODO: budgetHandler.Update
	budgets.DELETE("/:id", r.notImplemented)  // TODO: budgetHandler.Delete
	budgets.GET("/current", r.notImplemented) // TODO: budgetHandler.GetCurrent

	// 予算提案ルート
	suggestions := group.Group("/budget-suggestions")
	suggestions.GET("", r.notImplemented)           // TODO: suggestionHandler.List
	suggestions.POST("/generate", r.notImplemented) // TODO: suggestionHandler.Generate
}

// レポートルート（TODO: 専用ファイルに分離）
func (r *Router) registerReportRoutes(group *gin.RouterGroup) {
	reports := group.Group("/reports")
	reports.GET("/assets/snapshots", r.notImplemented)        // TODO: reportHandler.AssetSnapshots
	reports.GET("/assets/forecasts/latest", r.notImplemented) // TODO: reportHandler.LatestForecast

	// サマリールート
	summary := group.Group("/summary")
	summary.GET("/monthly", r.notImplemented) // TODO: summaryHandler.Monthly
}

// 通知設定ルート（TODO: 専用ファイルに分離）
func (r *Router) registerNotificationRoutes(group *gin.RouterGroup) {
	notifications := group.Group("/notifications")
	notifications.GET("/settings", r.notImplemented) // TODO: notificationHandler.GetSettings
	notifications.PUT("/settings", r.notImplemented) // TODO: notificationHandler.UpdateSettings
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

// healthCheck ヘルスチェックエンドポイント
// @Summary ヘルスチェック
// @Description APIの稼働状態を確認します
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{} "ヘルスチェック結果"
// @Router /health [get]
func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
		"app":    "FinanceTracker API",
	})
}

// apiInfo API情報エンドポイント
// @Summary API情報
// @Description APIの基本情報を返します
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{} "API情報"
// @Router / [get]
func (r *Router) apiInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to FinanceTracker API",
		"version": "1.0.0",
		"docs":    "/docs/index.html",
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
