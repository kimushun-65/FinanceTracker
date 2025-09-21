package router

import "github.com/gin-gonic/gin"

// DomainRouter ドメイン別ルーターのインターフェース
type DomainRouter interface {
	// RegisterRoutes 公開ルートを登録
	RegisterRoutes(group *gin.RouterGroup)

	// RegisterProtectedRoutes 認証が必要なルートを登録
	RegisterProtectedRoutes(group *gin.RouterGroup)
}

// ルートの一覧（将来の実装用）
const (
	// 認証関連
	RouteAuthLogin       = "/auth/login"
	RouteAuthCallback    = "/auth/callback"
	RouteAuthLogout      = "/auth/logout"
	RouteAuthCurrentUser = "/auth/user"

	// ユーザー関連
	RouteUserMe       = "/users/me"
	RouteUserUpdateMe = "/users/me"

	// アカウント関連
	RouteAccountList   = "/accounts"
	RouteAccountCreate = "/accounts"
	RouteAccountGet    = "/accounts/:id"
	RouteAccountUpdate = "/accounts/:id"
	RouteAccountDelete = "/accounts/:id"

	// トランザクション関連
	RouteTransactionList   = "/transactions"
	RouteTransactionCreate = "/transactions"
	RouteTransactionGet    = "/transactions/:id"
	RouteTransactionUpdate = "/transactions/:id"
	RouteTransactionDelete = "/transactions/:id"
)
