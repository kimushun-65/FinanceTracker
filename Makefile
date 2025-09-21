.PHONY: help up down restart build logs ps clean migrate-check migrate-diff migrate-auto migrate-apply migrate-status backend-shell db-shell migrate-shell seed seed-prod

# デフォルトターゲット
help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Docker Compose コマンド
up: ## すべてのサービスを起動
	docker-compose up

down: ## すべてのサービスを停止
	docker-compose down

restart: ## すべてのサービスを再起動
	docker-compose restart

build: ## すべてのサービスをビルド
	docker-compose build --no-cache

logs: ## ログを表示（フォロー）
	docker-compose logs -f

ps: ## サービスの状態を表示
	docker-compose ps

clean: ## ボリュームを含めて削除（データベースも削除）
	docker-compose down -v

# 個別サービス操作
backend-up: ## バックエンドのみ起動
	docker-compose up -d backend

backend-logs: ## バックエンドのログを表示
	docker-compose logs -f backend

backend-restart: ## バックエンドを再起動
	docker-compose restart backend

backend-shell: ## バックエンドコンテナに入る
	docker-compose exec backend sh

db-shell: ## データベースコンテナに入る（PostgreSQL）
	docker-compose exec db psql -U postgres financetracker_db

pgadmin-up: ## pgAdminを起動
	docker-compose up -d pgadmin

# Migration commands（Docker上で実行）
migrate-shell: ## Migrateコンテナに入る
	docker-compose run --rm migrate sh

migrate-check: ## スキーマ差分をチェック
	@echo "=== Checking schema differences ==="
	docker-compose run --rm migrate go run cmd/migrate/main.go check

migrate-diff: ## マイグレーション差分を生成
	@echo "=== Generating migration diff ==="
	docker-compose run --rm migrate go run cmd/migrate/main.go diff

migrate-auto: ## 自動でスキーマ差分をチェックし、必要に応じてマイグレーション生成・適用
	@echo "=== Running auto migration ==="
	docker-compose run --rm migrate go run cmd/migrate/main.go auto

migrate-apply: ## Atlasマイグレーションを適用
	@echo "=== Applying Atlas migrations ==="
	docker-compose run --rm migrate go run cmd/migrate/main.go apply


migrate-status: ## マイグレーション状態を確認（実際のDB）
	@echo "=== Migration status (actual database) ==="
	docker-compose run --rm migrate atlas migrate status \
		--url "postgres://postgres:postgres@db:5432/financetracker_db?sslmode=disable" \
		--dir "file://cmd/migrate/migrations"

migrate-validate: ## マイグレーションを検証
	@echo "=== Validating migrations ==="
	docker-compose run --rm migrate go run cmd/migrate/main.go validate

migrate-hash: ## マイグレーションのハッシュを再計算（ファイル編集後に必要）
	@echo "=== Recalculating migration hash ==="
	docker-compose run --rm migrate sh -c "cd /app && atlas migrate hash --dir file://cmd/migrate/migrations"

# Atlas CLIの直接実行（必要に応じて）
atlas-migrate-new: ## 新しいマイグレーションファイルを手動作成
	@read -p "Enter migration name: " name; \
	docker-compose run --rm migrate atlas migrate new --env dev $$name

atlas-schema-inspect: ## 現在のスキーマを確認
	docker-compose run --rm migrate atlas schema inspect --env dev

# 本番環境用コマンド
migrate-prod-apply: ## 本番環境にマイグレーションを適用（DATABASE_URLが必要）
	@echo "=== Applying migrations to production ==="
	@echo "WARNING: This will apply migrations to production database!"
	@read -p "Are you sure? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		atlas migrate apply --env prod; \
	else \
		echo "Aborted."; \
	fi

migrate-rollback: ## マイグレーションをロールバック
	@echo "=== Rolling back last migration ==="
	docker-compose run --rm migrate atlas migrate down --env dev

migrate-history: ## マイグレーション履歴を表示（実際のDB）
	@echo "=== Migration history (actual database) ==="
	docker-compose exec db psql -U postgres financetracker_db \
		-c "SELECT version, description, executed_at FROM atlas_schema_revisions ORDER BY executed_at DESC;"

# シードデータ投入
seed: ## シードデータを投入
	@echo "=== Running seed data ==="
	docker-compose run --rm backend sh -c "cd /app && go run cmd/seed/*.go"

