package value

import (
	"regexp"
	"strings"
	
	"financetracker/internal/domain/common"
)

// Auth0ID Auth0ユーザーIDを表現する値オブジェクト
type Auth0ID struct {
	value string
}

// Auth0のユーザーID形式の正規表現
// 通常は "auth0|" または "google-oauth2|" などのプロバイダー接頭辞が付く
var auth0IDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+\|[a-zA-Z0-9_-]+$`)

// NewAuth0ID 新しいAuth0IDインスタンスを作成
func NewAuth0ID(value string) (*Auth0ID, error) {
	auth0ID := strings.TrimSpace(value)
	
	if err := validateAuth0ID(auth0ID); err != nil {
		return nil, err
	}

	return &Auth0ID{
		value: auth0ID,
	}, nil
}

// Value Auth0 ID値を取得
func (a Auth0ID) Value() string {
	return a.value
}

// GetProvider プロバイダー名を取得（例：auth0, google-oauth2）
func (a Auth0ID) GetProvider() string {
	parts := strings.Split(a.value, "|")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

// GetProviderUserID プロバイダー固有のユーザーIDを取得
func (a Auth0ID) GetProviderUserID() string {
	parts := strings.Split(a.value, "|")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

// IsAuth0Native Auth0ネイティブユーザーかどうかを判定
func (a Auth0ID) IsAuth0Native() bool {
	return strings.HasPrefix(a.value, "auth0|")
}

// IsGoogleOAuth Googleソーシャルログインユーザーかどうかを判定
func (a Auth0ID) IsGoogleOAuth() bool {
	return strings.HasPrefix(a.value, "google-oauth2|")
}

// IsSocialLogin ソーシャルログインユーザーかどうかを判定
func (a Auth0ID) IsSocialLogin() bool {
	return !a.IsAuth0Native()
}

// String 文字列表現
func (a Auth0ID) String() string {
	return a.value
}

// Equals Auth0IDの同一性判定
func (a Auth0ID) Equals(other Auth0ID) bool {
	return a.value == other.value
}

// validateAuth0ID Auth0 IDのバリデーション
func validateAuth0ID(auth0ID string) error {
	// 空文字チェック
	if auth0ID == "" {
		return common.NewValidationError("auth0_id", auth0ID, "Auth0 user ID is required")
	}

	// 長さチェック（Auth0の制限に基づく）
	if len(auth0ID) > 128 {
		return common.NewValidationError("auth0_id", auth0ID, "Auth0 user ID must be 128 characters or less")
	}

	// フォーマットチェック
	if !auth0IDRegex.MatchString(auth0ID) {
		return common.NewValidationError("auth0_id", auth0ID, 
			"Auth0 user ID must be in format 'provider|user_id' (e.g., auth0|123456, google-oauth2|abcdef)")
	}

	// パイプ文字の数をチェック
	pipeCount := strings.Count(auth0ID, "|")
	if pipeCount != 1 {
		return common.NewValidationError("auth0_id", auth0ID, "Auth0 user ID must contain exactly one pipe symbol")
	}

	// プロバイダー部分とユーザーID部分をチェック
	parts := strings.Split(auth0ID, "|")
	provider := parts[0]
	userID := parts[1]

	if provider == "" {
		return common.NewValidationError("auth0_id", auth0ID, "Auth0 provider cannot be empty")
	}

	if userID == "" {
		return common.NewValidationError("auth0_id", auth0ID, "Auth0 user ID part cannot be empty")
	}

	return nil
}