# Auth0 統合ガイド（Google認証）

## 概要

FinSightアプリケーションでは、ユーザー認証にAuth0を通じたGoogle認証（ソーシャルログイン）のみを使用します。パスワードベースの認証は提供しません。このドキュメントでは、Google認証に特化したAuth0の統合方法と実装の詳細について説明します。

## アーキテクチャ

### 認証フロー

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │     │    Auth0    │     │   API       │
│   (SPA)     │     │             │     │   Server    │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                     │
       │ 1. ログイン要求   │                     │
       │──────────────────>│                     │
       │                   │                     │
       │ 2. Google認証へ   │                     │
       │<──────────────────│ リダイレクト       │
       │                   │                     │
       │ 3. Googleで認証   │                     │
       │──────────────────>│ (Google OAuth2)    │
       │                   │                     │
       │ 4. トークン返却   │                     │
       │<──────────────────│                     │
       │ (ID/Access Token) │                     │
       │                   │                     │
       │ 5. API呼び出し    │                     │
       │───────────────────┼────────────────────>│
       │ (Bearer Token)    │                     │
       │                   │ 6. トークン検証    │
       │                   │<────────────────────│
       │                   │ (JWKS)             │
       │                   │                     │
       │ 7. レスポンス     │                     │
       │<──────────────────┼─────────────────────│
       │                   │                     │
```

## フロントエンド実装（Next.js）

### 1. Auth0 SDK のインストール

```bash
npm install @auth0/nextjs-auth0
```

### 2. 環境変数の設定

```env
# .env.local
AUTH0_SECRET='use [openssl rand -hex 32] to generate a 32 bytes value'
AUTH0_BASE_URL='https://app.finsight.com'
AUTH0_ISSUER_BASE_URL='https://your-tenant.auth0.com'
AUTH0_CLIENT_ID='your_client_id'
AUTH0_CLIENT_SECRET='your_client_secret'
AUTH0_AUDIENCE='https://api.finsight.com'
AUTH0_SCOPE='openid profile email read:transactions write:transactions'
AUTH0_CONNECTION='google-oauth2'  # Google認証のみ使用
```

### 3. Auth0 Provider の設定

```typescript
// pages/_app.tsx
import { UserProvider } from '@auth0/nextjs-auth0/client';

export default function App({ Component, pageProps }) {
  return (
    <UserProvider>
      <Component {...pageProps} />
    </UserProvider>
  );
}
```

### 4. API Routes の設定

```typescript
// pages/api/auth/[...auth0].ts
import { handleAuth, handleCallback } from '@auth0/nextjs-auth0';

