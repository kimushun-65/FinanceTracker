# DI層実装プラン

## 実装手順とディレクトリ構成

### 新しいディレクトリ構成

```
backend/internal/
├── di/                           # 🆕 DI（Dependency Injection）層
│   ├── container.go              # メインDIコンテナ
│   ├── config.go                 # 統一設定管理
│   ├── interfaces.go             # DIコンテナインターフェース
│   └── providers/                # プロバイダー群（責務分割）
│       ├── database.go           # データベース・外部サービス
│       ├── repository.go         # リポジトリ初期化
│       ├── service.go            # アプリケーションサービス
│       ├── handler.go            # HTTPハンドラー
│       └── middleware.go         # ミドルウェア
├── domain/                       # ドメイン層（既存）
├── application/                  # アプリケーション層（既存）
├── infrastructure/               # インフラストラクチャ層（既存）
└── interface/                    # インターフェース層（既存）
```

## 段階的実装計画

### Phase 1: DI基盤構築 📅 Day 1-2

#### 1.1 ディレクトリとファイル作成
```bash
mkdir -p backend/internal/di/providers
touch backend/internal/di/container.go
touch backend/internal/di/config.go
touch backend/internal/di/interfaces.go
touch backend/internal/di/providers/database.go
touch backend/internal/di/providers/repository.go
touch backend/internal/di/providers/service.go
touch backend/internal/di/providers/handler.go
touch backend/internal/di/providers/middleware.go
```

#### 1.2 基本設定管理実装
**ファイル**: `internal/di/config.go`
```go
package di

import (
    "os"
    "strconv"
    "time"
    "fmt"
)

// Config holds all application configuration
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Auth0    Auth0Config
    JWT      JWTConfig
    Log      LogConfig
}

type ServerConfig struct {
    Port         string
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    Environment  string
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
    MaxOpenConns int
    MaxIdleConns int
}

type Auth0Config struct {
    Domain       string
    ClientID     string
    ClientSecret string
    CallbackURL  string
    Audience     string
}

type JWTConfig struct {
    Secret      string
    ExpiryHours time.Duration
}

type LogConfig struct {
    Level  string
    Format string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
    config := &Config{
        Server: ServerConfig{
            Port:         getEnv("SERVER_PORT", "8080"),
            ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", "10s"),
            WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", "10s"),
            Environment:  getEnv("ENVIRONMENT", "development"),
        },
        Database: DatabaseConfig{
            Host:         getEnv("DB_HOST", "localhost"),
            Port:         getIntEnv("DB_PORT", 5432),
            User:         getEnv("DB_USER", "postgres"),
            Password:     getEnv("DB_PASSWORD", ""),
            DBName:       getEnv("DB_NAME", "financetracker"),
            SSLMode:      getEnv("DB_SSLMODE", "disable"),
            MaxOpenConns: getIntEnv("DB_MAX_OPEN_CONNS", 25),
            MaxIdleConns: getIntEnv("DB_MAX_IDLE_CONNS", 5),
        },
        Auth0: Auth0Config{
            Domain:       getEnv("AUTH0_DOMAIN", ""),
            ClientID:     getEnv("AUTH0_CLIENT_ID", ""),
            ClientSecret: getEnv("AUTH0_CLIENT_SECRET", ""),
            CallbackURL:  getEnv("AUTH0_CALLBACK_URL", "http://localhost:3000/auth/callback"),
            Audience:     getEnv("AUTH0_AUDIENCE", ""),
        },
        JWT: JWTConfig{
            Secret:      getEnv("JWT_SECRET", "development-secret-key"),
            ExpiryHours: getDurationEnv("JWT_EXPIRY_HOURS", "24h"),
        },
        Log: LogConfig{
            Level:  getEnv("LOG_LEVEL", "info"),
            Format: getEnv("LOG_FORMAT", "json"),
        },
    }
    
    return config, config.validate()
}

func (c *Config) validate() error {
    if c.Auth0.Domain == "" {
        return fmt.Errorf("AUTH0_DOMAIN is required")
    }
    if c.Auth0.ClientID == "" {
        return fmt.Errorf("AUTH0_CLIENT_ID is required")
    }
    if c.Auth0.ClientSecret == "" {
        return fmt.Errorf("AUTH0_CLIENT_SECRET is required")
    }
    return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getDurationEnv(key string, defaultValue string) time.Duration {
    value := getEnv(key, defaultValue)
    if duration, err := time.ParseDuration(value); err == nil {
        return duration
    }
    // Fallback to default
    if duration, err := time.ParseDuration(defaultValue); err == nil {
        return duration
    }
    return 30 * time.Second // Last resort default
}
```

