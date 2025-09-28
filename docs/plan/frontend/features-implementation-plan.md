# Frontend Features層実装計画書（APIカスタムフック）

## 概要

Feature-Sliced Design (FSD) に基づき、各entityに対応するReact QueryのカスタムフックでCRUD操作を提供するfeatures層を実装します。この層はentities層のAPIクライアントを使用し、データフェッチング機能を提供します。

**注意**: この計画書はAPIカスタムフック部分のみに焦点を当てており、UI コンポーネントは含まれません。

## ディレクトリ構造

```
frontend/src/features/
├── account-management/
│   ├── useAccounts.ts           # アカウント一覧取得
│   ├── useAccount.ts            # 個別アカウント取得
│   ├── useCreateAccount.ts      # アカウント作成
│   ├── useUpdateAccount.ts      # アカウント更新
│   ├── useDeleteAccount.ts      # アカウント削除
│   └── index.ts

├── transaction-management/
│   ├── useTransactions.ts       # 取引一覧取得（フィルタ対応）
│   ├── useTransaction.ts        # 個別取引取得
│   ├── useCreateTransaction.ts  # 取引作成
│   ├── useUpdateTransaction.ts  # 取引更新
│   ├── useDeleteTransaction.ts  # 取引削除
│   ├── useTransactionAggregates.ts # 集計データ取得
│   ├── useTransactionFilters.ts # フィルタ状態管理
│   └── index.ts

├── category-management/
│   ├── useCategories.ts         # カテゴリ一覧取得
│   ├── useCategoryMasters.ts    # マスターカテゴリ取得
│   ├── useCategory.ts           # 個別カテゴリ取得
│   ├── useCreateCategory.ts     # カテゴリ作成
│   ├── useUpdateCategory.ts     # カテゴリ更新
│   ├── useDeleteCategory.ts     # カテゴリ削除
│   └── index.ts

├── user-profile/
│   ├── useUserProfile.ts        # ユーザープロフィール取得
│   ├── useUpdateProfile.ts      # プロフィール更新
│   └── index.ts

├── budget-management/
│   ├── useBudgets.ts           # 予算一覧取得
│   ├── useBudget.ts            # 個別予算取得
│   ├── useBudgetSummary.ts     # 予算サマリー取得
│   ├── useCreateBudget.ts      # 予算作成
│   ├── useUpdateBudget.ts      # 予算更新
│   ├── useDeleteBudget.ts      # 予算削除
│   └── index.ts

└── index.ts                        # features層の公開API集約
```

## 実装方針

### 1. カスタムフックの設計パターン

#### 基本的なCRUDフック

**クエリフック例**
```typescript
// features/account-management/useAccounts.ts
import { useQuery } from '@tanstack/react-query';
import { accountApi, accountKeys } from '@/entities/account';

export const useAccounts = () => {
  return useQuery({
    queryKey: accountKeys.lists(),
    queryFn: accountApi.list,
    staleTime: 5 * 60 * 1000, // 5分
  });
};
```

**ミューテーションフック例**
```typescript
// features/account-management/useCreateAccount.ts
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { accountApi, accountKeys, type CreateAccountPayload } from '@/entities/account';

export const useCreateAccount = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: accountApi.create,
    onSuccess: () => {
      // キャッシュを無効化してリフェッチ
      queryClient.invalidateQueries({ queryKey: accountKeys.lists() });
    },
  });
};
```

#### フィルタ対応のクエリフック

```typescript
// features/transaction-management/useTransactions.ts
import { useQuery } from '@tanstack/react-query';
import { transactionApi, transactionKeys, type TransactionListParams } from '@/entities/transaction';

export const useTransactions = (params: TransactionListParams = {}) => {
  return useQuery({
    queryKey: transactionKeys.list(params),
    queryFn: () => transactionApi.list(params),
    staleTime: 2 * 60 * 1000, // 2分
    enabled: true, // パラメータに応じて有効/無効を制御可能
  });
};
```

#### 楽観的更新パターン

```typescript
// features/account-management/useUpdateAccount.ts
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { accountApi, accountKeys, type Account, type UpdateAccountPayload } from '@/entities/account';

export const useUpdateAccount = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, ...payload }: { id: string } & UpdateAccountPayload) =>
      accountApi.update(id, payload),
    
    // 楽観的更新
    onMutate: async ({ id, ...newData }) => {
      await queryClient.cancelQueries({ queryKey: accountKeys.detail(id) });
      
      const previousAccount = queryClient.getQueryData<Account>(accountKeys.detail(id));
      
      queryClient.setQueryData<Account>(accountKeys.detail(id), (old) => ({
        ...old!,
        ...newData,
        updatedAt: new Date().toISOString(),
      }));

      return { previousAccount };
    },
    
    onError: (err, { id }, context) => {
      if (context?.previousAccount) {
        queryClient.setQueryData(accountKeys.detail(id), context.previousAccount);
      }
    },
    
    onSuccess: (data, { id }) => {
      queryClient.setQueryData(accountKeys.detail(id), data);
      queryClient.invalidateQueries({ queryKey: accountKeys.lists() });
    },
  });
};
```


