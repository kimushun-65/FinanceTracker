package handler

import (
	"net/http"

	"financetracker/internal/application/dto"
	"financetracker/internal/application/service"
	"financetracker/internal/infrastructure/auth0"

	"github.com/gin-gonic/gin"
)

// AuthHandler 認証関連のハンドラー
type AuthHandler struct {
	authService *service.AuthService
	auth0Client *auth0.Client
	middleware  *auth0.AuthMiddleware
	clientID    string
	redirectURI string
}

// NewAuthHandler 新しい認証ハンドラーを作成
func NewAuthHandler(authService *service.AuthService, auth0Client *auth0.Client, middleware *auth0.AuthMiddleware, clientID, redirectURI string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		auth0Client: auth0Client,
		middleware:  middleware,
		clientID:    clientID,
		redirectURI: redirectURI,
	}
}

// Login Auth0ログインページへリダイレクト
// @Summary Auth0ログイン
// @Description Auth0のログインページへリダイレクトします
// @Tags auth
// @Success 307 {string} string "Auth0ログインページへリダイレクト"
// @Router /auth/login [get]
func (h *AuthHandler) Login(c *gin.Context) {
	// CSRF防止用のstateパラメーターを生成
	state := generateRandomString(32) // 実装が必要

	// stateをセッションまたはキャッシュに保存
	c.SetCookie("auth_state", state, 600, "/", "", false, true)

	// Auth0ログインURLを構築
	loginURL := h.auth0Client.BuildLoginURL(state, h.redirectURI)

	c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

// Callback Auth0からのコールバックを処理
// @Summary Auth0コールバック
// @Description Auth0からの認証コールバックを処理し、ユーザー情報を保存します
// @Tags auth
// @Param code query string true "認証コード"
// @Param state query string true "CSRF防止用のstate"
// @Success 200 {object} map[string]interface{} "ユーザー情報とトークン"
// @Failure 400 {object} map[string]interface{} "リクエストエラー"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 500 {object} map[string]interface{} "内部エラー"
// @Router /auth/callback [get]
func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing authorization code",
		})
		return
	}

	// stateパラメーターを検証
	storedState, err := c.Cookie("auth_state")
	if err != nil || storedState != state {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid state parameter",
		})
		return
	}

	// stateクッキーをクリア
	c.SetCookie("auth_state", "", -1, "/", "", false, true)

	// 認可コードをトークンに交換
	// 通常はAuth0のトークンエンドポイントを呼び出す
	// 現時点では成功メッセージを返す
	c.JSON(http.StatusOK, gin.H{
		"message": "authentication successful",
		"code":    code,
	})
}

// Logout Auth0からログアウト
// @Summary ログアウト
// @Description ユーザーをログアウトし、Auth0のログアウトページへリダイレクトします
// @Tags auth
// @Success 307 {string} string "Auth0ログアウトページへリダイレクト"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// セッションクッキーをクリア
	// Auth0ログアウトURLを構築
	returnTo := "http://localhost:3000" // 環境変数から取得すべき
	logoutURL := h.auth0Client.BuildLogoutURL(returnTo)

	c.Redirect(http.StatusTemporaryRedirect, logoutURL)
}

// GetCurrentUser 現在の認証済みユーザー情報を取得
// @Summary 現在のユーザー情報取得
// @Description 認証済みユーザーの情報を取得します
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "ユーザー情報"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 500 {object} map[string]interface{} "内部エラー"
// @Router /auth/user [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	claims, exists := h.middleware.GetClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "no authenticated user",
		})
		return
	}

	// ヘッダーからアクセストークンを取得
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 7 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing authorization header",
		})
		return
	}

	accessToken := authHeader[7:] // "Bearer "プレフィックスを削除

	// Auth0からユーザー情報を取得
	userInfo, err := h.auth0Client.GetUserInfo(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user info",
		})
		return
	}

	// Auth0の型からDTOに変換
	userInfoDTO := &dto.UserInfo{
		Sub:           userInfo.Sub,
		Name:          userInfo.Name,
		Email:         userInfo.Email,
		EmailVerified: userInfo.EmailVerified,
		Picture:       userInfo.Picture,
	}
	
	tokenClaimsDTO := &dto.TokenClaims{
		Subject: claims.Subject,
		Roles:   claims.Roles,
	}

	// データベースとユーザーを同期
	user, err := h.authService.SyncUser(c.Request.Context(), userInfoDTO, tokenClaimsDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to sync user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// RefreshToken トークンをリフレッシュ
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// このエンドポイントはトークンリフレッシュを処理
	// 実装はトークン戦略に依存
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

// generateRandomString stateパラメーター用のランダム文字列を生成するヘルパー関数
// generateRandomString ランダム文字列を生成
func generateRandomString(length int) string {
	// これは簡略化されたバージョン。本番環境ではcrypto/randを使用
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