#### 1.3 コンテナインターフェース定義
**ファイル**: `internal/di/interfaces.go`
```go
package di

import (
    "context"
    "gorm.io/gorm"
)

// ContainerInterface defines the interface for dependency injection container
type ContainerInterface interface {
    // Configuration
    GetConfig() *Config
    
    // Database
    GetDB() *gorm.DB
    
    // Graceful shutdown
    Shutdown(ctx context.Context) error
}

// HealthChecker interface for health checks
type HealthChecker interface {
    CheckHealth(ctx context.Context) error
}

// Initializer interface for components that need initialization
type Initializer interface {
    Initialize() error
}

// Cleaner interface for components that need cleanup
type Cleaner interface {
    Cleanup() error
}
```

#### 1.4 メインコンテナ構造定義
**ファイル**: `internal/di/container.go` (初期版)
```go
package di

import (
    "context"
    "fmt"
    "log"
    
    "gorm.io/gorm"
    
    // Domain repository interfaces (既存)
    userRepo "github.com/your-org/financetracker/internal/domain/user/repository"
    
    // Infrastructure (既存)
    "github.com/your-org/financetracker/internal/infrastructure/auth0"
    
    // Application services (実装予定)
    "github.com/your-org/financetracker/internal/application/service"
    
    // Interface handlers (実装予定)
    "github.com/your-org/financetracker/internal/interface/handler"
    "github.com/your-org/financetracker/internal/interface/middleware"
)

// Container holds all dependencies
type Container struct {
    // Configuration
    config *Config
    
    // Infrastructure
    db          *gorm.DB
    auth0Client *auth0.Client
    
    // Repositories (domain interfaces)
    userRepo userRepo.UserRepository
    // TODO: Add other repositories as they are implemented
    
    // Services
    authService service.AuthService
    // TODO: Add other services as they are implemented
    
    // Handlers
    authHandler *handler.AuthHandler
    // TODO: Add other handlers as they are implemented
    
    // Middleware
    authMiddleware *middleware.AuthMiddleware
    // TODO: Add other middleware as needed
}

// NewContainer creates and initializes a new DI container
func NewContainer() (*Container, error) {
    // Load configuration
    config, err := LoadConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }
    
    container := &Container{
        config: config,
    }
    
    // Initialize in dependency order
    if err := container.initInfrastructure(); err != nil {
        return nil, fmt.Errorf("failed to initialize infrastructure: %w", err)
    }
    
    if err := container.initRepositories(); err != nil {
        return nil, fmt.Errorf("failed to initialize repositories: %w", err)
    }
    
    if err := container.initServices(); err != nil {
        return nil, fmt.Errorf("failed to initialize services: %w", err)
    }
    
    if err := container.initHandlers(); err != nil {
        return nil, fmt.Errorf("failed to initialize handlers: %w", err)
    }
    
    if err := container.initMiddleware(); err != nil {
        return nil, fmt.Errorf("failed to initialize middleware: %w", err)
    }
    
    log.Println("DI Container initialized successfully")
    return container, nil
}

// GetConfig returns the application configuration
func (c *Container) GetConfig() *Config {
    return c.config
}

// GetDB returns the database connection
func (c *Container) GetDB() *gorm.DB {
    return c.db
}

// GetAuthHandler returns the auth handler
func (c *Container) GetAuthHandler() *handler.AuthHandler {
    return c.authHandler
}

// Shutdown gracefully shuts down the container
func (c *Container) Shutdown(ctx context.Context) error {
    log.Println("Shutting down DI Container...")
    
    // Close database connection
    if c.db != nil {
        if sqlDB, err := c.db.DB(); err == nil {
            if err := sqlDB.Close(); err != nil {
                log.Printf("Error closing database: %v", err)
            }
        }
    }
    
    log.Println("DI Container shutdown complete")
    return nil
}
```

### Phase 2: プロバイダー実装 📅 Day 2-3

