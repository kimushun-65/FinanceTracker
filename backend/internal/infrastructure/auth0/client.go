package auth0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/square/go-jose/v3"
)

type Client struct {
	domain     string
	audience   string
	jwks       *jose.JSONWebKeySet
	jwksMutex  sync.RWMutex
	jwksExpiry time.Time
}

func NewClient(domain, audience string) *Client {
	return &Client{
		domain:   domain,
		audience: audience,
	}
}

func (c *Client) GetJWKS(ctx context.Context) (*jose.JSONWebKeySet, error) {
	c.jwksMutex.RLock()
	if c.jwks != nil && time.Now().Before(c.jwksExpiry) {
		defer c.jwksMutex.RUnlock()
		return c.jwks, nil
	}
	c.jwksMutex.RUnlock()

	c.jwksMutex.Lock()
	defer c.jwksMutex.Unlock()

	// Double-check after acquiring write lock
	if c.jwks != nil && time.Now().Before(c.jwksExpiry) {
		return c.jwks, nil
	}

	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", c.domain)
	req, err := http.NewRequestWithContext(ctx, "GET", jwksURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var jwks jose.JSONWebKeySet
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	c.jwks = &jwks
	c.jwksExpiry = time.Now().Add(1 * time.Hour)

	return &jwks, nil
}

func (c *Client) GetDomain() string {
	return c.domain
}

func (c *Client) GetAudience() string {
	return c.audience
}

func (c *Client) GetIssuer() string {
	return fmt.Sprintf("https://%s/", c.domain)
}

type TokenClaims struct {
	Issuer    string   `json:"iss"`
	Subject   string   `json:"sub"`
	Audience  []string `json:"aud"`
	ExpiresAt int64    `json:"exp"`
	IssuedAt  int64    `json:"iat"`
	Scope     string   `json:"scope,omitempty"`
	Azp       string   `json:"azp,omitempty"`

	// Custom claims
	UserID string   `json:"https://api.financetracker.local/user_id,omitempty"`
	Roles  []string `json:"https://api.financetracker.local/roles,omitempty"`
}

func (c *Client) ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	jwks, err := c.GetJWKS(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %w", err)
	}

	token, err := jose.ParseSigned(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	var claims TokenClaims
	var found bool

	for _, key := range jwks.Keys {
		payload, err := token.Verify(key)
		if err == nil {
			if err := json.Unmarshal(payload, &claims); err != nil {
				continue
			}
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("no matching key found")
	}

	// Validate standard claims
	if err := c.validateClaims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func (c *Client) validateClaims(claims *TokenClaims) error {
	// Check issuer
	if claims.Issuer != c.GetIssuer() {
		return fmt.Errorf("invalid issuer: %s", claims.Issuer)
	}

	// Check audience
	audienceFound := false
	for _, aud := range claims.Audience {
		if aud == c.audience {
			audienceFound = true
			break
		}
	}
	if !audienceFound {
		return fmt.Errorf("invalid audience")
	}

	// Check expiration
	if time.Now().Unix() > claims.ExpiresAt {
		return fmt.Errorf("token expired")
	}

	return nil
}

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
	UpdatedAt     string `json:"updated_at"`
}

func (c *Client) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	userInfoURL := fmt.Sprintf("https://%s/userinfo", c.domain)
	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

func (c *Client) BuildLoginURL(state, redirectURI string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", "") // This should be injected from env
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", "openid profile email")
	params.Set("state", state)
	params.Set("audience", c.audience)

	return fmt.Sprintf("https://%s/authorize?%s", c.domain, params.Encode())
}

func (c *Client) BuildLogoutURL(returnTo string) string {
	params := url.Values{}
	params.Set("client_id", "") // This should be injected from env
	params.Set("returnTo", returnTo)

	return fmt.Sprintf("https://%s/v2/logout?%s", c.domain, params.Encode())
}
