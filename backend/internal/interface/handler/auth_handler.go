package handler

import (
	"net/http"

	"financetracker/internal/application/service"
	"financetracker/internal/infrastructure/auth0"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	auth0Client *auth0.Client
	middleware  *auth0.AuthMiddleware
	clientID    string
	redirectURI string
}

func NewAuthHandler(authService *service.AuthService, auth0Client *auth0.Client, middleware *auth0.AuthMiddleware, clientID, redirectURI string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		auth0Client: auth0Client,
		middleware:  middleware,
		clientID:    clientID,
		redirectURI: redirectURI,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	// Generate state parameter for CSRF protection
	state := generateRandomString(32) // You'll need to implement this

	// Store state in session or cache
	c.SetCookie("auth_state", state, 600, "/", "", false, true)

	// Build Auth0 login URL
	loginURL := h.auth0Client.BuildLoginURL(state, h.redirectURI)

	c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing authorization code",
		})
		return
	}

	// Verify state parameter
	storedState, err := c.Cookie("auth_state")
	if err != nil || storedState != state {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid state parameter",
		})
		return
	}

	// Clear state cookie
	c.SetCookie("auth_state", "", -1, "/", "", false, true)

	// Exchange code for tokens
	// This would typically involve calling Auth0's token endpoint
	// For now, we'll return a success message
	c.JSON(http.StatusOK, gin.H{
		"message": "authentication successful",
		"code":    code,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear any session cookies
	// Build Auth0 logout URL
	returnTo := "http://localhost:3000" // This should come from env
	logoutURL := h.auth0Client.BuildLogoutURL(returnTo)

	c.Redirect(http.StatusTemporaryRedirect, logoutURL)
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	claims, exists := h.middleware.GetClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "no authenticated user",
		})
		return
	}

	// Get access token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 7 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing authorization header",
		})
		return
	}

	accessToken := authHeader[7:] // Remove "Bearer " prefix

	// Get user info from Auth0
	userInfo, err := h.auth0Client.GetUserInfo(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user info",
		})
		return
	}

	// Sync user with database
	user, err := h.authService.SyncUser(c.Request.Context(), userInfo, claims)
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

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// This endpoint would handle token refresh
	// Implementation depends on your token strategy
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

// Helper function to generate random strings for state parameter
func generateRandomString(length int) string {
	// This is a simplified version. In production, use crypto/rand
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