export default handleAuth({
  async callback(req, res) {
    try {
      await handleCallback(req, res, {
        afterCallback: async (req, res, session) => {
          // ユーザー情報をバックエンドに同期
          const response = await fetch(`${process.env.API_BASE_URL}/auth/sync`, {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${session.accessToken}`,
              'Content-Type': 'application/json'
            }
          });

          if (!response.ok) {
            throw new Error('Failed to sync user');
          }

          return session;
        }
      });
    } catch (error) {
      res.redirect('/error');
    }
  }
});
```

### 5. 認証状態の利用とGoogle専用ログイン

```typescript
// components/NavBar.tsx
import { useUser } from '@auth0/nextjs-auth0/client';

export default function NavBar() {
  const { user, error, isLoading } = useUser();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>{error.message}</div>;

  const handleLogin = () => {
    // Google認証のみを強制
    window.location.href = '/api/auth/login?connection=google-oauth2';
  };

  return (
    <nav>
      {user ? (
        <>
          <img 
            src={user.picture} 
            alt={user.name}
            className="w-8 h-8 rounded-full"
          />
          <span>Welcome, {user.name}</span>
          <a href="/api/auth/logout">Logout</a>
        </>
      ) : (
        <button onClick={handleLogin}>
          <img src="/google-logo.svg" alt="Google" className="w-5 h-5 mr-2" />
          Sign in with Google
        </button>
      )}
    </nav>
  );
}
```

### 6. ログインボタンコンポーネント

```typescript
// components/GoogleLoginButton.tsx
import { FcGoogle } from 'react-icons/fc';

export default function GoogleLoginButton() {
  const handleGoogleLogin = () => {
    window.location.href = '/api/auth/login?connection=google-oauth2';
  };

  return (
    <button
      onClick={handleGoogleLogin}
      className="flex items-center justify-center w-full px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
    >
      <FcGoogle className="w-5 h-5 mr-2" />
      Continue with Google
    </button>
  );
}
```

### 7. API呼び出しフック

```typescript
// hooks/useApi.ts
import { useEffect, useState } from 'react';
import { useUser } from '@auth0/nextjs-auth0/client';

export function useApi(url: string) {
  const { user } = useUser();
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!user) return;

    const fetchData = async () => {
      try {
        const response = await fetch(`/api/auth/token`);
        const { accessToken } = await response.json();

        const apiResponse = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}${url}`, {
          headers: {
            'Authorization': `Bearer ${accessToken}`,
            'Content-Type': 'application/json'
          }
        });

        if (!apiResponse.ok) {
          throw new Error('API request failed');
        }

        const data = await apiResponse.json();
        setData(data);
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [user, url]);

  return { data, loading, error };
}
```

## バックエンド実装（Go + Gin）

### 1. JWT検証ミドルウェア

```go
// middleware/auth.go
package middleware

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/auth0/go-jwt-middleware/v2"
    "github.com/auth0/go-jwt-middleware/v2/jwks"
    "github.com/auth0/go-jwt-middleware/v2/validator"
    "github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
    issuerURL := fmt.Sprintf("https://%s/", os.Getenv("AUTH0_DOMAIN"))
    audience := os.Getenv("AUTH0_AUDIENCE")

    provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

    jwtValidator, err := validator.New(
        provider.KeyFunc,
        validator.RS256,
        issuerURL,
        []string{audience},
        validator.WithAllowedClockSkew(5*time.Minute),
    )

    if err != nil {
        log.Fatalf("Failed to create validator: %v", err)
    }

    return func(c *gin.Context) {
        authorization := c.GetHeader("Authorization")
        if authorization == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code": "UNAUTHORIZED",
                    "message": "Authorization header is required",
                },
            })
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authorization, "Bearer ")
        if token == authorization {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code": "UNAUTHORIZED",
                    "message": "Bearer token is required",
                },
            })
            c.Abort()
            return
        }

        tokenClaims, err := jwtValidator.ValidateToken(c.Request.Context(), token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code": "UNAUTHORIZED",
                    "message": "Invalid token",
                },
            })
            c.Abort()
            return
        }

        claims := tokenClaims.(*validator.ValidatedClaims)
        // Google認証の場合、SubjectはGoogle OAuth2のIDを含む
        // 例: "google-oauth2|115744780651980123456"
        c.Set("user_id", claims.RegisteredClaims.Subject)
        c.Set("email", claims.CustomClaims["email"])
        c.Set("provider", "google-oauth2")
        c.Set("permissions", claims.CustomClaims["permissions"])
        
        c.Next()
    }
}
```

### 2. パーミッションチェック

```go
// middleware/permission.go
package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        permissions, exists := c.Get("permissions")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{
                "error": gin.H{
                    "code": "FORBIDDEN",
                    "message": "No permissions found",
                },
            })
            c.Abort()
            return
        }

        perms, ok := permissions.([]interface{})
        if !ok {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": gin.H{
                    "code": "INTERNAL_SERVER_ERROR",
                    "message": "Invalid permissions format",
                },
            })
            c.Abort()
            return
        }

        hasPermission := false
        for _, p := range perms {
            if pStr, ok := p.(string); ok && pStr == permission {
                hasPermission = true
                break
            }
        }

        if !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{
                "error": gin.H{
                    "code": "FORBIDDEN",
                    "message": fmt.Sprintf("Permission '%s' is required", permission),
                },
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 3. ルーティング設定

```go
// routes/routes.go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/your-org/finsight/middleware"
    "github.com/your-org/finsight/handlers"
)

func SetupRoutes(router *gin.Engine) {
    v1 := router.Group("/v1")

    // Auth0コールバック（認証不要）
    auth := v1.Group("/auth")
    {
        auth.POST("/callback", handlers.Auth0Callback)
        auth.POST("/sync", middleware.AuthRequired(), handlers.SyncUser)
    }

    // ユーザー管理（認証必須）
    users := v1.Group("/users")
    users.Use(middleware.AuthRequired())
    {
        users.GET("/me", handlers.GetProfile)
        users.PUT("/me", handlers.UpdateProfile)
        users.POST("/me/request-password-reset", handlers.RequestPasswordReset)
        users.DELETE("/me", handlers.DeleteAccount)
    }

    // 取引管理（認証＋権限必須）
    transactions := v1.Group("/transactions")
    transactions.Use(middleware.AuthRequired())
    {
        transactions.GET("", 
            middleware.RequirePermission("read:transactions"), 
            handlers.ListTransactions)
        transactions.POST("", 
            middleware.RequirePermission("write:transactions"), 
            handlers.CreateTransaction)
        transactions.DELETE("/:id", 
            middleware.RequirePermission("delete:transactions"), 
            handlers.DeleteTransaction)
    }

    // 以下、他のエンドポイントも同様に設定
}
```

## Auth0 管理画面設定

### 1. Application設定

```
Application Type: Single Page Application

Allowed Callback URLs:
- https://app.finsight.com/api/auth/callback
- http://localhost:3000/api/auth/callback (開発環境)

Allowed Logout URLs:
- https://app.finsight.com
- http://localhost:3000 (開発環境)

Allowed Web Origins:
- https://app.finsight.com
- http://localhost:3000 (開発環境)

Refresh Token Rotation: Enabled
Refresh Token Expiration: 2592000 (30日)

// 接続設定
Connections:
- Google OAuth2: Enabled ✓
- Database: Disabled ✗
- Other Social: Disabled ✗
```

### 2. API設定

```
API Identifier: https://api.finsight.com
Signing Algorithm: RS256
Token Expiration: 86400 (24時間)
Allow Skipping User Consent: Enabled
```

### 3. Google OAuth2 接続設定

Auth0ダッシュボード → Authentication → Social → Google:

```
Client ID: [Google Cloud ConsoleのOAuth2クライアントID]
Client Secret: [Google Cloud ConsoleのOAuth2クライアントシークレット]

許可するドメイン:
- gmail.com
- [会社のG Suiteドメイン（オプション）]

Attribute Mappings:
- email → email
- email_verified → email_verified
- name → name
- given_name → given_name
- family_name → family_name
- picture → picture
- locale → locale
```

### 4. Rules/Actions設定

#### Google認証ユーザーのメタデータ追加

```javascript
// Add Google User Metadata to Tokens
exports.onExecutePostLogin = async (event, api) => {
  const namespace = 'https://finsight.com';
  
  // Google認証のみを許可
  if (event.connection.strategy !== 'google-oauth2') {
    api.access.deny('Only Google authentication is allowed');
    return;
  }
  
  if (event.authorization) {
    // IDトークンにカスタムクレーム追加
    api.idToken.setCustomClaim(`${namespace}/user_id`, event.user.user_id);
    api.idToken.setCustomClaim(`${namespace}/email`, event.user.email);
    api.idToken.setCustomClaim(`${namespace}/provider`, 'google-oauth2');
    api.idToken.setCustomClaim(`${namespace}/picture`, event.user.picture);
    
    // アクセストークンにパーミッション追加
    api.accessToken.setCustomClaim(`${namespace}/permissions`, event.authorization.roles);
    api.accessToken.setCustomClaim(`${namespace}/email`, event.user.email);
  }
  
  // 初回ログイン時の処理
  if (event.stats.logins_count === 1) {
    // ウェルカムメール送信などの処理
    api.user.setAppMetadata('onboarding_completed', false);
  }
};
```

#### 特定ドメイン制限（オプション）

```javascript
// Restrict to specific email domains
exports.onExecutePostLogin = async (event, api) => {
  const allowedDomains = ['gmail.com', 'yourcompany.com'];
  const emailDomain = event.user.email.split('@')[1];
  
  if (!allowedDomains.includes(emailDomain)) {
    api.access.deny(`Access denied. Email domain ${emailDomain} is not allowed.`);
  }
};
```

## セキュリティベストプラクティス

### 1. トークンの保存

- アクセストークンは httpOnly Cookie に保存
- リフレッシュトークンは Secure Cookie に保存
- LocalStorage や SessionStorage は使用しない

### 2. CORS設定

```go
// main.go
config := cors.DefaultConfig()
config.AllowOrigins = []string{
    "https://app.finsight.com",
    "http://localhost:3000", // 開発環境のみ
}
config.AllowHeaders = []string{
    "Origin",
    "Content-Type",
    "Accept",
    "Authorization",
}
config.AllowCredentials = true

router.Use(cors.New(config))
```

### 3. レート制限

```go
// middleware/ratelimit.go
limiter := tollbooth.NewLimiter(100, &limiter.ExpirableOptions{
    DefaultExpirationTTL: time.Hour,
})

router.Use(LimitHandler(limiter))
```

## 開発環境セットアップ

### 1. Auth0テナント作成

1. https://auth0.com でアカウント作成
2. 新規テナント作成（dev-finsight）
3. Application と API を作成

### 2. ローカル開発

```bash
# .env.development
AUTH0_DOMAIN=dev-finsight.auth0.com
AUTH0_CLIENT_ID=your_dev_client_id
AUTH0_CLIENT_SECRET=your_dev_client_secret
AUTH0_AUDIENCE=http://localhost:8080
```

### 3. テスト用ユーザー作成

Google認証を使用するため、テストには実際のGoogleアカウントを使用：
1. 開発用のGoogleアカウントでログイン
2. Auth0ダッシュボードでユーザーが自動作成されることを確認
3. Users & Roles → Users から該当ユーザーにロール・パーミッションを付与

## トラブルシューティング

### よくある問題と解決方法

1. **「Invalid token」エラー**
   - Audienceが一致しているか確認
   - トークンの有効期限を確認
   - JWKSエンドポイントへのアクセスを確認

2. **CORSエラー**
   - Auth0のAllowed Web Originsを確認
   - APIサーバーのCORS設定を確認

3. **Google認証が失敗する**
   - Google OAuth2の接続が有効になっているか確認
   - Google Cloud ConsoleでのOAuth2設定を確認
   - リダイレクトURIがGoogle側で許可されているか確認

4. **ユーザー情報が取得できない**
   - openid profile emailスコープが含まれているか確認
   - Google OAuth2のAttribute Mappingsを確認

5. **「Only Google authentication is allowed」エラー**
   - ログインURLに`connection=google-oauth2`パラメータが含まれているか確認
   - Auth0 Actionsで他の認証方法をブロックしているか確認

## Google Cloud Console設定

### OAuth 2.0 クライアントID作成

1. [Google Cloud Console](https://console.cloud.google.com/)にアクセス
2. APIs & Services → Credentials
3. Create Credentials → OAuth client ID
4. Application type: Web application
5. Authorized redirect URIs:
   ```
   https://your-tenant.auth0.com/login/callback
   https://your-tenant.us.auth0.com/login/callback
   ```
6. クライアントIDとシークレットをAuth0に設定

### 必要なAPIの有効化

- Google+ API（非推奨だが一部機能で必要な場合あり）
- Google Identity Toolkit API
- People API（プロフィール情報取得用）