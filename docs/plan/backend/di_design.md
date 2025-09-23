# DI（Dependency Injection）層設計書

## 概要
FinanceTrackerバックエンドにDI層を導入することで、Clean Architectureの依存関係管理を効率化し、テスタビリティと保守性を向上させます。

## DIコンテナの必要性

### 🚨 現在の問題点
1. **依存関係の散在**: main.goで手動で依存関係を構築
2. **循環依存のリスク**: 複雑な依存関係での構築順序ミス
3. **テストの困難さ**: モックオブジェクトの注入が困難
4. **コードの重複**: 似たような依存関係構築コードの重複
5. **設定管理の分散**: 環境変数の読み込みが各所で発生

### ✅ DIコンテナによる解決
1. **中央集権的管理**: 全ての依存関係を1箇所で管理
2. **自動的な依存関係解決**: 正しい順序での初期化
3. **テスタビリティ**: モックオブジェクトの簡単な注入
4. **設定の統一化**: 環境変数・設定の一元管理
5. **Clean Architectureの強化**: 依存関係逆転の適切な実装

## ディレクトリ構成

```
backend/internal/di/
├── container.go              # メインDIコンテナ
├── config.go                 # 設定管理
├── providers/                # プロバイダー群
│   ├── database.go           # データベース関連プロバイダー
│   ├── repository.go         # リポジトリプロバイダー
│   ├── service.go            # サービスプロバイダー
│   ├── handler.go            # ハンドラープロバイダー
│   └── middleware.go         # ミドルウェアプロバイダー
└── interfaces.go             # DIコンテナインターフェース
```

## DIコンテナ設計

### コンテナ構造

```go
package di

import (
    "gorm.io/gorm"
    
    // Domain repositories
    userRepo "github.com/your-org/financetracker/internal/domain/user/repository"
    accountRepo "github.com/your-org/financetracker/internal/domain/account/repository"
    transactionRepo "github.com/your-org/financetracker/internal/domain/transaction/repository"
    categoryRepo "github.com/your-org/financetracker/internal/domain/category/repository"
    budgetRepo "github.com/your-org/financetracker/internal/domain/budget/repository"
    assetRepo "github.com/your-org/financetracker/internal/domain/asset/repository"
    notificationRepo "github.com/your-org/financetracker/internal/domain/notification/repository"
    
    // Application services
    "github.com/your-org/financetracker/internal/application/service"
    "github.com/your-org/financetracker/internal/application/transaction"
    
    // Infrastructure
    "github.com/your-org/financetracker/internal/infrastructure/auth0"
    
    // Interface
    "github.com/your-org/financetracker/internal/interface/handler"
    "github.com/your-org/financetracker/internal/interface/middleware"
)

// Container holds all dependencies
type Container struct {
    // Configuration
    Config *Config
    
    // Database
    DB *gorm.DB
    
    // External Services
    Auth0Client *auth0.Client
    
    // Repositories - using domain interfaces
    UserRepo         userRepo.UserRepository
    AccountRepo      accountRepo.AccountRepository
    TransactionRepo  transactionRepo.TransactionRepository
    CategoryRepo     categoryRepo.CategoryRepository
    CategoryMasterRepo categoryRepo.CategoryMasterRepository
    BudgetRepo       budgetRepo.BudgetRepository
    BudgetSuggestionRepo budgetRepo.BudgetSuggestionRepository
    AssetSnapshotRepo assetRepo.AssetSnapshotRepository
    AssetForecastRepo assetRepo.AssetForecastRepository
    NotificationRepo notificationRepo.NotificationSettingsRepository
    
    // Transaction Manager
    TxManager transaction.Manager
    
    // Application Services
    AuthService         service.AuthService
    UserService         service.UserService
    AccountService      service.AccountService
    TransactionService  service.TransactionService
    CategoryService     service.CategoryService
    BudgetService       service.BudgetService
    AssetService        service.AssetService
    NotificationService service.NotificationService
    ReportService       service.ReportService
    
    // HTTP Handlers
    AuthHandler         *handler.AuthHandler
    UserHandler         *handler.UserHandler
    AccountHandler      *handler.AccountHandler
    TransactionHandler  *handler.TransactionHandler
    CategoryHandler     *handler.CategoryHandler
    BudgetHandler       *handler.BudgetHandler
    AssetHandler        *handler.AssetHandler
    ReportHandler       *handler.ReportHandler
    NotificationHandler *handler.NotificationHandler
    
    // Middleware
    AuthMiddleware     *middleware.AuthMiddleware
    CORSMiddleware     *middleware.CORSMiddleware
    LoggerMiddleware   *middleware.LoggerMiddleware
    ErrorMiddleware    *middleware.ErrorMiddleware
}

// NewContainer creates and initializes all dependencies
func NewContainer() (*Container, error) {
    c := &Container{}
    
    // Load configuration
    if err := c.loadConfig(); err != nil {
        return nil, err
    }
    
    // Initialize infrastructure
    if err := c.initInfrastructure(); err != nil {
        return nil, err
    }
    
    // Initialize repositories
    c.initRepositories()
    
    // Initialize transaction manager
    c.initTransactionManager()
    
    // Initialize services
    c.initServices()
    
    // Initialize handlers
    c.initHandlers()
    
    // Initialize middleware
    c.initMiddleware()
    
    return c, nil
}
```

