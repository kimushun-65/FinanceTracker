# Phase 2: ドメイン層実装計画（更新版）

## 概要
FinanceTrackerのドメイン層をクリーンアーキテクチャに基づいて実装します。
ドメイン層は、ビジネスルールとエンティティを定義し、外部の実装詳細から独立させます。

**注**: この計画は`docs/requirements/domain-design.md`との整合性チェックを反映し、YAGNI原則に基づいてシンプル化した更新版です。

## 実装順序と依存関係

### 1. 共通ドメイン要素（基盤）
**優先度**: 最高
**理由**: 全てのエンティティが依存する基本構造

1. `internal/domain/common/base_entity.go`
   - ID (UUID)
   - CreatedAt, UpdatedAt
   - 基本的なバリデーションメソッド

2. 値オブジェクト
   - `internal/domain/common/value/money.go`
     - 金額の表現（通貨、精度管理）
     - 計算メソッド（加算、減算、比較）
   - `internal/domain/common/value/email.go`
     - メールアドレスのバリデーション
   - `internal/domain/common/value/hex_color.go`
     - カラーコードのバリデーション
   - `internal/domain/common/value/time.go`
     - 時刻の表現（HH:mm形式）

3. `internal/domain/common/repository/base_repository.go`
   - 共通リポジトリインターフェース
   - ページネーション
   - ソート機能

4. `internal/domain/common/errors.go`
   - ドメイン固有のエラー定義
   - NotFoundError
   - ValidationError
   - ConflictError


### 2. Userドメイン
**優先度**: 高
**理由**: 全ての機能の基盤となる認証・認可に関わる

1. `internal/domain/user/entity/user.go`
   - ユーザーエンティティ
   - メソッド: IsActive(), UpdateProfile(name, email)

2. `internal/domain/user/value/user_id.go`
   - ユーザーID値オブジェクト

3. `internal/domain/user/value/auth0_id.go`
   - Auth0 ID値オブジェクト

4. `internal/domain/user/repository/user_repository.go`
   - ユーザーリポジトリインターフェース


### 3. Accountドメイン
**優先度**: 高
**理由**: 取引記録の前提となる

1. `internal/domain/account/entity/account.go`
   - 口座エンティティ
   - メソッド: UpdateBalance(amount Money), IsActive(), CanWithdraw(amount Money)

2. `internal/domain/account/entity/account_movement.go`
   - 口座入出金履歴エンティティ
   - メソッド: IsDeposit(), IsWithdrawal()

3. `internal/domain/account/value/account_type.go`
   - 口座タイプ値オブジェクト（checking, savings, investment, credit_card, loan）
   - メソッド: IsAsset(), IsLiability(), CanGoNegative()

4. `internal/domain/account/value/account_name.go`
   - 口座名値オブジェクト

5. `internal/domain/account/value/balance.go`
   - 残高値オブジェクト
   - メソッド: Add(), Subtract(), GetDifference(), IsPositive()

6. `internal/domain/account/repository/account_repository.go`
   - 口座リポジトリインターフェース


### 4. Transactionドメイン
**優先度**: 高
**理由**: コア機能

1. `internal/domain/transaction/entity/transaction.go`
   - 取引エンティティ
   - メソッド: Validate(), CanModify()（1年以内のみ編集可能）

2. `internal/domain/transaction/value/transaction_type.go`
   - 取引タイプ値オブジェクト（income, expense）
   - メソッド: IsIncome(), IsExpense(), GetSign()

3. `internal/domain/transaction/repository/transaction_repository.go`
   - 取引リポジトリインターフェース


### 5. Categoryドメイン
**優先度**: 中
**理由**: 取引の分類に使用

1. `internal/domain/category/entity/category.go`
   - カテゴリエンティティ
   - メソッド: IsActive(), IsCustom()

2. `internal/domain/category/entity/category_master.go`
   - マスターカテゴリエンティティ
   - メソッド: IsActive()

3. `internal/domain/category/value/category_name.go`
   - カテゴリ名値オブジェクト

4. `internal/domain/category/repository/category_repository.go`
   - カテゴリリポジトリインターフェース


### 6. Budgetドメイン
**優先度**: 中
**理由**: 予算管理機能

1. `internal/domain/budget/entity/budget.go`
   - 予算エンティティ
   - メソッド: CalculateRemaining(spent Money), IsOverBudget(spent Money)

2. `internal/domain/budget/entity/budget_suggestion.go`
   - 予算提案エンティティ
   - メソッド: CanAccept(), CanReject()

3. `internal/domain/budget/value/suggestion_status.go`
   - 提案ステータス値オブジェクト（pending, accepted, rejected）
   - メソッド: CanTransitionTo(), IsFinal()

4. `internal/domain/budget/repository/budget_repository.go`
   - 予算リポジトリインターフェース


