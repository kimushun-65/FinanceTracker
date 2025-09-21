# Goプロジェクト基本セットアップガイド

## 概要
このドキュメントは、Goで新規プロジェクトを開始する際の標準的なセットアップ手順をまとめたものです。
Docker環境でのクリーンアーキテクチャを採用したWebアプリケーション開発を想定しています。

## 目次
1. [プロジェクト初期化](#1-プロジェクト初期化)
2. [Docker環境構築](#2-docker環境構築)
3. [プロジェクト構造設計](#3-プロジェクト構造設計)
4. [基本パッケージ実装](#4-基本パッケージ実装)
5. [データベースセットアップ](#5-データベースセットアップ)
6. [HTTPサーバー構築](#6-httpサーバー構築)
7. [開発支援ツール設定](#7-開発支援ツール設定)
8. [テスト環境構築](#8-テスト環境構築)

---

## 1. プロジェクト初期化

### 1.1 プロジェクトディレクトリ作成
```bash
mkdir myproject
cd myproject
mkdir backend frontend docs
cd backend
```

### 1.2 Go Modules初期化
```bash
go mod init github.com/username/myproject
```

### 1.3 .gitignore作成
```gitignore
# Binaries
*.exe
*.dll
*.so
*.dylib
/bin/
/dist/

# Test binary
*.test

# Output of go coverage
*.out

# Dependency directories
/vendor/

# IDE
.idea/
.vscode/
*.swp
*.swo

# Environment
.env
.env.local

# OS
.DS_Store
Thumbs.db

# Temporary
/tmp/
*.log
```

### 1.4 環境変数テンプレート作成
```bash
touch .env.example
```

```env
# Application
APP_ENV=development
APP_PORT=8080
APP_NAME=myproject

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=myproject_db
DB_SSLMODE=disable

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000

# Auth (if needed)
AUTH_SECRET_KEY=your-secret-key
```

---

## 2. Docker環境構築

### 2.1 Dockerfile作成
```dockerfile
FROM golang:1.23-alpine

# 必要なツールのインストール
RUN apk update && apk add --no-cache \
    git \
    curl \
    postgresql-client \
    make \
    && rm -rf /var/cache/apk/*

# 開発ツールのインストール
RUN go install github.com/air-verse/air@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# タイムゾーン設定
ENV TZ=Asia/Tokyo
RUN apk add --no-cache tzdata

WORKDIR /app

# Go modulesの依存関係を先にコピー（キャッシュ効率化）
COPY go.mod go.sum ./
RUN go mod download

# ソースコード全体をコピー
COPY . .

# 開発用エントリーポイント
CMD ["air", "-c", ".air.toml"]
```

### 2.2 docker-compose.yml作成
```yaml
services:
  # Backend API
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: myproject_backend
    env_file:
      - ./backend/.env
    ports:
      - "8080:8080"
    volumes:
      - ./backend:/app
      - go_modules:/go/pkg/mod
    depends_on:
      db:
        condition: service_healthy
    networks:
      - myproject_network

  # PostgreSQL Database
  db:
    image: postgres:15-alpine
    container_name: myproject_postgres
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=myproject_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - myproject_network

  # Database GUI (optional)
  pgadmin:
    image: dpage/pgadmin4:latest
    container_name: myproject_pgadmin
    ports:
      - "5050:80"
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@example.com
      PGADMIN_DEFAULT_PASSWORD: admin
    depends_on:
      - db
    networks:
      - myproject_network

volumes:
  postgres_data:
  go_modules:

networks:
  myproject_network:
    driver: bridge
```

### 2.3 Air設定ファイル作成
```toml
# .air.toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/api"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = true

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
```

---

## 3. プロジェクト構造設計

### 3.1 クリーンアーキテクチャベースの構造
```bash
mkdir -p cmd/api
mkdir -p internal/{domain,application,infrastructure,interface}
mkdir -p pkg/{config,logger,errors}
mkdir -p scripts
mkdir -p docs
```

### 3.2 ディレクトリ構造
```
backend/
├── cmd/
│   └── api/          # アプリケーションエントリポイント
│       └── main.go
├── internal/         # 内部パッケージ（外部からimport不可）
│   ├── domain/       # ビジネスロジック層
│   │   ├── entity/   # エンティティ
│   │   ├── value/    # 値オブジェクト
│   │   └── repository/ # リポジトリインターフェース
│   ├── application/  # アプリケーション層
│   │   ├── usecase/  # ユースケース
│   │   └── dto/      # データ転送オブジェクト
│   ├── infrastructure/ # インフラストラクチャ層
│   │   ├── persistence/ # データベース実装
│   │   └── external/   # 外部サービス
│   └── interface/    # インターフェース層
│       ├── handler/  # HTTPハンドラー
│       ├── middleware/ # ミドルウェア
│       └── router/   # ルーティング
├── pkg/             # 外部パッケージ（他プロジェクトから利用可能）
│   ├── config/      # 設定管理
│   ├── logger/      # ロギング
│   └── errors/      # エラー定義
├── scripts/         # ユーティリティスクリプト
├── docs/           # ドキュメント
├── Makefile        # ビルドタスク
├── go.mod
└── go.sum
```

---

## 4. 基本パッケージ実装

### 4.1 設定管理 (pkg/config/config.go)
```go
package config

import (
    "os"
    "strconv"
)

type Config struct {
    AppEnv  string
    AppPort string
    AppName string
    
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
    
    LogLevel string
    
    CORSAllowedOrigins []string
}

func Load() *Config {
    return &Config{
        AppEnv:  getEnv("APP_ENV", "development"),
        AppPort: getEnv("APP_PORT", "8080"),
        AppName: getEnv("APP_NAME", "myproject"),
        
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "myproject_db"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        
        LogLevel: getEnv("LOG_LEVEL", "debug"),
        
        CORSAllowedOrigins: []string{getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")},
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### 4.2 ロガー設定 (pkg/logger/logger.go)
```go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func New() *zap.Logger {
    config := zap.NewDevelopmentConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    
    logger, _ := config.Build()
    return logger
}
```

### 4.3 エラー定義 (pkg/errors/errors.go)
```go
package errors

import (
    "fmt"
    "net/http"
)

type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Status  int    `json:"-"`
}

func (e *AppError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common errors
var (
    ErrNotFound     = &AppError{"NOT_FOUND", "Resource not found", http.StatusNotFound}
    ErrBadRequest   = &AppError{"BAD_REQUEST", "Invalid request", http.StatusBadRequest}
    ErrInternal     = &AppError{"INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError}
    ErrUnauthorized = &AppError{"UNAUTHORIZED", "Unauthorized access", http.StatusUnauthorized}
)
```

---

## 5. データベースセットアップ

### 5.1 マイグレーションツールの選択
以下のいずれかを選択：
- **Atlas**: スキーマ管理に特化
- **golang-migrate**: シンプルなSQLマイグレーション
- **GORM AutoMigrate**: 開発初期の簡易的な方法

### 5.2 Atlas使用例

#### atlas.hcl作成
```hcl
env "dev" {
  url = "postgres://postgres:postgres@localhost:5432/myproject_db?sslmode=disable"
  dev = "docker://postgres/15/dev"
}
```

#### Makefile追加
```makefile
.PHONY: migrate-diff migrate-apply

migrate-diff:
	atlas migrate diff $(name) \
		--env dev \
		--to file://schema.hcl

migrate-apply:
	atlas migrate apply \
		--env dev \
		--url $${DATABASE_URL}
```

### 5.3 初期テーブル設計例
```sql
-- users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create update trigger
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
```

---

## 6. HTTPサーバー構築

### 6.1 メインエントリポイント (cmd/api/main.go)
```go
package main

import (
    "log"
    
    "github.com/username/myproject/internal/interface/router"
    "github.com/username/myproject/pkg/config"
    "github.com/username/myproject/pkg/logger"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize logger
    logger := logger.New()
    defer logger.Sync()
    
    // Setup router
    r := router.New(cfg, logger)
    
    // Start server
    logger.Info("Starting server on :" + cfg.AppPort)
    if err := r.Run(":" + cfg.AppPort); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

### 6.2 ルーター設定 (internal/interface/router/router.go)
```go
package router

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    
    "github.com/username/myproject/internal/interface/middleware"
    "github.com/username/myproject/pkg/config"
)

func New(cfg *config.Config, logger *zap.Logger) *gin.Engine {
    if cfg.AppEnv == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    r := gin.New()
    
    // Global middleware
    r.Use(middleware.Logger(logger))
    r.Use(middleware.ErrorHandler())
    r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })
    
    // API routes
    api := r.Group("/api")
    {
        // Add your routes here
    }
    
    return r
}
```

---

## 7. 開発支援ツール設定

### 7.1 Makefile作成
```makefile
.PHONY: help up down build test lint

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

up: ## Start docker containers
	docker-compose up -d

down: ## Stop docker containers
	docker-compose down

build: ## Build application
	go build -o bin/api cmd/api/main.go

test: ## Run tests
	go test -v ./...

lint: ## Run linter
	golangci-lint run ./...

migrate: ## Run database migrations
	@echo "Running migrations..."
	# Add your migration command here

seed: ## Seed database
	go run cmd/seed/main.go
```

### 7.2 GitHub Actions設定 (.github/workflows/ci.yml)
```yaml
name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.23'
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v ./...
    
    - name: Run linter
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
```

---

## 8. テスト環境構築

### 8.1 テスト用設定
```go
// internal/test/setup.go
package test

import (
    "testing"
    "database/sql"
)

func SetupTestDB(t *testing.T) *sql.DB {
    // Test database setup
    t.Helper()
    // Implementation here
    return nil
}
```

### 8.2 テスト例
```go
// internal/domain/entity/user_test.go
package entity_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
    tests := []struct {
        name    string
        user    *entity.User
        wantErr bool
    }{
        // Test cases here
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

---

## チェックリスト

### 初期セットアップ完了確認
- [ ] Go Modules初期化
- [ ] .gitignore作成
- [ ] 環境変数設定
- [ ] Docker環境構築
- [ ] プロジェクト構造作成
- [ ] 基本パッケージ実装
- [ ] データベース接続確認
- [ ] HTTPサーバー起動確認
- [ ] ヘルスチェックエンドポイント動作確認
- [ ] 開発ツール設定
- [ ] CI/CD設定
- [ ] README.md作成

### 推奨される次のステップ
1. ドメインモデルの設計・実装
2. ユースケースの定義
3. リポジトリインターフェースの定義
4. インフラストラクチャ層の実装
5. HTTPハンドラーの実装
6. 認証・認可の実装
7. ロギング・監視の強化
8. パフォーマンスチューニング

---

## 参考リソース
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture in Go](https://github.com/bxcodec/go-clean-arch)

## トラブルシューティング

### よくある問題と解決方法

#### 1. Dockerコンテナが起動しない
```bash
# ログを確認
docker-compose logs -f backend

# コンテナを再ビルド
docker-compose build --no-cache backend
docker-compose up -d
```

#### 2. データベース接続エラー
```bash
# PostgreSQLの状態確認
docker-compose ps db
docker-compose logs db

# 接続情報の確認
docker exec -it myproject_postgres psql -U postgres -d myproject_db
```

#### 3. ホットリロードが効かない
```bash
# Air の設定確認
docker exec -it myproject_backend air -v

# ファイルシステムの権限確認
ls -la ./backend
```