### 設定管理

```go
package di

import (
    "os"
    "strconv"
    "time"
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
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type Auth0Config struct {
    Domain       string
    ClientID     string
    ClientSecret string
    CallbackURL  string
}

type JWTConfig struct {
    Secret      string
    ExpiryHours time.Duration
}

type LogConfig struct {
    Level  string
    Format string // json or console
}

func (c *Container) loadConfig() error {
    c.Config = &Config{
        Server: ServerConfig{
            Port:         getEnv("SERVER_PORT", "8080"),
            ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
            WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getIntEnv("DB_PORT", 5432),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            DBName:   getEnv("DB_NAME", "financetracker"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        Auth0: Auth0Config{
            Domain:       getEnv("AUTH0_DOMAIN", ""),
            ClientID:     getEnv("AUTH0_CLIENT_ID", ""),
            ClientSecret: getEnv("AUTH0_CLIENT_SECRET", ""),
            CallbackURL:  getEnv("AUTH0_CALLBACK_URL", ""),
        },
        JWT: JWTConfig{
            Secret:      getEnv("JWT_SECRET", "default-secret"),
            ExpiryHours: getDurationEnv("JWT_EXPIRY_HOURS", 24*time.Hour),
        },
        Log: LogConfig{
            Level:  getEnv("LOG_LEVEL", "info"),
            Format: getEnv("LOG_FORMAT", "json"),
        },
    }
    return nil
}

// Helper functions for environment variables
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

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}
```

## プロバイダー分割設計

### 1. データベースプロバイダー

```go
// providers/database.go
package providers

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/your-org/financetracker/internal/di"
)

func (c *Container) initInfrastructure() error {
    // Database connection
    db, err := c.setupDatabase()
    if err != nil {
        return err
    }
    c.DB = db
    
    // Auth0 client
    auth0Client, err := c.setupAuth0()
    if err != nil {
        return err
    }
    c.Auth0Client = auth0Client
    
    return nil
}

func (c *Container) setupDatabase() (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Config.Database.Host,
        c.Config.Database.Port,
        c.Config.Database.User,
        c.Config.Database.Password,
        c.Config.Database.DBName,
        c.Config.Database.SSLMode,
    )
    
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func (c *Container) setupAuth0() (*auth0.Client, error) {
    return auth0.NewClient(
        c.Config.Auth0.Domain,
        c.Config.Auth0.ClientID,
        c.Config.Auth0.ClientSecret,
    )
}
```

### 2. リポジトリプロバイダー

```go
// providers/repository.go
package providers

import (
    "github.com/your-org/financetracker/internal/infrastructure/gorm/repository"
)

func (c *Container) initRepositories() {
    c.UserRepo = repository.NewUserRepository(c.DB)
    c.AccountRepo = repository.NewAccountRepository(c.DB)
    c.TransactionRepo = repository.NewTransactionRepository(c.DB)
    c.CategoryRepo = repository.NewCategoryRepository(c.DB)
    c.CategoryMasterRepo = repository.NewCategoryMasterRepository(c.DB)
    c.BudgetRepo = repository.NewBudgetRepository(c.DB)
    c.BudgetSuggestionRepo = repository.NewBudgetSuggestionRepository(c.DB)
    c.AssetSnapshotRepo = repository.NewAssetSnapshotRepository(c.DB)
    c.AssetForecastRepo = repository.NewAssetForecastRepository(c.DB)
    c.NotificationRepo = repository.NewNotificationSettingsRepository(c.DB)
}
```

### 3. サービスプロバイダー

