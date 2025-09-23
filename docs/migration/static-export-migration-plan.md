# 静的エクスポート移行計画

## 概要
Next.js SSRモードから静的エクスポートモードへの移行を行い、AWS Amplifyでのデプロイを成功させる。

## 移行による影響

### 使えなくなる機能
- `middleware.ts` によるサーバーサイドの認証チェック
- サーバーサイドでの自動リダイレクト
- API Routes（`/api/*`）
- `getServerSideProps` / `getServerSideAuth`

### 維持される機能
- Auth0によるクライアントサイド認証
- HttpOnlyクッキーを使用したトークン管理（バックエンドAPI経由）
- ユーザー認証状態の管理
- 保護されたルートへのアクセス制御（クライアントサイド）

## 実装手順

### 1. Next.js設定の変更

**ファイル**: `frontend/next.config.ts`
```typescript
import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  output: 'export',
  distDir: 'out',
  images: {
    unoptimized: true,
  },
};

export default nextConfig;
```

### 2. Amplify設定の更新

**ファイル**: `amplify.yml`
```yaml
version: 1
frontend:
  phases:
    preBuild:
      commands:
        - cd frontend
        - npm ci
    build:
      commands:
        - npm run build
  artifacts:
    baseDirectory: frontend/out
    files:
      - '**/*'
  cache:
    paths:
      - frontend/node_modules/**/*
```

### 3. middleware.tsの無効化または削除

**オプション1**: ファイルをリネーム
```bash
mv frontend/src/middleware.ts frontend/src/middleware.ts.bak
```

**オプション2**: 内容をコメントアウト
```typescript
// middleware.tsの全内容をコメントアウト
```

### 4. クライアントサイド認証ガードの実装

**新規ファイル**: `frontend/src/shared/components/ProtectedRoute.tsx`
```typescript
import { useAuth0 } from '@auth0/auth0-react';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth0();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/');
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-12 w-12 animate-spin rounded-full border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return <>{children}</>;
};
```

### 5. 保護されたページの更新

**例**: `frontend/src/app/dashboard/page.tsx`
```typescript
import { ProtectedRoute } from '@/shared/components/ProtectedRoute';
import { DashboardContainer } from '@/page-components/dashboard';

export default function DashboardPage() {
  return (
    <ProtectedRoute>
      <DashboardContainer />
    </ProtectedRoute>
  );
}
```

### 6. useAuthWithCookieフックの調整

現在の実装は基本的にクライアントサイドで動作するため、大きな変更は不要。
ただし、ページロード時の認証チェックロジックを確認する必要がある。

### 7. 環境変数の確認

Amplifyの環境変数はそのまま使用可能：
- `NEXT_PUBLIC_API_BASE_URL`
- `NEXT_PUBLIC_AUTH0_AUDIENCE`
- `NEXT_PUBLIC_AUTH0_CLIENT_ID`
- `NEXT_PUBLIC_AUTH0_DOMAIN`
- `NEXT_PUBLIC_AUTH0_REDIRECT_URI`

## テスト項目

1. **ログイン/ログアウトフロー**
   - [ ] ログインボタンからAuth0へのリダイレクト
   - [ ] Auth0からのコールバック処理
   - [ ] トークンのHttpOnlyクッキーへの保存
   - [ ] ログアウト処理とクッキーの削除

2. **ルート保護**
   - [ ] 未認証時の`/dashboard`へのアクセス → ホームページへリダイレクト
   - [ ] 認証済みの場合の正常なアクセス
   - [ ] ページリロード時の認証状態の維持

3. **API通信**
   - [ ] バックエンドAPIへの認証付きリクエスト
   - [ ] CORS設定の確認

## デプロイ手順

1. 上記の変更をすべて実装
2. ローカルでテスト
   ```bash
   cd frontend
   npm run build
   npm run start
   ```
3. コミット・プッシュ
   ```bash
   git add .
   git commit -m "feat: Migrate to static export for Amplify deployment"
   git push origin main
   ```
4. Amplifyでの自動デプロイを確認

## ロールバック計画

もし問題が発生した場合：
1. `next.config.ts`から`output: 'export'`を削除
2. `amplify.yml`を元に戻す
3. `middleware.ts`を復活
4. SSRモードに戻す、またはVercelへの移行を検討

## 将来的な改善案

1. **Vercelへの移行**
   - Next.jsのフルサポート
   - SSR/ISRの完全なサポート
   - より簡単なデプロイプロセス

2. **AWS App Runnerの利用**
   - Dockerコンテナベースのデプロイ
   - SSRのフルサポート
   - カスタムドメインのサポート

3. **セキュリティの強化**
   - バックエンドでのより厳密な認証チェック
   - Rate limitingの実装
   - セキュリティヘッダーの追加