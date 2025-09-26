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

// Client Auth0クライアントを表現する構造体
type Client struct {
	domain     string
	clientID   string
	audience   string
	jwks       *jose.JSONWebKeySet
	jwksMutex  sync.RWMutex
	jwksExpiry time.Time
}

// NewClient 新しいAuth0クライアントインスタンスを作成
func NewClient(domain, clientID, audience string) *Client {
	return &Client{
		domain:   domain,
		clientID: clientID,
		audience: audience,
	}
}

// GetJWKS Auth0からJSON Web Key Setを取得しキャッシュする
func (c *Client) GetJWKS(ctx context.Context) (*jose.JSONWebKeySet, error) {
	c.jwksMutex.RLock()
	if c.jwks != nil && time.Now().Before(c.jwksExpiry) {
		defer c.jwksMutex.RUnlock()
		return c.jwks, nil
	}
	c.jwksMutex.RUnlock()

	c.jwksMutex.Lock()
	defer c.jwksMutex.Unlock()

	// 書き込みロック取得後に再度チェック
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

// GetDomain Auth0ドメインを取得
func (c *Client) GetDomain() string {
	return c.domain
}

// GetAudience Auth0 Audienceを取得
func (c *Client) GetAudience() string {
	return c.audience
}

// GetIssuer Auth0発行者URLを取得
func (c *Client) GetIssuer() string {
	return fmt.Sprintf("https://%s/", c.domain)
}

// TokenClaims JWTトークンのクレーム情報
type TokenClaims struct {
	Issuer    string   `json:"iss"`
	Subject   string   `json:"sub"`
	Audience  []string `json:"aud"`
	ExpiresAt int64    `json:"exp"`
	IssuedAt  int64    `json:"iat"`
	Scope     string   `json:"scope,omitempty"`
	Azp       string   `json:"azp,omitempty"`

	// カスタムクレーム
	UserID string   `json:"https://api.financetracker.local/user_id,omitempty"`
	Roles  []string `json:"https://api.financetracker.local/roles,omitempty"`
}

// ValidateToken JWTトークンを検証しクレームを取得
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

	// 標準クレームを検証
	if err := c.validateClaims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

// validateClaims クレームの妥当性を検証
func (c *Client) validateClaims(claims *TokenClaims) error {
	// 発行者をチェック
	if claims.Issuer != c.GetIssuer() {
		return fmt.Errorf("invalid issuer: %s", claims.Issuer)
	}

	// Audienceをチェック
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

	// 有効期限をチェック
	if time.Now().Unix() > claims.ExpiresAt {
		return fmt.Errorf("token expired")
	}

	return nil
}

// UserInfo Auth0ユーザー情報
type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	UpdatedAt     string `json:"updated_at"`
}

// GetUserInfo アクセストークンを使用してユーザー情報を取得
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

// BuildLoginURL Auth0ログインURLを構築
func (c *Client) BuildLoginURL(state, redirectURI string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", c.clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", "openid profile email")
	params.Set("state", state)
	params.Set("audience", c.audience)

	return fmt.Sprintf("https://%s/authorize?%s", c.domain, params.Encode())
}

// BuildLogoutURL Auth0ログアウトURLを構築
func (c *Client) BuildLogoutURL(returnTo string) string {
	params := url.Values{}
	params.Set("client_id", c.clientID)
	params.Set("returnTo", returnTo)

	return fmt.Sprintf("https://%s/v2/logout?%s", c.domain, params.Encode())
}