```go
// providers/service.go
package providers

import (
    "github.com/your-org/financetracker/internal/application/service"
    "github.com/your-org/financetracker/internal/application/transaction"
)

func (c *Container) initTransactionManager() {
    c.TxManager = transaction.NewManager(c.DB)
}

func (c *Container) initServices() {
    c.AuthService = service.NewAuthService(
        c.UserRepo,
        c.Auth0Client,
    )
    
    c.UserService = service.NewUserService(
        c.UserRepo,
        c.TxManager,
    )
    
    c.AccountService = service.NewAccountService(
        c.AccountRepo,
        c.UserRepo,
        c.TxManager,
    )
    
    c.TransactionService = service.NewTransactionService(
        c.TransactionRepo,
        c.AccountRepo,
        c.CategoryRepo,
        c.TxManager,
    )
    
    c.CategoryService = service.NewCategoryService(
        c.CategoryRepo,
        c.CategoryMasterRepo,
        c.UserRepo,
        c.TxManager,
    )
    
    c.BudgetService = service.NewBudgetService(
        c.BudgetRepo,
        c.BudgetSuggestionRepo,
        c.TransactionRepo,
        c.CategoryRepo,
        c.TxManager,
    )
    
    c.AssetService = service.NewAssetService(
        c.AssetSnapshotRepo,
        c.AssetForecastRepo,
        c.AccountRepo,
        c.TxManager,
    )
    
    c.NotificationService = service.NewNotificationService(
        c.NotificationRepo,
        c.UserRepo,
        c.TxManager,
    )
    
    c.ReportService = service.NewReportService(
        c.TransactionRepo,
        c.AccountRepo,
        c.BudgetRepo,
        c.AssetSnapshotRepo,
        c.CategoryRepo,
    )
}
```

### 4. ハンドラープロバイダー

```go
// providers/handler.go
package providers

import (
    "github.com/your-org/financetracker/internal/interface/handler"
)

func (c *Container) initHandlers() {
    c.AuthHandler = handler.NewAuthHandler(c.AuthService)
    c.UserHandler = handler.NewUserHandler(c.UserService)
    c.AccountHandler = handler.NewAccountHandler(c.AccountService)
    c.TransactionHandler = handler.NewTransactionHandler(c.TransactionService)
    c.CategoryHandler = handler.NewCategoryHandler(c.CategoryService)
    c.BudgetHandler = handler.NewBudgetHandler(c.BudgetService)
    c.AssetHandler = handler.NewAssetHandler(c.AssetService)
    c.ReportHandler = handler.NewReportHandler(c.ReportService)
    c.NotificationHandler = handler.NewNotificationHandler(c.NotificationService)
}
```

### 5. ミドルウェアプロバイダー

```go
// providers/middleware.go
package providers

import (
    "github.com/your-org/financetracker/internal/interface/middleware"
)

func (c *Container) initMiddleware() {
    c.AuthMiddleware = middleware.NewAuthMiddleware(c.Auth0Client)
    c.CORSMiddleware = middleware.NewCORSMiddleware()
    c.LoggerMiddleware = middleware.NewLoggerMiddleware(c.Config.Log)
    c.ErrorMiddleware = middleware.NewErrorMiddleware()
}
```

## main.go の簡略化

### Before (DIコンテナなし)
```go
func main() {
    // 設定読み込み
    config := loadConfig()
    
    // データベース接続
    db := setupDatabase(config)
    
    // Auth0クライアント
    auth0Client := setupAuth0(config)
    
    // リポジトリ初期化
    userRepo := repository.NewUserRepository(db)
    accountRepo := repository.NewAccountRepository(db)
    // ... 多数のリポジトリ
    
    // サービス初期化
    authService := service.NewAuthService(userRepo, auth0Client)
    userService := service.NewUserService(userRepo, txManager)
    // ... 多数のサービス
    
    // ハンドラー初期化
    authHandler := handler.NewAuthHandler(authService)
    userHandler := handler.NewUserHandler(userService)
    // ... 多数のハンドラー
    
    // ルーター設定
    router := setupRouter(authHandler, userHandler /* ... */)
    
    // サーバー起動
    router.Run(":8080")
}
```

