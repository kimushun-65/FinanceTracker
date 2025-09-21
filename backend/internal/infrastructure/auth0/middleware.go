package auth0

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	ClaimsContextKey contextKey = "auth0_claims"
	UserIDContextKey contextKey = "user_id"
)

type AuthMiddleware struct {
	client *Client
}

func NewAuthMiddleware(client *Client) *AuthMiddleware {
	return &AuthMiddleware{
		client: client,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c.Request)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		claims, err := m.client.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("invalid token: %v", err),
			})
			c.Abort()
			return
		}

		// Set claims in context
		c.Set(string(ClaimsContextKey), claims)

		// Extract user ID
		userID := claims.Subject
		if claims.UserID != "" {
			userID = claims.UserID
		}
		c.Set(string(UserIDContextKey), userID)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := m.GetClaims(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "no claims found",
			})
			c.Abort()
			return
		}

		if !m.hasScope(claims.Scope, scope) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "insufficient scope",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := m.GetClaims(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "no claims found",
			})
			c.Abort()
			return
		}

		hasRole := false
		for _, r := range claims.Roles {
			if r == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "insufficient role",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func (m *AuthMiddleware) hasScope(grantedScope, requiredScope string) bool {
	if grantedScope == "" {
		return false
	}

	scopes := strings.Split(grantedScope, " ")
	for _, scope := range scopes {
		if scope == requiredScope {
			return true
		}
	}

	return false
}

func (m *AuthMiddleware) GetClaims(c *gin.Context) (*TokenClaims, bool) {
	claims, exists := c.Get(string(ClaimsContextKey))
	if !exists {
		return nil, false
	}

	tokenClaims, ok := claims.(*TokenClaims)
	if !ok {
		return nil, false
	}

	return tokenClaims, true
}

func (m *AuthMiddleware) GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(string(UserIDContextKey))
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	if !ok {
		return "", false
	}

	return id, true
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID := ctx.Value(UserIDContextKey)
	if userID == nil {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

func SetUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}
