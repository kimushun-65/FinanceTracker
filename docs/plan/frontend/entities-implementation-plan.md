# Frontend Entities層実装計画書

## 概要

Feature-Sliced Design (FSD) に基づき、バックエンドのドメインモデルと整合性を保ちながら、フロントエンドのentities層を実装します。entities層は、ビジネスエンティティの型定義、APIクライアント、ビジネスロジックを管理する責務を持ちます。

## 実装方針

### 1. ディレクトリ構造

```
src/entities/
├── shared/                             # 共通リソース
│   ├── types/
│   │   ├── base.ts                    # BaseEntity, ApiResponse等の基本型
│   │   ├── api.ts                     # APIエラー型、ページネーション型
│   │   └── index.ts
│   ├── value-objects/
│   │   ├── money/
│   │   │   ├── money.types.ts         # Money, Currency型定義
│   │   │   ├── money.utils.ts         # 金額計算、フォーマット関数
│   │   │   ├── money.constants.ts     # 通貨定数
│   │   │   └── index.ts
│   │   ├── email/
│   │   │   ├── email.types.ts         # Email型定義
│   │   │   ├── email.validators.ts    # メールバリデーション
│   │   │   └── index.ts
│   │   ├── date/
│   │   │   ├── date.types.ts          # 日付関連の型
│   │   │   ├── date.utils.ts          # 日付フォーマット、計算
│   │   │   └── index.ts
│   │   └── index.ts
│   └── constants/
│       ├── api.constants.ts           # APIバージョン、タイムアウト設定
│       └── index.ts
│
├── user/
│   ├── model/
│   │   ├── user.types.ts              # User, UpdateUserPayload等
│   │   ├── user.schema.ts             # zodスキーマ定義
│   │   └── index.ts
│   ├── api/
│   │   ├── user.client.ts             # userApi実装
│   │   ├── user.endpoints.ts          # エンドポイント定義
│   │   ├── user.keys.ts               # React Queryキー定義
│   │   └── index.ts
│   ├── lib/
│   │   ├── user.validations.ts        # ユーザー名、メール等のバリデーション
│   │   ├── user.transformers.ts       # APIレスポンス変換
│   │   ├── user.utils.ts              # 表示名取得等のユーティリティ
│   │   └── index.ts
│   └── index.ts                       # Public API (型、API、ユーティリティのみ export)
│
├── account/
│   ├── model/
│   │   ├── account.types.ts           # Account, AccountType, CreateAccountPayload等
│   │   ├── account.schema.ts          # zodスキーマ
│   │   ├── account.constants.ts       # アカウントタイプ定数、制限値
│   │   └── index.ts
│   ├── api/
│   │   ├── account.client.ts          # accountApi実装
│   │   ├── account.endpoints.ts       # エンドポイント定義
│   │   ├── account.keys.ts            # React Queryキー定義
│   │   └── index.ts
│   ├── lib/
│   │   ├── account.validations.ts     # 口座名、残高等のバリデーション
│   │   ├── account.transformers.ts    # APIレスポンス変換
│   │   ├── account.calculations.ts    # 残高計算、出金可能判定
│   │   ├── account.formatters.ts      # 口座情報フォーマット
│   │   └── index.ts
│   ├── ui/
│   │   ├── AccountTypeIcon/
│   │   │   ├── AccountTypeIcon.tsx
│   │   │   └── index.ts
│   │   ├── AccountTypeBadge/
│   │   │   ├── AccountTypeBadge.tsx
│   │   │   └── index.ts
│   │   └── index.ts
│   └── index.ts
│
├── transaction/
│   ├── model/
│   │   ├── transaction.types.ts       # Transaction, TransactionType等
│   │   ├── transaction.schema.ts      # zodスキーマ
│   │   ├── transaction.constants.ts   # 取引タイプ定数、制限値
│   │   └── index.ts
│   ├── api/
│   │   ├── transaction.client.ts      # transactionApi実装
│   │   ├── transaction.endpoints.ts   # エンドポイント定義
│   │   ├── transaction.keys.ts        # React Queryキー定義
│   │   ├── transaction.filters.ts     # クエリパラメータ型、フィルタ関数
│   │   └── index.ts
│   ├── lib/
│   │   ├── transaction.validations.ts # 金額、説明等のバリデーション
│   │   ├── transaction.transformers.ts# APIレスポンス変換
│   │   ├── transaction.aggregations.ts# 月次集計、カテゴリ別集計
│   │   ├── transaction.formatters.ts  # 取引情報フォーマット
│   │   ├── transaction.sorters.ts     # ソート関数
│   │   └── index.ts
│   ├── ui/
│   │   ├── TransactionTypeBadge/
│   │   │   ├── TransactionTypeBadge.tsx
│   │   │   └── index.ts
│   │   ├── TransactionAmountDisplay/
│   │   │   ├── TransactionAmountDisplay.tsx
│   │   │   └── index.ts
│   │   └── index.ts
│   └── index.ts
│
├── category/
│   ├── model/
│   │   ├── category.types.ts          # Category, CategoryMaster, CategoryType等
│   │   ├── category.schema.ts         # zodスキーマ
│   │   ├── category.constants.ts      # カテゴリタイプ、デフォルトアイコン
│   │   └── index.ts
│   ├── api/
│   │   ├── category.client.ts         # categoryApi実装
│   │   ├── category.endpoints.ts      # エンドポイント定義
│   │   ├── category.keys.ts           # React Queryキー定義
│   │   └── index.ts
│   ├── lib/
│   │   ├── category.validations.ts    # カテゴリ名のバリデーション
│   │   ├── category.transformers.ts   # APIレスポンス変換
│   │   ├── category.utils.ts          # カテゴリ関連ユーティリティ
│   │   └── index.ts
│   ├── ui/
│   │   ├── CategoryIcon/
│   │   │   ├── CategoryIcon.tsx       # アイコン表示コンポーネント
│   │   │   └── index.ts
│   │   ├── CategoryColorBadge/
│   │   │   ├── CategoryColorBadge.tsx # カラー表示バッジ
│   │   │   └── index.ts
│   │   └── index.ts
│   └── index.ts
│
├── budget/
│   ├── model/
│   │   ├── budget.types.ts            # Budget, PeriodType, CreateBudgetPayload等
│   │   ├── budget.schema.ts           # zodスキーマ
│   │   ├── budget.constants.ts        # 期間タイプ定数
│   │   └── index.ts
│   ├── api/
│   │   ├── budget.client.ts           # budgetApi実装
│   │   ├── budget.endpoints.ts        # エンドポイント定義
│   │   ├── budget.keys.ts             # React Queryキー定義
│   │   └── index.ts
│   ├── lib/
│   │   ├── budget.validations.ts      # 予算額、期間のバリデーション
│   │   ├── budget.transformers.ts     # APIレスポンス変換
│   │   ├── budget.calculations.ts     # 予算消化率計算、残額計算
│   │   ├── budget.checkers.ts         # 予算期間チェック、超過判定
│   │   └── index.ts
│   ├── ui/
│   │   ├── BudgetProgressBar/
│   │   │   ├── BudgetProgressBar.tsx  # 予算消化率プログレスバー
│   │   │   └── index.ts
│   │   ├── BudgetStatusBadge/
│   │   │   ├── BudgetStatusBadge.tsx  # 予算状態バッジ
│   │   │   └── index.ts
│   │   └── index.ts
│   └── index.ts
│
├── budget-suggestion/                  # 予算提案
│   ├── model/
│   │   ├── suggestion.types.ts
│   │   ├── suggestion.schema.ts
│   │   └── index.ts
│   ├── api/
│   │   ├── suggestion.client.ts
│   │   ├── suggestion.endpoints.ts
│   │   └── index.ts
│   ├── lib/
│   │   └── index.ts
│   └── index.ts
│
├── asset-snapshot/                     # 資産スナップショット
│   ├── model/
│   │   ├── snapshot.types.ts
│   │   ├── snapshot.schema.ts
│   │   └── index.ts
│   ├── api/
│   │   ├── snapshot.client.ts
│   │   ├── snapshot.endpoints.ts
│   │   └── index.ts
│   ├── lib/
│   │   ├── snapshot.aggregations.ts   # 資産集計
│   │   └── index.ts
│   └── index.ts
│
├── notification-settings/              # 通知設定
│   ├── model/
│   │   ├── notification.types.ts
│   │   ├── notification.schema.ts
│   │   └── index.ts
│   ├── api/
│   │   ├── notification.client.ts
│   │   ├── notification.endpoints.ts
│   │   └── index.ts
│   └── index.ts
│
└── index.ts                           # エンティティ層の公開API集約
```

