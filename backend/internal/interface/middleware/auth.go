// Package middleware provides HTTP middleware implementations for the FinanceTracker API.
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"financetracker/pkg/config"
	"financetracker/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWKSResponse represents the response from Auth0 JWKS endpoint.
type JWKSResponse struct {
	Keys []json.RawMessage `json:"keys"`
}

// Auth0Claims represents the JWT claims from Auth0.
type Auth0Claims struct {
	jwt.RegisteredClaims
	Scope       string   `json:"scope"`
	Permissions []string `json:"permissions"`
}

// Auth returns a middleware that validates JWT tokens from Auth0.
func Auth(cfg *config.Config) gin.HandlerFunc {
	// Construct the JWKS URL
	auth0Domain := cfg.Auth0Domain
	if !strings.HasPrefix(auth0Domain, "https://") {
		auth0Domain = "https://" + auth0Domain
	}
	jwksURL := auth0Domain + "/.well-known/jwks.json"

	// Parse Auth0 audience
	expectedAudience := cfg.Auth0Audience
	expectedIssuer := auth0Domain + "/"

	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if err := c.Error(errors.NewUnauthorizedError("Missing authorization header")); err != nil {
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			if err := c.Error(errors.NewUnauthorizedError("Invalid authorization header format")); err != nil {
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse token with claims
		token, err := jwt.ParseWithClaims(tokenString, &Auth0Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Verify signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.NewUnauthorizedError("Unexpected signing method")
			}

			// Get the key ID from token header
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, errors.NewUnauthorizedError("Missing key ID")
			}

			// Fetch JWKS (in production, this should be cached)
			publicKey, err := fetchPublicKey(jwksURL, kid)
			if err != nil {
				return nil, err
			}

			return publicKey, nil
		})

		if err != nil || !token.Valid {
			if err := c.Error(errors.NewUnauthorizedError("Invalid token")); err != nil {
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		// Extract and validate claims
		claims, ok := token.Claims.(*Auth0Claims)
		if !ok {
			if err := c.Error(errors.NewUnauthorizedError("Invalid token claims")); err != nil {
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		// Verify issuer
		if claims.Issuer != expectedIssuer {
			if err := c.Error(errors.NewUnauthorizedError("Invalid token issuer")); err != nil {
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		// Verify audience
		audValid := false
		for _, aud := range claims.Audience {
			if aud == expectedAudience {
				audValid = true
				break
			}
		}
		if !audValid {
			if err := c.Error(errors.NewUnauthorizedError("Invalid token audience")); err != nil {
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		// Check token expiration
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			if err := c.Error(errors.NewUnauthorizedError("Token has expired")); err != nil {
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		// Extract user ID from subject
		userID := claims.Subject
		if userID == "" {
			if err := c.Error(errors.NewUnauthorizedError("Missing user ID in token")); err != nil {
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("UserID", userID)
		c.Set("Claims", claims)
		c.Set("Permissions", claims.Permissions)

		c.Next()
	}
}

// fetchPublicKey fetches the public key from Auth0 JWKS endpoint.
// In production, this should be cached to avoid repeated HTTP requests.
func fetchPublicKey(jwksURL, kid string) (interface{}, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksURL, http.NoBody)
	if err != nil {
		return nil, errors.NewInternalError("Failed to create JWKS request", err)
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.NewExternalServiceError("Auth0", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error but don't return it
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewExternalServiceError("Auth0", nil)
	}

	// Parse JWKS response
	var jwks JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, errors.NewInternalError("Failed to parse JWKS", err)
	}

	// Find the key with matching kid
	for _, key := range jwks.Keys {
		var keyData map[string]interface{}
		if err := json.Unmarshal(key, &keyData); err != nil {
			continue
		}

		if keyData["kid"] == kid {
			// Parse the key (simplified - in production use a proper JWT library)
			// This is a placeholder - actual implementation would parse RSA key from JWK
			return parseRSAPublicKeyFromJWK(keyData)
		}
	}

	return nil, errors.NewUnauthorizedError("Public key not found")
}

// parseRSAPublicKeyFromJWK parses RSA public key from JWK.
// This is a simplified placeholder - use a proper JWT library in production.
func parseRSAPublicKeyFromJWK(_ map[string]interface{}) (interface{}, error) {
	// TODO: Implement proper JWK to RSA public key conversion
	// This would typically use a library like github.com/lestrrat-go/jwx
	return nil, errors.NewNotImplementedError("JWK parsing")
}

// RequirePermission returns a middleware that checks for specific permissions.
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get permissions from context
		permissions, exists := c.Get("Permissions")
		if !exists {
			if err := c.Error(errors.NewForbiddenError("No permissions found")); err != nil {
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		// Check if user has the required permission
		userPermissions, ok := permissions.([]string)
		if !ok {
			if err := c.Error(errors.NewInternalError("Invalid permissions format", nil)); err != nil {
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		hasPermission := false
		for _, p := range userPermissions {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			if err := c.Error(errors.NewForbiddenError("Insufficient permissions")); err != nil {
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth is a middleware that validates JWT tokens if present but doesn't require them.
func OptionalAuth(cfg *config.Config) gin.HandlerFunc {
	authMiddleware := Auth(cfg)

	return func(c *gin.Context) {
		// Check if Authorization header exists
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No auth header, continue without authentication
			c.Next()
			return
		}

		// If auth header exists, validate it
		authMiddleware(c)
	}
}