#### 2.1 データベースプロバイダー
**ファイル**: `internal/di/providers/database.go`
```go
package providers

import (
    "fmt"
    "time"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    
    "github.com/your-org/financetracker/internal/di"
    "github.com/your-org/financetracker/internal/infrastructure/auth0"
)

// initInfrastructure initializes database and external services
func (c *Container) initInfrastructure() error {
    // Initialize database
    db, err := c.setupDatabase()
    if err != nil {
        return fmt.Errorf("database setup failed: %w", err)
    }
    c.db = db
    
    // Initialize Auth0 client
    auth0Client, err := c.setupAuth0Client()
    if err != nil {
        return fmt.Errorf("auth0 setup failed: %w", err)
    }
    c.auth0Client = auth0Client
    
    return nil
}

func (c *Container) setupDatabase() (*gorm.DB, error) {
    config := c.config.Database
    
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        config.Host,
        config.Port,
        config.User,
        config.Password,
        config.DBName,
        config.SSLMode,
    )
    
    // Configure GORM logger based on environment
    var gormLogger logger.Interface
    if c.config.Server.Environment == "development" {
        gormLogger = logger.Default.LogMode(logger.Info)
    } else {
        gormLogger = logger.Default.LogMode(logger.Error)
    }
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: gormLogger,
    })
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    sqlDB.SetMaxOpenConns(config.MaxOpenConns)
    sqlDB.SetMaxIdleConns(config.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db, nil
}

func (c *Container) setupAuth0Client() (*auth0.Client, error) {
    return auth0.NewClient(
        c.config.Auth0.Domain,
        c.config.Auth0.ClientID,
        c.config.Auth0.ClientSecret,
    )
}
```

#### 2.2 リポジトリプロバイダー
**ファイル**: `internal/di/providers/repository.go`
```go
package providers

import (
    "github.com/your-org/financetracker/internal/infrastructure/gorm/repository"
)

// initRepositories initializes all repository implementations
func (c *Container) initRepositories() error {
    // User repository (already implemented)
    c.userRepo = repository.NewUserRepository(c.db)
    
    // TODO: Initialize other repositories as they are implemented
    // c.accountRepo = repository.NewAccountRepository(c.db)
    // c.transactionRepo = repository.NewTransactionRepository(c.db)
    // c.categoryRepo = repository.NewCategoryRepository(c.db)
    // c.budgetRepo = repository.NewBudgetRepository(c.db)
    
    return nil
}
```

#### 2.3 サービスプロバイダー（段階的実装）
**ファイル**: `internal/di/providers/service.go`
```go
package providers

import (
    "github.com/your-org/financetracker/internal/application/service"
)

// initServices initializes all application services
func (c *Container) initServices() error {
    // Auth service (already implemented)
    c.authService = service.NewAuthService(c.userRepo, c.auth0Client)
    
    // TODO: Initialize other services as they are implemented
    // c.userService = service.NewUserService(c.userRepo, c.txManager)
    // c.accountService = service.NewAccountService(c.accountRepo, c.userRepo, c.txManager)
    
    return nil
}
```

### Phase 3: 既存コード統合 📅 Day 3-4

#### 3.1 main.go の更新
**ファイル**: `cmd/api/main.go`
```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/gin-gonic/gin"
    
    "github.com/your-org/financetracker/internal/di"
    "github.com/your-org/financetracker/internal/interface/router"
)

func main() {
    // Initialize DI container
    container, err := di.NewContainer()
    if err != nil {
        log.Fatal("Failed to initialize DI container:", err)
    }
    
    // Set Gin mode based on environment
    if container.GetConfig().Server.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    // Setup router
    r := router.SetupRouter(container)
    
    // Create HTTP server
    srv := &http.Server{
        Addr:         ":" + container.GetConfig().Server.Port,
        Handler:      r,
        ReadTimeout:  container.GetConfig().Server.ReadTimeout,
        WriteTimeout: container.GetConfig().Server.WriteTimeout,
    }
    
    // Start server in a goroutine
    go func() {
        log.Printf("Server starting on port %s", container.GetConfig().Server.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()
    
    // Wait for interrupt signal for graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    // Create shutdown context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Shutdown server
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    // Shutdown DI container
    if err := container.Shutdown(ctx); err != nil {
        log.Printf("Container shutdown error: %v", err)
    }
    
    log.Println("Server exited")
}
```