### 2. 型定義の設計

#### BaseEntity（共通基底型）
```typescript
// src/entities/shared/types/base.ts
export type BaseEntity = {
  id: string; // UUID
  createdAt: string; // ISO 8601
  updatedAt: string; // ISO 8601
};
```

#### 値オブジェクトの設計
```typescript
// src/entities/shared/value-objects/money.ts
export type Money = {
  amount: string; // decimal string to avoid floating point issues
  currency: Currency;
};

export type Currency = 'JPY' | 'USD' | 'EUR';

// src/entities/shared/value-objects/email.ts
export type Email = string & { readonly __brand: unique symbol };

export const createEmail = (value: string): Email => {
  if (!isValidEmail(value)) {
    throw new Error('Invalid email format');
  }
  return value as Email;
};
```

### 3. 各エンティティの実装計画

#### User Entity

```typescript
// src/entities/user/model/types.ts
import { BaseEntity } from '@/entities/shared';
import { Email } from '@/entities/shared/value-objects';

export type User = BaseEntity & {
  auth0Id: string;
  email: Email;
  name: string;
  emailVerified: boolean;
};

export type UpdateUserPayload = {
  name?: string;
  email?: string;
};
```

#### Account Entity

```typescript
// src/entities/account/model/types.ts
export type AccountType = 'checking' | 'investment' | 'cash';

export type Account = BaseEntity & {
  userId: string;
  name: string;
  accountType: AccountType;
  balance: Money;
};

export type CreateAccountPayload = {
  name: string;
  accountType: AccountType;
  initialBalance?: Money;
};

export type UpdateAccountPayload = {
  name?: string;
  accountType?: AccountType;
  balance?: Money;
};
```