seed-prod: ## 本番環境にシードデータを投入
	@echo "=== Running production seed data ==="
	@echo "WARNING: This will seed data to production database!"
	@read -p "Are you sure? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		docker-compose run --rm backend sh -c "cd /app && go run cmd/seed-prod/main.go"; \
	else \
		echo "Aborted."; \
	fi

# 開発用コマンド
dev: ## 開発環境を起動（DB + Backend + Frontend）
	docker-compose up -d db
	@echo "Waiting for database to be ready..."
	@sleep 10
	docker-compose up -d backend frontend pgadmin
	@echo "Services are starting up..."
	@echo "Backend:    http://localhost:8080"
	@echo "Frontend:   http://localhost:3000"
	@echo "pgAdmin:    http://localhost:5050"

dev-logs: ## 開発環境のログを表示
	docker-compose logs -f backend frontend

test: ## テスト実行
	docker-compose run --rm backend go test ./...

fmt: ## コードフォーマット
	docker-compose run --rm backend go fmt ./...

lint: ## リント実行
	docker-compose run --rm backend golangci-lint run

ci-check: ## GitHub Actionsで実行されるバックエンドの全チェックをDocker上で一括実行
	@echo "🚀 Running backend CI checks on Docker (equivalent to GitHub Actions)"
	@echo ""
	@echo "=== Code formatting ==="
	docker-compose run --rm backend sh -c "go fmt ./... && go run golang.org/x/tools/cmd/goimports@v0.26.0 -w ."
	@echo "=== Checking formatting ==="
	@if [ -n "$$(git diff --name-only)" ]; then \
		echo "❌ Code formatting issues found. Please commit the formatting changes."; \
		git diff --name-only; \
		exit 1; \
	else \
		echo "✅ Code formatting check passed"; \
	fi
	@echo "=== Running golangci-lint ==="
	docker-compose run --rm backend golangci-lint run --config=.golangci.yml --timeout=5m
	@echo "=== Running go vet ==="
	docker-compose run --rm backend go vet ./...
	@echo "=== Building ==="
	docker-compose run --rm backend go build  ./...
	@echo "=== Architecture check ==="
	cd backend && chmod +x ./scripts/check-architecture.sh && ./scripts/check-architecture.sh
	@echo ""
	@echo "🎉 All backend CI checks passed successfully!"

# データベース操作
db-create: ## データベースを作成
	docker-compose exec db psql -U postgres -c "CREATE DATABASE financetracker_db WITH ENCODING 'UTF8';"
	docker-compose exec db psql -U postgres -c "CREATE DATABASE financetracker_test WITH ENCODING 'UTF8';"
	docker-compose exec db psql -U postgres -c "CREATE DATABASE financetracker_dev WITH ENCODING 'UTF8';"

db-drop: ## データベースを削除
	docker-compose exec db psql -U postgres -c "DROP DATABASE IF EXISTS financetracker_db;"
	docker-compose exec db psql -U postgres -c "DROP DATABASE IF EXISTS financetracker_test;"
	docker-compose exec db psql -U postgres -c "DROP DATABASE IF EXISTS financetracker_dev;"

db-reset: ## データベースをリセット（削除して再作成）
	@make db-drop
	@make db-create
	@make migrate-apply

# 初期セットアップ
setup: ## 初期セットアップ（ビルド、DB作成、マイグレーション）
	@echo "=== Starting initial setup ==="
	@make build
	@make up -d
	@echo "Waiting for database to be ready..."
	@sleep 15
	@make migrate-apply
	@make seed
	@echo "=== Setup completed ==="
	@echo "Backend:    http://localhost:8080"
	@echo "Frontend:   http://localhost:3000"
	@echo "pgAdmin:    http://localhost:5050"

migrate-gorm: ## GORMでマイグレーションを実行（緊急時のみ）
	@echo "=== Running GORM migration ==="
	docker-compose run --rm migrate go run cmd/migrate/main.go gorm

migrate-init: ## 初期スキーマからAtlasマイグレーションを生成
	@echo "=== Generating initial migration from schema.hcl ==="
	docker-compose run --rm migrate atlas migrate diff initial \
		--env dev \
		--to file://cmd/migrate/schema.hcl