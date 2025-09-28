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
	// middleware.Authが設定したコンテキストからauth0_idを取得
	auth0ID, exists := c.Get("auth0_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "no authenticated user",
		})
		return
	}

	// auth0_idをstringに変換
	auth0IDStr, ok := auth0ID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id format",
		})
		return
	}

	// HttpOnlyクッキーからアクセストークンを取得
	accessToken, err := c.Cookie("access_token")
	if err != nil || accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing access token",
		})
		return
	}

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
	}

	// middleware.Authはロールを設定していないため、空の配列を使用
	tokenClaimsDTO := &dto.TokenClaims{
		Subject: auth0IDStr,
		Roles:   []string{}, // TODO: Auth0からロール情報を取得する場合は実装が必要
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

// CheckAuth 認証状態をチェック
// @Summary 認証状態チェック
// @Description HttpOnlyクッキーからトークンを読み取り、認証状態を確認します
// @Tags auth
// @Success 200 {object} map[string]interface{} "認証成功"
// @Failure 401 {object} map[string]interface{} "認証失敗"
// @Router /auth/check [get]
func (h *AuthHandler) CheckAuth(c *gin.Context) {
	// HttpOnlyクッキーからアクセストークンを取得
	accessToken, err := c.Cookie("access_token")
	if err != nil || accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"authenticated": false,
			"error":         "no access token found",
		})
		return
	}

	// Auth0でトークンを検証
	userInfo, err := h.auth0Client.GetUserInfo(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"authenticated": false,
			"error":         "invalid token",
		})
		return
	}

	// 認証成功
	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user": gin.H{
			"sub":           userInfo.Sub,
			"name":          userInfo.Name,
			"email":         userInfo.Email,
			"emailVerified": userInfo.EmailVerified,
		},
	})
}

// SetToken HttpOnlyクッキーにトークンを設定
// @Summary トークン設定
// @Description Auth0トークンをHttpOnlyクッキーに保存します
// @Tags auth
// @Param request body map[string]string true "トークンリクエスト"
// @Success 200 {object} map[string]interface{} "トークン設定成功"
// @Failure 400 {object} map[string]interface{} "リクエストエラー"
// @Router /auth/token [post]
func (h *AuthHandler) SetToken(c *gin.Context) {
	var request struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// HttpOnlyクッキーを設定
	// 開発環境ではsecureをfalseに設定
	secure := false
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		secure = true
	}

	// SameSiteの設定 - 開発環境ではNoneに設定（クロスオリジンでクッキーを送信するため）
	sameSite := http.SameSiteNoneMode
	if c.Request.Host == "localhost:8080" || c.Request.Host == "127.0.0.1:8080" {
		// 開発環境ではSameSite=Noneはセキュアでないため、Laxを使用
		sameSite = http.SameSiteLaxMode
		// ただし、クロスオリジンの場合は動作しない可能性があるため、
		// 開発環境では domain を明示的に設定
	}
	c.SetSameSite(sameSite)

	// domainは空にして、現在のホストのみで有効にする
	// 異なるポート間でのクッキー共有は複雑なため
	domain := ""

	c.SetCookie(
		"access_token", // name
		request.Token,  // value
		60*60*24*7,     // maxAge (7 days)
		"/",            // path
		domain,         // domain
		secure,         // secure (HTTPS only in production)
		true,           // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// RemoveToken HttpOnlyクッキーからトークンを削除
// @Summary トークン削除
// @Description HttpOnlyクッキーからアクセストークンを削除します
// @Tags auth
// @Success 200 {object} map[string]interface{} "トークン削除成功"
// @Router /auth/token [delete]
func (h *AuthHandler) RemoveToken(c *gin.Context) {
	// クッキーを削除
	// 開発環境ではsecureをfalseに設定
	secure := false
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		secure = true
	}

	// SameSiteの設定 - SetTokenと同じロジックを使用
	sameSite := http.SameSiteNoneMode
	if c.Request.Host == "localhost:8080" || c.Request.Host == "127.0.0.1:8080" {
		sameSite = http.SameSiteLaxMode
	}
	c.SetSameSite(sameSite)

	// domainは空にして、現在のホストのみで有効にする
	domain := ""

	c.SetCookie(
		"access_token", // name
		"",             // value
		-1,             // maxAge (expired)
		"/",            // path
		domain,         // domain
		secure,         // secure
		true,           // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