#### Transaction Entity

```typescript
// src/entities/transaction/model/types.ts
export type TransactionType = 'income' | 'expense';

export type Transaction = BaseEntity & {
  userId: string;
  accountId: string;
  categoryId: string;
  amount: Money;
  transactionType: TransactionType;
  description: string;
  transactionDate: string; // ISO 8601
};

export type CreateTransactionPayload = {
  accountId: string;
  categoryId: string;
  amount: Money;
  transactionType: TransactionType;
  description: string;
  transactionDate?: string;
};
```

#### Category Entity

```typescript
// src/entities/category/model/types.ts
export type CategoryType = 'income' | 'expense';

export type Category = BaseEntity & {
  userId: string;
  categoryMasterId: string;
  name: string;
  customName?: string;
  isActive: boolean;
};

export type CategoryMaster = BaseEntity & {
  name: string;
  categoryType: CategoryType;
  icon: string;
  color?: string; // HexColor
  displayOrder: number;
};
```

#### Budget Entity

```typescript
// src/entities/budget/model/types.ts
export type PeriodType = 'monthly' | 'yearly';

export type Budget = BaseEntity & {
  userId: string;
  categoryId: string;
  amount: Money;
  period: PeriodType;
  startDate: string;
  endDate?: string;
  isActive: boolean;
};
```

### 4. APIクライアントの実装方針

