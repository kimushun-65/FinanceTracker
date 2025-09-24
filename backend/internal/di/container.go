package di

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// ドメインリポジトリ
	accountRepo "financetracker/internal/domain/account/repository"
	assetRepo "financetracker/internal/domain/asset/repository"
	budgetRepo "financetracker/internal/domain/budget/repository"
	categoryRepo "financetracker/internal/domain/category/repository"
	notificationRepo "financetracker/internal/domain/notification/repository"
	transactionRepo "financetracker/internal/domain/transaction/repository"
	userRepo "financetracker/internal/domain/user/repository"

	// アプリケーションサービス
	appHandler "financetracker/internal/application/handler"
	"financetracker/internal/application/service"

	// インフラストラクチャ
	"financetracker/internal/infrastructure/auth0"
	postgresRepo "financetracker/internal/infrastructure/gorm/repository"

	// インターフェース
	"financetracker/internal/interface/handler"
	loggerPkg "financetracker/pkg/logger"
)

// Container は全ての依存関係を保持する
type Container struct {
	// 設定
	Config *Config

	// データベース
	DB *gorm.DB

	// 外部サービス
	Auth0Client *auth0.Client

	// リポジトリ - ドメインインターフェースを使用
	UserRepo             userRepo.UserRepository
	AccountRepo          accountRepo.AccountRepository
	TransactionRepo      transactionRepo.TransactionRepository
	CategoryRepo         categoryRepo.CategoryRepository
	CategoryMasterRepo   categoryRepo.CategoryMasterRepository
	BudgetRepo           budgetRepo.BudgetRepository
	BudgetSuggestionRepo budgetRepo.BudgetSuggestionRepository
	AssetSnapshotRepo    assetRepo.AssetSnapshotRepository
	AssetForecastRepo    assetRepo.AssetForecastRepository
	NotificationRepo     notificationRepo.NotificationSettingsRepository

	// アプリケーションサービス
	AuthService        *service.AuthService
	UserService        *service.UserService
	AccountService     *service.AccountService
	TransactionService *service.TransactionService

	// HTTPハンドラー
	AuthHandler        *handler.AuthHandler
	UserHandler        *handler.UserHandler
	AccountHandler     *handler.AccountHandler
	TransactionHandler *handler.TransactionHandler

	// 共有ロジックハンドラー
	UserLogicHandler *appHandler.UserLogicHandler

	// ロガー
	Logger *loggerPkg.Logger
}

// NewContainer は全ての依存関係を初期化して作成する
func NewContainer() (*Container, error) {
	c := &Container{}

	// 設定を読み込む
	if err := c.loadConfig(); err != nil {
		return nil, err
	}

	// インフラストラクチャを初期化
	if err := c.initInfrastructure(); err != nil {
		return nil, err
	}

	// リポジトリを初期化
	c.initRepositories()

	// サービスを初期化
	c.initServices()

	// ハンドラーを初期化
	c.initHandlers()

	return c, nil
}

// Close は全てのリソースをクローズする
func (c *Container) Close() error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// initInfrastructure はインフラストラクチャの初期化を行う
func (c *Container) initInfrastructure() error {
	// データベース接続を設定
	db, err := c.setupDatabase()
	if err != nil {
		return fmt.Errorf("データベース接続の設定に失敗しました: %w", err)
	}
	c.DB = db

	// Auth0クライアントを設定
	auth0Client, err := c.setupAuth0()
	if err != nil {
		return fmt.Errorf("Auth0クライアントの設定に失敗しました: %w", err)
	}
	c.Auth0Client = auth0Client

	return nil
}

// setupDatabase はデータベース接続を設定する
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

	// ログレベルを設定
	logLevel := logger.Silent
	if c.Config.Log.Level == "debug" {
		logLevel = logger.Info
	}

	// GORM設定
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// データベースに接続
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	// 接続プールの設定
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 最大アイドル接続数
	sqlDB.SetMaxIdleConns(10)
	// 最大接続数
	sqlDB.SetMaxOpenConns(100)
	// 接続の最大生存時間
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// setupAuth0 はAuth0クライアントを設定する
func (c *Container) setupAuth0() (*auth0.Client, error) {
	if c.Config.Auth0.Domain == "" {
		return nil, fmt.Errorf("AUTH0_DOMAINが設定されていません")
	}

	return auth0.NewClient(
		c.Config.Auth0.Domain,
		c.Config.Auth0.ClientID,
		c.Config.Auth0.ClientSecret,
	), nil
}

// initRepositories は全てのリポジトリを初期化する
func (c *Container) initRepositories() {
	// ユーザーリポジトリ
	c.UserRepo = postgresRepo.NewUserRepository(c.DB)

	// アカウントリポジトリ
	c.AccountRepo = postgresRepo.NewAccountRepository(c.DB)

	// 取引リポジトリ
	c.TransactionRepo = postgresRepo.NewTransactionRepository(c.DB)

	// カテゴリリポジトリ
	c.CategoryRepo = postgresRepo.NewCategoryRepository(c.DB)
	c.CategoryMasterRepo = postgresRepo.NewCategoryMasterRepository(c.DB)

	// 予算リポジトリ
	c.BudgetRepo = postgresRepo.NewBudgetRepository(c.DB)
}

// initServices は全てのアプリケーションサービスを初期化する
func (c *Container) initServices() {
	// ロガーを初期化
	c.Logger = loggerPkg.NewLogger(c.Config.Log.Level)

	// 認証サービス
	c.AuthService = service.NewAuthService(
		c.UserRepo,
	)

	// ユーザーサービス
	c.UserService = service.NewUserService(
		c.UserRepo,
		c.Logger,
	)

	// 口座サービス
	c.AccountService = service.NewAccountService(
		c.AccountRepo,
		c.Logger,
	)

	// 取引サービス
	c.TransactionService = service.NewTransactionService(
		c.TransactionRepo,
		c.Logger,
	)
}

// initHandlers は全てのHTTPハンドラーを初期化する
func (c *Container) initHandlers() {
	// 認証ハンドラー
	c.AuthHandler = handler.NewAuthHandler(
		c.AuthService,
		c.Auth0Client,
		&auth0.AuthMiddleware{},
		c.Config.Auth0.Domain,
		c.Config.Auth0.Audience,
	)

	// 共有ロジックハンドラー
	c.UserLogicHandler = appHandler.NewUserLogicHandler(
		c.UserService,
		c.Logger,
	)

	// ユーザーハンドラー
	c.UserHandler = handler.NewUserHandler(
		c.UserService,
		c.Logger,
	)

	// 口座ハンドラー
	c.AccountHandler = handler.NewAccountHandler(
		c.AccountService,
		c.Logger,
	)

	// トランザクションハンドラー
	c.TransactionHandler = handler.NewTransactionHandler(
		c.TransactionService,
		c.Logger,
	)
}