## 実装原則

### 1. エンティティ設計原則
- ビジネスルールをエンティティ内にカプセル化
- 不変条件の保証（コンストラクタでのバリデーション）
- メソッドはビジネスロジックのみ（永続化は含まない）

### 2. 値オブジェクト設計原則
- イミュータブル（不変）
- 等価性は値で判定
- 自己検証（コンストラクタでバリデーション）

### 3. リポジトリインターフェース設計原則
- ドメイン層に定義（実装はインフラ層）
- ドメインモデルのみを扱う（ORMモデルは扱わない）
- 永続化の詳細を隠蔽

### 4. エラーハンドリング
- ドメイン固有のエラーを定義
- ビジネスルール違反を明確に表現
- 技術的詳細は含まない

## テスト戦略

### 1. エンティティテスト
- ビジネスルールの検証
- 境界値テスト
- 異常系テスト

### 2. 値オブジェクトテスト
- バリデーションテスト
- 等価性テスト
- 演算テスト（該当する場合）

### 3. リポジトリテスト
- モックを使用した単体テスト
- インターフェースの契約検証

## ディレクトリ構造

```
internal/domain/
├── common/
│   ├── base_entity.go
│   ├── errors.go
│   ├── value/
│   │   ├── money.go
│   │   ├── email.go
│   │   ├── hex_color.go
│   │   └── time.go
│   └── repository/
│       └── base_repository.go
├── user/
│   ├── entity/
│   │   └── user.go
│   ├── value/
│   │   ├── user_id.go
│   │   └── auth0_id.go
│   └── repository/
│       └── user_repository.go
├── account/
│   ├── entity/
│   │   ├── account.go
│   │   └── account_movement.go
│   ├── value/
│   │   ├── account_type.go
│   │   ├── account_name.go
│   │   └── balance.go
│   └── repository/
│       └── account_repository.go
├── transaction/
│   ├── entity/
│   │   └── transaction.go
│   ├── value/
│   │   └── transaction_type.go
│   └── repository/
│       └── transaction_repository.go
├── category/
│   ├── entity/
│   │   ├── category.go
│   │   └── category_master.go
│   ├── value/
│   │   └── category_name.go
│   └── repository/
│       └── category_repository.go
└── budget/
    ├── entity/
    │   ├── budget.go
    │   └── budget_suggestion.go
    ├── value/
    │   └── suggestion_status.go
    └── repository/
        └── budget_repository.go
```

## 実装スケジュール（推定）

1. **Day 1**: 共通ドメイン要素（4-5時間）
   - base_entity.go
   - 値オブジェクト（money, email, hex_color, time）
   - base_repository.go
   - errors.go

2. **Day 2**: User & Accountドメイン（4-5時間）
   - Userエンティティと関連要素
   - Accountエンティティと関連要素（AccountMovement含む）

3. **Day 3**: Transaction & Categoryドメイン（4-5時間）
   - Transactionエンティティと関連要素
   - Categoryエンティティと関連要素（CategoryMaster含む）

4. **Day 4**: Budgetドメイン & テスト（4-5時間）
   - Budgetエンティティと関連要素（BudgetSuggestion含む）
   - 全体的なテスト実装

## 成功基準

1. **アーキテクチャ準拠**
   - 依存関係の方向が内向き（ドメイン層は外部に依存しない）
   - ビジネスロジックがドメイン層に集約

2. **テストカバレッジ**
   - 各エンティティ・値オブジェクトのテストカバレッジ90%以上

3. **保守性**
   - 明確な責務分離
   - 拡張しやすい設計

4. **シンプルさ**
   - YAGNI原則に従った実装
   - 過度な抽象化の回避

## 注意事項

1. **GORMモデルとの分離**
   - ドメインエンティティとGORMモデルは完全に分離
   - マッピングはインフラ層で実装

2. **過度な抽象化の回避**
   - YAGNIの原則に従い、必要最小限の抽象化
   - 実際の要求に基づいた設計

3. **パフォーマンス考慮**
   - N+1問題を意識したリポジトリインターフェース設計
   - 必要に応じた集約の最適化

4. **ドメイン設計書との整合性**
   - `docs/requirements/domain-design.md`のビジネスルールを必ず反映
   - エンティティの振る舞いや値オブジェクトの詳細は設計書を参照

## YAGNIに基づく除外事項

以下の機能は現時点では実装しません（将来必要になった際に追加）：

1. **ドメインイベント**
   - 集約間の連携は当面アプリケーション層で直接実装
   - イベント駆動アーキテクチャは複雑性に見合うメリットが現時点では不明

2. **複雑な集約境界**
   - 単純なエンティティ関係から開始
   - 実際の運用で問題が発生してから最適化

この原則により、開発速度を向上し、実際の要求に基づいて段階的に複雑性を追加します。