### After (DIコンテナあり)
```go
func main() {
    // DIコンテナ初期化
    container, err := di.NewContainer()
    if err != nil {
        log.Fatal("Failed to initialize container:", err)
    }
    
    // ルーター設定
    router := setupRouter(container)
    
    // サーバー起動
    log.Printf("Server starting on port %s", container.Config.Server.Port)
    router.Run(":" + container.Config.Server.Port)
}

func setupRouter(container *di.Container) *gin.Engine {
    router := gin.New()
    
    // ミドルウェア適用
    router.Use(container.LoggerMiddleware.Logger())
    router.Use(container.CORSMiddleware.CORS())
    router.Use(container.ErrorMiddleware.ErrorHandler())
    
    // ルート定義
    v1 := router.Group("/api/v1")
    v1.Use(container.AuthMiddleware.RequireAuth())
    {
        v1.GET("/users/me", container.UserHandler.GetMe)
        v1.PUT("/users/me", container.UserHandler.UpdateMe)
        
        v1.GET("/accounts", container.AccountHandler.List)
        v1.POST("/accounts", container.AccountHandler.Create)
        // ... その他のルート
    }
    
    return router
}
```

## テスト用DIコンテナ

```go
// di/test_container.go
package di

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/your-org/financetracker/internal/infrastructure/gorm/repository"
)

// NewTestContainer creates a container for testing
func NewTestContainer() (*Container, error) {
    c := &Container{}
    
    // テスト用設定
    c.Config = &Config{
        Database: DatabaseConfig{
            // SQLite in-memory for testing
        },
        JWT: JWTConfig{
            Secret:      "test-secret",
            ExpiryHours: 1 * time.Hour,
        },
    }
    
    // テスト用データベース (SQLite in-memory)
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    c.DB = db
    
    // Auto-migrate for testing
    if err := c.DB.AutoMigrate(/* all models */); err != nil {
        return nil, err
    }
    
    // Mock Auth0 client
    c.Auth0Client = &mockAuth0Client{}
    
    // Initialize with real implementations
    c.initRepositories()
    c.initTransactionManager()
    c.initServices()
    c.initHandlers()
    c.initMiddleware()
    
    return c, nil
}

// MockContainer allows for dependency injection in tests
type MockContainer struct {
    *Container
}

func (m *MockContainer) SetUserRepo(repo userRepo.UserRepository) {
    m.UserRepo = repo
    // Re-initialize dependent services
    m.AuthService = service.NewAuthService(m.UserRepo, m.Auth0Client)
    m.UserService = service.NewUserService(m.UserRepo, m.TxManager)
}
```

## 利点とメリット

### 🎯 開発効率の向上
1. **依存関係の一元管理**: 新しいサービスやリポジトリの追加が簡単
2. **設定の統一化**: 環境変数の管理が一箇所に集約
3. **初期化順序の自動化**: 依存関係の正しい初期化順序が保証

### 🧪 テスタビリティの向上
1. **モックの簡単な注入**: テスト用のコンテナでモックオブジェクトを簡単に注入
2. **独立したテスト環境**: テスト用データベース（SQLite）の使用
3. **部分的なモッキング**: 特定の依存関係のみをモックに置き換え

### 🏗️ アーキテクチャの強化
1. **Clean Architectureの徹底**: 依存関係逆転の適切な実装
2. **単一責任原則**: 各コンポーネントの責務が明確
3. **拡張性**: 新しい機能の追加が容易

### 🔧 保守性の向上
1. **設定変更の容易さ**: 環境変数の追加・変更が簡単
2. **デバッグの効率化**: 依存関係の問題の特定が容易
3. **リファクタリングの安全性**: 依存関係の変更の影響範囲が明確

## 実装スケジュール

### Phase 1: DI基盤構築 (1-2日)
1. `internal/di/` ディレクトリ作成
2. `Container` 構造体定義
3. 設定管理システム実装
4. 基本的なプロバイダー実装

### Phase 2: 既存機能の移行 (2-3日)
1. 認証関連の依存関係移行
2. ユーザー管理の依存関係移行
3. `main.go` の簡略化
4. テストの更新

### Phase 3: 新機能実装 (継続的)
1. 新しいサービス・リポジトリの追加時にDIコンテナを活用
2. テスト用コンテナの拡充
3. モック機能の強化

## まとめ

DI層の導入により、FinanceTrackerバックエンドは以下の改善を実現できます：

1. **Clean Architectureの完全実装**
2. **開発効率とテスタビリティの大幅向上**
3. **保守性と拡張性の強化**
4. **設定管理の統一化**

これにより、今後のビジネスロジック実装（トランザクション管理、口座管理等）がより効率的かつ堅牢に行えるようになります。