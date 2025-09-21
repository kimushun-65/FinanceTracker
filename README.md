# FinSight - 賢い家計管理をシンプルに

FinSight は、個人の財務管理を簡単かつ効率的に行うための家計簿アプリケーションです。支出の自動カテゴリ分類、リアルタイム予算管理、AI による支出最適化提案など、先進的な機能を提供します。

## 🎯 プロジェクト概要

### ミッション

ユーザーが日々の支出を簡単に記録し、財務状況を把握し、より良い金銭管理の決定を下せるようサポートすること。

### 主要機能

- 📊 **取引管理**: 収入・支出の記録と管理
- 💰 **予算設定**: カテゴリ別の月次予算設定と追跡
- 📈 **資産管理**: 複数口座の統合管理と資産推移の可視化
- 🤖 **AI 提案**: 過去のデータに基づく予算最適化提案
- 📱 **モバイル対応**: シンプルなモバイルアプリでの素早い入力
- 📧 **レポート機能**: 月次レポートのメール送信

## 🚀 技術スタック

### バックエンド

- **言語**: Go 1.21+
- **フレームワーク**: Gin (Web), GORM (ORM)
- **データベース**: PostgreSQL 15+
- **認証**: Auth0 (Google OAuth)
- **マイグレーション**: Atlas

### フロントエンド

- **Web**: Next.js 14+ with TypeScript
- **状態管理**: Redux Toolkit
- **UI ライブラリ**: Material-UI

### インフラストラクチャ

- **コンテナ**: Docker/Docker Compose
- **CI/CD**: GitHub Actions
- **クラウド**: AWS CDK
- **モニタリング**: CloudWatch

## 🛠️ Makefileコマンド一覧

### 基本操作
| コマンド | 説明 |
|---------|------|
| `make help` | 利用可能なコマンド一覧を表示 |
| `make setup` | 初期セットアップ（ビルド、DB作成、マイグレーション） |

### Docker操作
| コマンド | 説明 |
|---------|------|
| `make up` | すべてのサービスを起動 |
| `make down` | すべてのサービスを停止 |
| `make restart` | すべてのサービスを再起動 |
| `make build` | すべてのサービスをビルド |
| `make logs` | ログを表示（フォロー） |
| `make ps` | サービスの状態を表示 |
| `make clean` | ボリュームを含めて削除（データベースも削除） |

### 開発環境
| コマンド | 説明 |
|---------|------|
| `make dev` | 開発環境を起動（DB + Backend + Frontend） |
| `make dev-logs` | 開発環境のログを表示 |
| `make backend-up` | バックエンドのみ起動 |
| `make backend-logs` | バックエンドのログを表示 |
| `make backend-restart` | バックエンドを再起動 |
| `make backend-shell` | バックエンドコンテナに入る |
| `make pgadmin-up` | pgAdminを起動 |

### データベース操作
| コマンド | 説明 |
|---------|------|
| `make db-shell` | PostgreSQLコンテナに入る |
| `make db-create` | データベースを作成 |
| `make db-drop` | データベースを削除 |
| `make db-reset` | データベースをリセット（削除して再作成） |

### マイグレーション（Atlas）
| コマンド | 説明 |
|---------|------|
| `make migrate-check` | スキーマ差分をチェック |
| `make migrate-diff` | マイグレーション差分を生成 |
| `make migrate-auto` | 自動でスキーマ差分をチェックし、必要に応じてマイグレーション生成・適用 |
| `make migrate-apply` | マイグレーションを適用 |
| `make migrate-status` | マイグレーション状態を確認 |
| `make migrate-validate` | マイグレーションを検証 |
| `make migrate-rollback` | マイグレーションをロールバック |
| `make migrate-history` | マイグレーション履歴を表示 |
| `make atlas-migrate-new` | 新しいマイグレーションファイルを手動作成 |
| `make atlas-schema-inspect` | 現在のスキーマを確認 |

### シードデータ
| コマンド | 説明 |
|---------|------|
| `make seed` | シードデータを投入 |
| `make seed-prod` | 本番環境にシードデータを投入（確認あり） |

### テスト・品質管理
| コマンド | 説明 |
|---------|------|
| `make test` | テスト実行 |
| `make fmt` | コードフォーマット |
| `make lint` | リント実行 |
| `make ci-check` | GitHub Actions相当の全チェックを実行 |

### 本番環境用
| コマンド | 説明 |
|---------|------|
| `make migrate-prod-apply` | 本番環境にマイグレーションを適用（確認あり） |

## 🚦 開発の始め方

1. リポジトリをクローン
```bash
git clone https://github.com/your-org/financetracker.git
cd financetracker
```

2. 環境変数の設定
```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
# 必要に応じて.envファイルを編集
```

3. 初期セットアップと起動
```bash
make setup  # 初回のみ
make dev    # 開発環境起動
```

4. アクセス
- Backend API: http://localhost:8080
- Frontend: http://localhost:3000
- pgAdmin: http://localhost:5050
  - Email: admin@financetracker.com
  - Password: admin

## 📁 プロジェクト構造
```
.
├── backend/           # Goバックエンド（オニオンアーキテクチャ）
│   ├── cmd/          # エントリポイント
│   ├── internal/     # アプリケーション内部実装
│   └── pkg/          # 共有パッケージ
├── frontend/         # Next.jsフロントエンド
├── cdk/              # AWS CDKインフラ
├── docs/             # ドキュメント
├── docker-compose.yml
└── Makefile
```

## 📋 機能一覧

### 1. ダッシュボード

- 今月の収支サマリー表示
- カテゴリ別支出グラフ
- 予算消化状況の可視化
- 最近の取引履歴

### 2. 取引管理

- 収入・支出の記録
- カテゴリ分類
- 日付別フィルタリング
- 検索機能

### 3. 予算管理

- カテゴリ別予算設定
- 前月実績の参照
- AI 予算提案
- 予算超過アラート

### 4. 資産管理

- 複数口座の残高管理
- 資産推移グラフ
- 1 年後の資産予測
- 口座別内訳表示

### 5. レポート

- 期間指定での収支確認
- カテゴリ別内訳
- メールでのレポート送信

### 6. 設定

- カテゴリの ON/OFF 切り替え
- メール通知設定
- プロフィール管理

## 📚 API ドキュメント

API の詳細については以下を参照してください：

- [API 設計書](docs/requirements/api-design.md)

## 📈 パフォーマンス要件

- **レスポンスタイム**: 95%のリクエストが 200ms 以内
- **同時接続数**: 1,000 ユーザー以上のサポート
- **可用性**: 99.9%のアップタイム
- **データ保持**: 5 年間の取引履歴保持

© 2025 FinSight. All rights reserved.
