// Package middleware provides HTTP middleware implementations for the FinanceTracker API.
package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
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

// authConfig holds configuration for Auth middleware
type authConfig struct {
	jwksURL          string
	expectedAudience string
	expectedIssuer   string
}

// newAuthConfig creates a new authConfig from the application config
func newAuthConfig(cfg *config.Config) *authConfig {
	auth0Domain := cfg.Auth0Domain
	if !strings.HasPrefix(auth0Domain, "https://") {
		auth0Domain = "https://" + auth0Domain
	}

	return &authConfig{
		jwksURL:          auth0Domain + "/.well-known/jwks.json",
		expectedAudience: cfg.Auth0Audience,
		expectedIssuer:   auth0Domain + "/",
	}
}

// Auth returns a middleware that validates JWT tokens from Auth0.
func Auth(cfg *config.Config) gin.HandlerFunc {
	authCfg := newAuthConfig(cfg)

	return func(c *gin.Context) {
		if err := authenticateRequest(c, authCfg); err != nil {
			handleAuthError(c, err)
			return
		}
		c.Next()
	}
}

// authenticateRequest handles the authentication flow for a request
func authenticateRequest(c *gin.Context, cfg *authConfig) error {
	// Extract token
	tokenString, err := extractTokenFromHeader(c)
	if err != nil {
		return err
	}

	// Parse and validate token
	token, err := parseToken(tokenString, cfg.jwksURL)
	if err != nil {
		return errors.NewUnauthorizedError("Invalid token")
	}

	if !token.Valid {
		return errors.NewUnauthorizedError("Invalid token")
	}

	// Validate claims
	claims, err := validateClaims(token, cfg.expectedIssuer, cfg.expectedAudience)
	if err != nil {
		return err
	}

	// Set user information in context
	setUserContext(c, claims)
	return nil
}

// setUserContext sets user information in the Gin context
func setUserContext(c *gin.Context, claims *Auth0Claims) {
	c.Set("UserID", claims.Subject)
	c.Set("auth0ID", claims.Subject) // UserHandlerが期待するキー
	c.Set("Claims", claims)
	c.Set("Permissions", claims.Permissions)
}

// extractTokenFromHeader extracts the JWT token from the Authorization header
func extractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.NewUnauthorizedError("Missing authorization header")
	}

	// Check Bearer prefix
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.NewUnauthorizedError("Invalid authorization header format")
	}

	return parts[1], nil
}

// parseToken parses the JWT token and fetches the public key
func parseToken(tokenString, jwksURL string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &Auth0Claims{}, func(token *jwt.Token) (interface{}, error) {
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
		return fetchPublicKey(jwksURL, kid)
	})
}

// validateClaims validates the token claims
func validateClaims(token *jwt.Token, expectedIssuer, expectedAudience string) (*Auth0Claims, error) {
	// Extract claims
	claims, ok := token.Claims.(*Auth0Claims)
	if !ok {
		return nil, errors.NewUnauthorizedError("Invalid token claims")
	}

	// Verify issuer
	if claims.Issuer != expectedIssuer {
		return nil, errors.NewUnauthorizedError("Invalid token issuer")
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
		return nil, errors.NewUnauthorizedError("Invalid token audience")
	}

	// Check token expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.NewUnauthorizedError("Token has expired")
	}

	// Verify user ID exists
	if claims.Subject == "" {
		return nil, errors.NewUnauthorizedError("Missing user ID in token")
	}

	return claims, nil
}

// handleAuthError handles authentication errors
func handleAuthError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		_ = c.Error(appErr)
	} else {
		_ = c.Error(errors.NewUnauthorizedError(err.Error()))
	}
	c.Abort()
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
		_ = resp.Body.Close()
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
func parseRSAPublicKeyFromJWK(jwk map[string]interface{}) (interface{}, error) {
	// Extract the key type
	kty, ok := jwk["kty"].(string)
	if !ok || kty != "RSA" {
		return nil, errors.NewUnauthorizedError("Invalid key type")
	}

	// Extract modulus (n) and exponent (e)
	nStr, ok := jwk["n"].(string)
	if !ok {
		return nil, errors.NewUnauthorizedError("Missing modulus")
	}

	eStr, ok := jwk["e"].(string)
	if !ok {
		return nil, errors.NewUnauthorizedError("Missing exponent")
	}

	// Decode base64url encoded values
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid modulus encoding")
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid exponent encoding")
	}

	// Convert to big integers
	n := new(big.Int).SetBytes(nBytes)

	// Convert exponent bytes to int
	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	// Create RSA public key
	pubKey := &rsa.PublicKey{
		N: n,
		E: e,
	}

	return pubKey, nil
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
