# Backend アーキテクチャ・コーディングルール

## 概要
このプロジェクトは**オニオンアーキテクチャ（Clean Architecture）** を採用しています。  
ビジネスロジックを中心に据え、外部の技術的詳細から独立させることで、保守性・テスタビリティ・拡張性を確保しています。

## ディレクトリ構造

```
backend/
├── cmd/
│   └── lambda/         # Lambda関数のエントリポイント
│       ├── users/
│       ├── accounts/
│       ├── transactions/
│       ├── categories/
│       ├── budgets/
│       ├── reports/
│       ├── auth/
│       └── notifications/
├── internal/
│   ├── domain/         # ドメイン層（最内層）
│   │   ├── user/       # ユーザーコンテキスト
│   │   ├── account/    # アカウントコンテキスト
│   │   ├── transaction/# 取引コンテキスト
│   │   ├── category/   # カテゴリコンテキスト
│   │   ├── budget/     # 予算コンテキスト
│   │   └── common/     # 共通ドメイン要素
│   ├── application/    # アプリケーション層
│   │   ├── dto/        # データ転送オブジェクト
│   │   ├── usecase/    # ユースケース実装
│   │   └── transaction/# トランザクション管理
│   ├── infrastructure/ # インフラストラクチャ層
│   │   ├── database/   # DB接続
│   │   ├── gorm/       # GORM実装
│   │   ├── auth0/      # Auth0実装
│   │   └── aws/        # AWS SDK実装
│   └── interface/      # インターフェース層（プレゼンテーション層）
│       └── lambda/     # Lambdaハンドラー
```

## レイヤー詳細と責務

### 1. ドメイン層 (internal/domain)
**最も内側の層。ビジネスロジックの核心部分。**

#### 構成要素
- **Entity（エンティティ）**: ビジネスオブジェクト。識別子を持つ。
- **Value Object（値オブジェクト）**: 不変のオブジェクト。値の等価性で判断。
- **Repository Interface**: データ永続化の抽象インターフェース。
- **Domain Errors**: ドメイン固有のエラー定義。

#### コーディングルール

##### Entity
```go
// domain/[context]/entity/[entity_name].go
package entity

import (
    "time"
    "[project]/internal/domain/[context]"
    "[project]/internal/domain/[context]/value"
)

// 集約ルートのコメント
type EntityName struct {
    ID          int
    ValueObject *value.ValueObjectName  // 値オブジェクトへの参照
    Name        string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// NewEntityName は新しいエンティティを作成します（ファクトリメソッド）
func NewEntityName(params...) (*EntityName, error) {
    // バリデーション
    if name == "" {
        return nil, context.ErrInvalidName
    }
    
    now := time.Now()
    return &EntityName{
        Name:      name,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

// ビジネスメソッド（振る舞い）
func (e *EntityName) DoBusinessAction() error {
    // ビジネスロジック
    e.UpdatedAt = time.Now()
    return nil
}
```

##### Value Object
```go
// domain/[context]/value/[value_name].go
package value

import "[project]/internal/domain/[context]"

// ValueNameType は値オブジェクトの型定義
type ValueNameType string

// 定数定義
const (
    ValueConstant1 ValueNameType = "value1"
    ValueConstant2 ValueNameType = "value2"
)

// ValueName は値オブジェクト
type ValueName struct {
    ID   int
    Name ValueNameType
}

// NewValueName は新しい値オブジェクトを作成します
func NewValueName(id int, name string) (*ValueName, error) {
    typedName := ValueNameType(name)
    if !typedName.IsValid() {
        return nil, context.ErrInvalidValue
    }
    
    return &ValueName{
        ID:   id,
        Name: typedName,
    }, nil
}

// IsValid は有効性をチェックします
func (v ValueNameType) IsValid() bool {
    switch v {
    case ValueConstant1, ValueConstant2:
        return true
    default:
        return false
    }
}

// ビジネスメソッド
func (v ValueName) CanDoSomething() bool {
    // ビジネスロジック
    return true
}

// Equals は同一性を判定します
func (v ValueName) Equals(other ValueName) bool {
    return v.ID == other.ID
}
```

##### Repository Interface
```go
// domain/[context]/repository/[entity]_repository.go
package repository

import (
    "context"
    "[project]/internal/domain/[context]/entity"
)

type EntityRepository interface {
    GetByID(ctx context.Context, id int) (*entity.EntityName, error)
    GetByIDs(ctx context.Context, ids []int) ([]*entity.EntityName, error)
    Create(ctx context.Context, entity *entity.EntityName) error
    Update(ctx context.Context, entity *entity.EntityName) error
    Delete(ctx context.Context, id int) error
    // 必要に応じて追加のメソッド
}
```