### 2. 状態管理フック

```typescript
// features/transaction-management/useTransactionFilters.ts
import { useState, useCallback } from 'react';
import type { TransactionListParams, TransactionType } from '@/entities/transaction';

export const useTransactionFilters = () => {
  const [filters, setFilters] = useState<TransactionListParams>({});

  const setAccountFilter = useCallback((accountId?: string) => {
    setFilters(prev => ({ ...prev, accountId }));
  }, []);

  const setTypeFilter = useCallback((type?: TransactionType) => {
    setFilters(prev => ({ ...prev, type }));
  }, []);

  const setDateRange = useCallback((startDate?: string, endDate?: string) => {
    setFilters(prev => ({ ...prev, startDate, endDate }));
  }, []);

  const clearFilters = useCallback(() => {
    setFilters({});
  }, []);

  return {
    filters,
    setAccountFilter,
    setTypeFilter,
    setDateRange,
    clearFilters,
  };
};
```

### 3. 集計フック例

```typescript
// features/transaction-management/useTransactionAggregates.ts
import { useQuery } from '@tanstack/react-query';
import { transactionApi, transactionKeys, type TransactionListParams } from '@/entities/transaction';
import { calculateMonthlyAggregates } from '@/entities/transaction';

export const useTransactionAggregates = (params: TransactionListParams = {}) => {
  return useQuery({
    queryKey: [...transactionKeys.list(params), 'aggregates'],
    queryFn: async () => {
      const response = await transactionApi.list(params);
      // entities層の集計関数を使用
      return calculateMonthlyAggregates(response.transactions);
    },
    staleTime: 5 * 60 * 1000, // 5分
  });
};
```


## エラーハンドリング戦略

### 1. グローバルエラーハンドリング

```typescript
// shared/api/errorHandler.ts
export const globalErrorHandler = (error: unknown) => {
  if (error instanceof Error) {
    // ログ送信、トースト表示など
    console.error('API Error:', error.message);
  }
};

// React Query設定でグローバルに適用
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      onError: globalErrorHandler,
      retry: (failureCount, error) => {
        // 認証エラーの場合はリトライしない
        if (error?.status === 401) return false;
        return failureCount < 3;
      },
    },
    mutations: {
      onError: globalErrorHandler,
    },
  },
});
```

### 2. フィーチャー固有のエラーハンドリング

```typescript
// features/account-management/useCreateAccount.ts
export const useCreateAccount = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: accountApi.create,
    onError: (error) => {
      // アカウント作成固有のエラー処理
      if (error.code === 'DUPLICATE_ACCOUNT_NAME') {
        toast.error('同じ名前のアカウントが既に存在します');
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: accountKeys.lists() });
      toast.success('アカウントを作成しました');
    },
  });
};
```

## キャッシュ戦略

### 1. ステイル時間設定

```typescript
// 各エンティティごとのキャッシュ設定
const CACHE_CONFIG = {
  accounts: {
    staleTime: 5 * 60 * 1000,     // 5分 - 残高が変わりにくい
    cacheTime: 10 * 60 * 1000,    // 10分
  },
  transactions: {
    staleTime: 2 * 60 * 1000,     // 2分 - 頻繁に追加される
    cacheTime: 5 * 60 * 1000,     // 5分
  },
  categories: {
    staleTime: 30 * 60 * 1000,    // 30分 - ほとんど変わらない
    cacheTime: 60 * 60 * 1000,    // 1時間
  },
  budgets: {
    staleTime: 10 * 60 * 1000,    // 10分 - 使用状況が変わる
    cacheTime: 20 * 60 * 1000,    // 20分
  },
};
```

### 2. キャッシュ無効化パターン

```typescript
// 関連データの無効化
export const useCreateTransaction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: transactionApi.create,
    onSuccess: (data) => {
      // 取引リストを無効化
      queryClient.invalidateQueries({ queryKey: transactionKeys.lists() });
      
      // 関連アカウントの残高を無効化
      queryClient.invalidateQueries({ 
        queryKey: accountKeys.detail(data.accountId) 
      });
      
      // 予算の使用状況を無効化
      queryClient.invalidateQueries({ 
        queryKey: budgetKeys.lists() 
      });
    },
  });
};
```

## パフォーマンス最適化

### 1. 無限スクロール対応