```typescript
// src/entities/account/api/client.ts
import { apiClient } from '@/shared/api';
import { Account, CreateAccountPayload } from '../model';
import { endpoints } from './endpoints';

export const accountApi = {
  list: async (): Promise<Account[]> => {
    const response = await apiClient.get(endpoints.list);
    return response.data;
  },
  
  get: async (id: string): Promise<Account> => {
    const response = await apiClient.get(endpoints.get(id));
    return response.data;
  },
  
  create: async (payload: CreateAccountPayload): Promise<Account> => {
    const response = await apiClient.post(endpoints.create, payload);
    return response.data;
  },
  
  update: async (id: string, payload: UpdateAccountPayload): Promise<Account> => {
    const response = await apiClient.put(endpoints.update(id), payload);
    return response.data;
  },
  
  delete: async (id: string): Promise<void> => {
    await apiClient.delete(endpoints.delete(id));
  },
};
```

### 5. ビジネスロジックの実装

```typescript
// src/entities/account/lib/validations.ts
export const validateAccountName = (name: string): boolean => {
  return name.length >= 1 && name.length <= 100;
};

export const canWithdraw = (account: Account, amount: Money): boolean => {
  if (account.accountType === 'cash') {
    return parseFloat(account.balance.amount) >= parseFloat(amount.amount);
  }
  return true;
};

// src/entities/transaction/lib/transformers.ts
export const groupTransactionsByMonth = (transactions: Transaction[]): Map<string, Transaction[]> => {
  // 月ごとにグループ化するロジック
};

export const calculateCategoryTotals = (transactions: Transaction[]): Map<string, Money> => {
  // カテゴリごとの合計を計算
};
```

### 6. UI コンポーネント（entities層に含まれる共通UI）

```typescript
// src/entities/account/ui/AccountTypeIcon.tsx
export const AccountTypeIcon: FC<{ type: AccountType }> = ({ type }) => {
  const icons = {
    checking: <BankIcon />,
    investment: <TrendingUpIcon />,
    cash: <WalletIcon />,
  };
  return icons[type];
};

// src/entities/transaction/ui/TransactionTypeBadge.tsx
export const TransactionTypeBadge: FC<{ type: TransactionType }> = ({ type }) => {
  return (
    <Badge variant={type === 'income' ? 'success' : 'destructive'}>
      {type === 'income' ? '収入' : '支出'}
    </Badge>
  );
};
```

## 実装優先順位

1. **Phase 1: 基盤整備**（必須）
   - shared/types/base.ts
   - shared/value-objects/（money, email）
   - shared/api/client.ts（共通APIクライアント設定）

2. **Phase 2: 認証・ユーザー関連**（高優先度）
   - user entity全体
   - 認証状態管理の基盤

3. **Phase 3: コア機能エンティティ**（高優先度）
   - account entity
   - transaction entity
   - category entity

4. **Phase 4: 拡張機能エンティティ**（中優先度）
   - budget entity
   - budgetSuggestion entity
   - assetSnapshot entity

5. **Phase 5: その他のエンティティ**（低優先度）
   - notificationSettings entity
   - assetForecast entity
   - accountMovement entity

## テスト戦略

1. **型定義のテスト**
   - TypeScriptの型チェックで基本的な整合性を保証
   - 値オブジェクトの生成関数は単体テストを実装

2. **APIクライアントのテスト**
   - MSW（Mock Service Worker）を使用したモックテスト
   - エラーケースの処理確認

3. **ビジネスロジックのテスト**
   - バリデーション関数の単体テスト
   - データ変換関数の単体テスト

## 注意事項

1. **バックエンドとの整合性**
   - バックエンドの値オブジェクトの制約を尊重
   - APIレスポンスの型は実際のレスポンスと厳密に一致させる

2. **エラーハンドリング**
   - APIエラーは適切に型付けしてハンドリング
   - バリデーションエラーはユーザーフレンドリーなメッセージに変換

3. **パフォーマンス**
   - 大量データを扱う場合は適切なページネーション
   - 必要に応じてメモ化や最適化を実施

4. **拡張性**
   - 新しいエンティティやフィールドの追加が容易な設計
   - ドメインロジックの変更が他層に影響しない設計

## 次のステップ

1. shared層の基盤コードを実装
2. user entityの実装とテスト
3. 認証フローとの統合確認
4. 他のエンティティを優先順位に従って順次実装