##### Domain Errors
```go
// domain/[context]/errors.go
package [context]

import "errors"

var (
    ErrInvalidName        = errors.New("名前が無効です")
    ErrNotFound          = errors.New("見つかりません")
    ErrAlreadyExists     = errors.New("既に存在します")
    ErrInvalidTransition = errors.New("無効な状態遷移です")
)
```

### 2. アプリケーション層 (internal/application)
**ユースケースを実装。ドメイン層を利用してビジネスフローを制御。**

#### 構成要素
- **UseCase/Service**: ユースケースの実装
- **DTO**: データ転送オブジェクト
- **Transaction Manager**: トランザクション管理

#### コーディングルール

##### UseCase/Service
```go
// application/usecase/[entity]_service.go
package usecase

import (
    "context"
    "[project]/internal/application/dto"
    "[project]/internal/domain/[context]/repository"
)

type EntityService struct {
    entityRepo   repository.EntityRepository
    txManager    transaction.Manager  // 必要に応じて
}

func NewEntityService(repo repository.EntityRepository) *EntityService {
    return &EntityService{
        entityRepo: repo,
    }
}

// ユースケースメソッド
func (s *EntityService) GetEntity(ctx context.Context, id int) (*dto.EntityDTO, error) {
    // リポジトリからエンティティ取得
    entity, err := s.entityRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // DTOに変換
    return dto.FromEntity(entity), nil
}

func (s *EntityService) CreateEntity(ctx context.Context, req *dto.CreateEntityRequest) (*dto.EntityDTO, error) {
    // ドメインエンティティの作成
    entity, err := entity.NewEntityName(req.Name, ...)
    if err != nil {
        return nil, err
    }
    
    // 永続化
    if err := s.entityRepo.Create(ctx, entity); err != nil {
        return nil, err
    }
    
    return dto.FromEntity(entity), nil
}
```

##### DTO
```go
// application/dto/[entity]_dto.go
package dto

import "[project]/internal/domain/[context]/entity"

// リクエストDTO
type CreateEntityRequest struct {
    Name string `json:"name" binding:"required"`
    // 他のフィールド
}

// レスポンスDTO
type EntityDTO struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    // 他のフィールド
}

// エンティティからDTOへの変換
func FromEntity(e *entity.EntityName) *EntityDTO {
    return &EntityDTO{
        ID:   e.ID,
        Name: e.Name,
    }
}
```

### 3. インフラストラクチャ層 (internal/infrastructure)
**技術的詳細の実装。リポジトリの具体実装、外部サービス連携など。**

#### 構成要素
- **Repository実装**: GORMを使用したリポジトリ実装
- **Model**: DBテーブルに対応するGORMモデル
- **外部サービス実装**: JWT、メール送信など

#### コーディングルール

##### GORMモデル
```go
// infrastructure/gorm/model/[entity].go
package model

import (
    "time"
    domainEntity "[project]/internal/domain/[context]/entity"
    "[project]/internal/domain/[context]/value"
)

// EntityName はテーブルのモデル
type EntityName struct {
    ID        int       `gorm:"primaryKey;autoIncrement"`
    Name      string    `gorm:"type:varchar(255);not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
    
    // Relations
    RelatedModel *RelatedModel `gorm:"foreignKey:RelatedID"`
}

func (EntityName) TableName() string {
    return "table_name"
}

// ToEntity はGORMモデルからドメインエンティティへ変換
func (m *EntityName) ToEntity() (*domainEntity.EntityName, error) {
    // 変換ロジック
    return &domainEntity.EntityName{
        ID:        m.ID,
        Name:      m.Name,
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }, nil
}

// FromEntity はドメインエンティティからGORMモデルへ変換
func FromEntity(e *domainEntity.EntityName) *EntityName {
    return &EntityName{
        ID:        e.ID,
        Name:      e.Name,
        CreatedAt: e.CreatedAt,
        UpdatedAt: e.UpdatedAt,
    }
}
```

##### Repository実装
```go
// infrastructure/gorm/repository/[entity]_repository.go
package repository

import (
    "context"
    "gorm.io/gorm"
    
    "[project]/internal/domain/[context]/entity"
    "[project]/internal/domain/[context]/repository"
    "[project]/internal/infrastructure/gorm/model"
)

type entityRepository struct {
    *BaseRepository
}

func NewEntityRepository(db *gorm.DB) repository.EntityRepository {
    return &entityRepository{
        BaseRepository: NewBaseRepository(db),
    }
}