#### 3.2 ルーター設定の更新
**ファイル**: `internal/interface/router/router.go` (更新)
```go
package router

import (
    "github.com/gin-gonic/gin"
    
    "github.com/your-org/financetracker/internal/di"
    "github.com/your-org/financetracker/internal/interface/handler"
)

// SetupRouter configures and returns the main router
func SetupRouter(container *di.Container) *gin.Engine {
    r := gin.New()
    
    // Add global middleware
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    
    // Add custom middleware (when implemented)
    // r.Use(container.GetCORSMiddleware().CORS())
    // r.Use(container.GetLoggerMiddleware().Logger())
    // r.Use(container.GetErrorMiddleware().ErrorHandler())
    
    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })
    
    // Setup auth routes
    setupAuthRoutes(r, container)
    
    // Setup API v1 routes
    setupV1Routes(r, container)
    
    return r
}

func setupAuthRoutes(r *gin.Engine, container *di.Container) {
    authHandler := container.GetAuthHandler()
    
    auth := r.Group("/auth")
    {
        auth.GET("/login", authHandler.Login)
        auth.GET("/callback", authHandler.Callback)
        auth.POST("/logout", authHandler.Logout)
        auth.GET("/user", authHandler.GetUser)
        auth.GET("/check", authHandler.CheckAuth)
        auth.POST("/token", authHandler.SetToken)
        auth.DELETE("/token", authHandler.DeleteToken)
    }
}

func setupV1Routes(r *gin.Engine, container *di.Container) {
    v1 := r.Group("/api/v1")
    
    // Add authentication middleware when available
    // v1.Use(container.GetAuthMiddleware().RequireAuth())
    
    {
        v1.GET("/", handler.NotImplemented) // API info endpoint
        
        // User routes (to be implemented)
        users := v1.Group("/users")
        {
            users.GET("/me", handler.NotImplemented)
            users.PUT("/me", handler.NotImplemented)
        }
        
        // Account routes (to be implemented)
        accounts := v1.Group("/accounts")
        {
            accounts.GET("", handler.NotImplemented)
            accounts.POST("", handler.NotImplemented)
            accounts.GET("/:id", handler.NotImplemented)
            accounts.PUT("/:id", handler.NotImplemented)
            accounts.DELETE("/:id", handler.NotImplemented)
            accounts.POST("/:id/movements", handler.NotImplemented)
        }
        
        // Transaction routes (to be implemented)
        transactions := v1.Group("/transactions")
        {
            transactions.GET("", handler.NotImplemented)
            transactions.POST("", handler.NotImplemented)
            transactions.GET("/:id", handler.NotImplemented)
            transactions.PUT("/:id", handler.NotImplemented)
            transactions.DELETE("/:id", handler.NotImplemented)
            transactions.GET("/summary/monthly", handler.NotImplemented)
        }
        
        // Other routes...
    }
}
```

### Phase 4: テスト環境構築 📅 Day 4-5

#### 4.1 テスト用DIコンテナ
**ファイル**: `internal/di/test_container.go`
```go
package di

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    
    "github.com/your-org/financetracker/internal/infrastructure/gorm/model"
)

// NewTestContainer creates a container for testing with in-memory database
func NewTestContainer() (*Container, error) {
    config := &Config{
        Server: ServerConfig{
            Port:        "8080",
            Environment: "test",
        },
        Auth0: Auth0Config{
            Domain:       "test.auth0.com",
            ClientID:     "test-client-id",
            ClientSecret: "test-client-secret",
            CallbackURL:  "http://localhost:3000/auth/callback",
        },
        JWT: JWTConfig{
            Secret:      "test-secret-key",
            ExpiryHours: 1 * time.Hour,
        },
        Log: LogConfig{
            Level:  "error",
            Format: "json",
        },
    }
    
    container := &Container{
        config: config,
    }
    
    // Setup in-memory SQLite database
    if err := container.setupTestDatabase(); err != nil {
        return nil, err
    }
    
    // Initialize test repositories
    if err := container.initRepositories(); err != nil {
        return nil, err
    }
    
    // Initialize test services with mock Auth0
    if err := container.initTestServices(); err != nil {
        return nil, err
    }
    
    return container, nil
}

func (c *Container) setupTestDatabase() error {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        return err
    }
    
    // Auto-migrate all models
    err = db.AutoMigrate(
        &model.User{},
        &model.Account{},
        &model.Transaction{},
        &model.Category{},
        &model.CategoryMaster{},
        &model.Budget{},
        &model.BudgetSuggestion{},
        &model.AssetSnapshot{},
        &model.AssetForecast{},
        &model.NotificationSettings{},
    )
    if err != nil {
        return err
    }
    
    c.db = db
    return nil
}

func (c *Container) initTestServices() error {
    // Mock Auth0 client for testing
    c.auth0Client = &mockAuth0Client{}
    
    // Initialize services with mock dependencies
    return c.initServices()
}

// mockAuth0Client for testing
type mockAuth0Client struct{}

func (m *mockAuth0Client) GetUserInfo(token string) (*auth0.UserInfo, error) {
    return &auth0.UserInfo{
        Sub:   "test-user-id",
        Email: "test@example.com",
        Name:  "Test User",
    }, nil
}

func (m *mockAuth0Client) ValidateToken(token string) (*auth0.Claims, error) {
    return &auth0.Claims{
        Sub: "test-user-id",
    }, nil
}
```

