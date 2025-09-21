package router

import (
	"github.com/gin-gonic/gin"
)

// AuthRouter 認証関連のルートを管理
type AuthRouter struct {
	handlers *Handlers
}

// NewAuthRouter 新しい認証ルーターを作成
func NewAuthRouter(handlers *Handlers) *AuthRouter {
	return &AuthRouter{
		handlers: handlers,
	}
}

// RegisterRoutes 認証関連のルートを登録
func (ar *AuthRouter) RegisterRoutes(group *gin.RouterGroup) {
	// 公開ルート（認証不要）
	public := group.Group("/auth")
	{
		if ar.handlers != nil && ar.handlers.AuthHandler != nil {
			public.GET("/login", ar.handlers.AuthHandler.Login)
			public.GET("/callback", ar.handlers.AuthHandler.Callback)
			public.POST("/logout", ar.handlers.AuthHandler.Logout)
		}
	}
}

// RegisterProtectedRoutes 認証が必要なルートを登録
func (ar *AuthRouter) RegisterProtectedRoutes(group *gin.RouterGroup) {
	if ar.handlers != nil && ar.handlers.AuthHandler != nil {
		group.GET("/auth/user", ar.handlers.AuthHandler.GetCurrentUser)
	}
}