func (r *entityRepository) GetByID(ctx context.Context, id int) (*entity.EntityName, error) {
    var m model.EntityName
    db := r.GetDB(ctx)
    
    if err := db.Preload("RelatedModel").First(&m, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    
    return m.ToEntity()
}

func (r *entityRepository) Create(ctx context.Context, e *entity.EntityName) error {
    m := model.FromEntity(e)
    db := r.GetDB(ctx)
    return db.Create(m).Error
}
```

### 4. インターフェース層 (internal/interface)
**外部とのインターフェース。HTTPコントローラー、ミドルウェアなど。**

#### 構成要素
- **Controller**: HTTPエンドポイント
- **Middleware**: 認証、ロギング、エラーハンドリング
- **Helper**: 共通ヘルパー関数

#### コーディングルール

##### Controller
```go
// interface/controller/[entity]_controller.go
package controller

import (
    "net/http"
    "github.com/gin-gonic/gin"
    
    "[project]/internal/application/dto"
    "[project]/internal/application/usecase"
)

type EntityController struct {
    entityService *usecase.EntityService
}

func NewEntityController(service *usecase.EntityService) *EntityController {
    return &EntityController{
        entityService: service,
    }
}

// HTTPハンドラー
func (c *EntityController) GetEntity(ctx *gin.Context) {
    reqCtx := ctx.Request.Context()
    
    // パラメータ取得
    id, err := getIntParam(ctx, "id")
    if err != nil {
        respondWithError(ctx, err)
        return
    }
    
    // サービス呼び出し
    entity, err := c.entityService.GetEntity(reqCtx, id)
    if err != nil {
        respondWithError(ctx, err)
        return
    }
    
    ctx.JSON(http.StatusOK, entity)
}

func (c *EntityController) CreateEntity(ctx *gin.Context) {
    reqCtx := ctx.Request.Context()
    
    // リクエストバインド
    var req dto.CreateEntityRequest
    if !bindJSON(ctx, &req) {
        return
    }
    
    // サービス呼び出し
    entity, err := c.entityService.CreateEntity(reqCtx, &req)
    if err != nil {
        respondWithError(ctx, err)
        return
    }
    
    ctx.JSON(http.StatusCreated, entity)
}
```

## 重要な設計原則

### 1. 依存性の方向
- 外側の層は内側の層に依存する（内側は外側を知らない）
- インターフェース → アプリケーション → ドメイン
- インフラストラクチャ → ドメイン（リポジトリインターフェースの実装）

### 2. ドメイン駆動設計（DDD）
- **境界づけられたコンテキスト**: user, account, transaction, category, budget
- **集約ルート**: 各コンテキストの主要エンティティ
- **値オブジェクト**: 不変で交換可能なオブジェクト

### 3. エラーハンドリング
- ドメイン層でエラーを定義
- アプリケーション層でエラーを変換・集約
- インターフェース層でHTTPステータスに変換

### 4. トランザクション管理
- アプリケーション層で管理
- リポジトリはトランザクション対応

### 5. 命名規則
- **パッケージ名**: 小文字、単数形
- **ファイル名**: snake_case
- **構造体・インターフェース**: PascalCase
- **メソッド・関数**: PascalCase（公開）、camelCase（非公開）
- **定数**: PascalCase または UPPER_SNAKE_CASE

### 6. テスト
- 各層で単体テスト作成
- モックを使用して層間の依存を切り離す
- ドメインロジックは特に重点的にテスト

## 実装の流れ

1. **ドメイン層から開始**
   - エンティティ、値オブジェクトの定義
   - リポジトリインターフェースの定義
   - ドメインエラーの定義

2. **アプリケーション層**
   - ユースケースの実装
   - DTOの定義
   - トランザクション管理の実装

3. **インフラストラクチャ層**
   - GORMモデルの定義
   - リポジトリの具体実装
   - 外部サービスの実装

4. **インターフェース層**
   - コントローラーの実装
   - ルーティングの設定
   - ミドルウェアの適用

## ベストプラクティス

1. **早期リターン**: エラーチェックは早期にreturn
2. **nilチェック**: ポインタ型は必ずnilチェック
3. **コンテキスト**: 全てのDB操作でcontext.Contextを使用
4. **ログ**: 適切なログレベルで記録
5. **並行性**: goroutineを使う場合は適切な同期処理
6. **リソース管理**: defer文でリソースの解放
7. **エラーメッセージ**: 日本語で分かりやすく

## 禁止事項

1. ドメイン層からインフラストラクチャ層への直接参照
2. コントローラーからリポジトリへの直接アクセス
3. ビジネスロジックのコントローラーへの記述
4. DTOのドメイン層での使用
5. GORMモデルのドメイン層での使用