```typescript
// features/transaction-management/useInfiniteTransactions.ts
import { useInfiniteQuery } from '@tanstack/react-query';
import { transactionApi, transactionKeys } from '@/entities/transaction';

export const useInfiniteTransactions = (params: TransactionListParams = {}) => {
  return useInfiniteQuery({
    queryKey: transactionKeys.list(params),
    queryFn: ({ pageParam = 0 }) => 
      transactionApi.list({ ...params, offset: pageParam }),
    getNextPageParam: (lastPage, pages) => {
      if (lastPage.hasMore) {
        return pages.length * 20; // 20件ずつ取得
      }
      return undefined;
    },
    staleTime: 2 * 60 * 1000,
  });
};
```

### 2. 選択的refetch

```typescript
// 特定の条件でのみrefetchを実行
export const useTransactions = (params: TransactionListParams = {}) => {
  return useQuery({
    queryKey: transactionKeys.list(params),
    queryFn: () => transactionApi.list(params),
    refetchOnWindowFocus: false, // ウィンドウフォーカス時はrefetchしない
    refetchOnMount: 'always',    // マウント時は常にrefetch
    refetchInterval: 5 * 60 * 1000, // 5分ごとに自動refetch
  });
};
```

## 実装優先順位

### Phase 1: 基盤整備（高優先度）
1. **account-management** - アカウント管理機能
   - useAccounts, useAccount, useCreateAccount, useUpdateAccount, useDeleteAccount
   - useAccountValidation, useAccountCalculations
   
2. **transaction-management** - 取引管理機能
   - useTransactions, useTransaction, useCreateTransaction, useUpdateTransaction, useDeleteTransaction
   - useTransactionAggregates, useTransactionFilters, useTransactionValidation, useTransactionSort
   
3. **category-management** - カテゴリ管理機能
   - useCategories, useCategoryMasters, useCategory, useCreateCategory, useUpdateCategory, useDeleteCategory
   - useCategoryValidation, useCategorySort

### Phase 2: 拡張機能（中優先度）
4. **budget-management** - 予算管理機能
   - useBudgets, useBudget, useBudgetSummary, useCreateBudget, useUpdateBudget, useDeleteBudget
   - useBudgetValidation, useBudgetCalculations
   
5. **user-profile** - ユーザープロフィール機能
   - useUserProfile, useUpdateProfile
   - useProfileValidation

### Phase 3: 高度な機能（低優先度）
6. その他の分析・レポート機能
7. 無限スクロール対応
8. リアルタイム更新機能

## テスト戦略

### 1. カスタムフックのテスト

```typescript
// features/account-management/__tests__/useAccounts.test.ts
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useAccounts } from '../useAccounts';

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  
  return ({ children }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

test('should fetch accounts successfully', async () => {
  const { result } = renderHook(() => useAccounts(), {
    wrapper: createWrapper(),
  });

  await waitFor(() => expect(result.current.isSuccess).toBe(true));
  expect(result.current.data).toBeDefined();
});
```

### 2. ミューテーションのテスト

```typescript
// features/account-management/__tests__/useCreateAccount.test.ts
import { renderHook, waitFor } from '@testing-library/react';
import { act } from '@testing-library/react';
import { useCreateAccount } from '../useCreateAccount';

test('should create account successfully', async () => {
  const { result } = renderHook(() => useCreateAccount(), {
    wrapper: createWrapper(),
  });

  const mockPayload = {
    name: 'Test Account',
    accountType: 'checking' as const,
  };

  act(() => {
    result.current.mutate(mockPayload);
  });

  await waitFor(() => expect(result.current.isSuccess).toBe(true));
  expect(result.current.data).toBeDefined();
});
```


## 注意事項

### 1. 型安全性の確保
- 全てのフックで適切な型定義を使用
- entities層の型を再利用してDRYを保つ
- ジェネリクス使用で再利用性を高める

### 2. エラーハンドリング
- ネットワークエラー、認証エラーの適切な処理
- ユーザーフレンドリーなエラーメッセージ
- リトライ戦略の実装

### 3. パフォーマンス
- 不要なリレンダリングの防止
- 適切なメモ化の使用
- バンドルサイズの最適化

## 次のステップ

1. **Phase 1の実装開始**
   - account-management featureの基本的なCRUDフック
   - バリデーションと計算フック
   
2. **テスト環境の整備**
   - MSWを使用したAPIモック
   - カスタムフックのテストユーティリティの作成
   
3. **段階的な機能拡張**
   - 基本的なCRUDフックの安定化後、高度な機能を追加
   - パフォーマンス最適化の実装

これでentities層の上に構築される強固なfeatures層のAPIカスタムフックが整い、保守性と拡張性の高いデータフェッチング機能を提供できます。