#### 4.2 統合テスト例
**ファイル**: `internal/di/container_test.go`
```go
package di_test

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/your-org/financetracker/internal/di"
)

func TestNewContainer(t *testing.T) {
    // Set test environment variables
    t.Setenv("AUTH0_DOMAIN", "test.auth0.com")
    t.Setenv("AUTH0_CLIENT_ID", "test-client-id")
    t.Setenv("AUTH0_CLIENT_SECRET", "test-client-secret")
    
    container, err := di.NewContainer()
    require.NoError(t, err)
    require.NotNil(t, container)
    
    // Test configuration loading
    config := container.GetConfig()
    assert.Equal(t, "test.auth0.com", config.Auth0.Domain)
    assert.Equal(t, "test-client-id", config.Auth0.ClientID)
    
    // Test database connection
    db := container.GetDB()
    assert.NotNil(t, db)
    
    // Test handler initialization
    authHandler := container.GetAuthHandler()
    assert.NotNil(t, authHandler)
}

func TestNewTestContainer(t *testing.T) {
    container, err := di.NewTestContainer()
    require.NoError(t, err)
    require.NotNil(t, container)
    
    // Test that it uses test configuration
    config := container.GetConfig()
    assert.Equal(t, "test", config.Server.Environment)
    
    // Test that database is connected
    db := container.GetDB()
    assert.NotNil(t, db)
    
    // Test that we can perform database operations
    sqlDB, err := db.DB()
    require.NoError(t, err)
    err = sqlDB.Ping()
    assert.NoError(t, err)
}
```

## 導入効果の測定

### Before（DIコンテナなし）
```go
// main.go - 100+ lines of dependency initialization
func main() {
    config := loadConfig()
    db := setupDatabase(config)
    auth0Client := setupAuth0(config)
    userRepo := repository.NewUserRepository(db)
    authService := service.NewAuthService(userRepo, auth0Client)
    authHandler := handler.NewAuthHandler(authService)
    // ... 多数の初期化コード
    router := setupRouter(authHandler, /* many handlers */)
    router.Run()
}
```

### After（DIコンテナあり）
```go
// main.go - 20 lines で完結
func main() {
    container, err := di.NewContainer()
    if err != nil {
        log.Fatal(err)
    }
    router := router.SetupRouter(container)
    router.Run()
}
```

## メリット実現

### 🎯 開発効率向上
1. **新機能追加時の工数削減**: 50-70%削減
2. **設定変更の時間短縮**: 環境変数の一元管理
3. **依存関係の可視化**: 明確な依存関係グラフ

### 🧪 テスタビリティ向上
1. **モックテストの簡素化**: テスト用コンテナで簡単にモック注入
2. **独立したテスト環境**: in-memoryデータベース使用
3. **テスト実行時間短縮**: 高速なテストセットアップ

### 🏗️ アーキテクチャ品質向上
1. **Clean Architectureの徹底**: 依存関係逆転の完全実現
2. **単一責任原則の遵守**: 各コンポーネントの役割明確化
3. **拡張性の確保**: 新機能の追加が容易

### 📈 保守性向上
1. **設定管理の統一**: 環境変数の体系的管理
2. **エラー処理の改善**: 初期化エラーの適切なハンドリング
3. **ログ品質の向上**: 構造化された依存関係ログ

この実装により、FinanceTrackerバックエンドは次のレベルの開発効率と品質を実現できます。