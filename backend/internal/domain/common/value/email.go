package value

import (
	"regexp"
	"strings"

	"financetracker/internal/domain/common"
)

// Email メールアドレスを表現する値オブジェクト
type Email struct {
	value string
}

// RFC 5322に準拠した簡易的なメールアドレス正規表現
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail 新しいEmailインスタンスを作成
func NewEmail(value string) (*Email, error) {
	email := strings.TrimSpace(value)
	
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	return &Email{
		value: strings.ToLower(email), // 小文字で正規化
	}, nil
}

// Value メールアドレスの値を取得
func (e Email) Value() string {
	return e.value
}

// GetDomain ドメイン部分を取得
func (e Email) GetDomain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// GetLocalPart ローカル部分（@より前）を取得
func (e Email) GetLocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// Equals メールアドレスの同一性判定
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// String 文字列表現
func (e Email) String() string {
	return e.value
}

// validateEmail メールアドレスのバリデーション
func validateEmail(email string) error {
	// 空文字チェック
	if email == "" {
		return common.NewValidationError("email", email, "email address is required")
	}

	// 長さチェック（RFC 5321では320文字まで、実用的には255文字）
	if len(email) > 255 {
		return common.NewValidationError("email", email, "email address must be 255 characters or less")
	}

	// フォーマットチェック
	if !emailRegex.MatchString(email) {
		return common.NewValidationError("email", email, "invalid email address format")
	}

	// @の数をチェック
	atCount := strings.Count(email, "@")
	if atCount != 1 {
		return common.NewValidationError("email", email, "email address must contain exactly one @ symbol")
	}

	// ドメイン部分の追加チェック
	parts := strings.Split(email, "@")
	domain := parts[1]
	
	// ドメインが空でないかチェック
	if domain == "" {
		return common.NewValidationError("email", email, "email domain cannot be empty")
	}

	// ドメインの先頭・末尾がピリオドでないかチェック
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return common.NewValidationError("email", email, "email domain cannot start or end with a period")
	}

	// 連続するピリオドがないかチェック
	if strings.Contains(domain, "..") {
		return common.NewValidationError("email", email, "email domain cannot contain consecutive periods")
	}

	return nil
}