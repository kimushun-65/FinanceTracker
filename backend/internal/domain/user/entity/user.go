package entity

import (
	"strings"

	"financetracker/internal/domain/common"
	"financetracker/internal/domain/common/value"
	userValue "financetracker/internal/domain/user/value"
)

// User ユーザーエンティティ
type User struct {
	common.BaseEntity
	auth0UserID userValue.Auth0ID
	email       value.Email
	name        string
}

// NewUser 新しいUserエンティティを作成
func NewUser(auth0UserID userValue.Auth0ID, email value.Email, name string) (*User, error) {
	if err := validateUserName(name); err != nil {
		return nil, err
	}

	return &User{
		BaseEntity:  common.NewBaseEntity(),
		auth0UserID: auth0UserID,
		email:       email,
		name:        strings.TrimSpace(name),
	}, nil
}

// ReconstructUser 既存のデータからUserエンティティを再構築（リポジトリから取得時に使用）
func ReconstructUser(
	baseEntity common.BaseEntity,
	auth0UserID userValue.Auth0ID,
	email value.Email,
	name string,
) *User {
	return &User{
		BaseEntity:  baseEntity,
		auth0UserID: auth0UserID,
		email:       email,
		name:        name,
	}
}

// Auth0UserID Auth0ユーザーIDを取得
func (u User) Auth0UserID() userValue.Auth0ID {
	return u.auth0UserID
}

// Email メールアドレスを取得
func (u User) Email() value.Email {
	return u.email
}

// Name ユーザー名を取得
func (u User) Name() string {
	return u.name
}

// UpdateProfile プロファイル情報を更新
func (u *User) UpdateProfile(name string, email value.Email) error {
	if err := validateUserName(name); err != nil {
		return err
	}

	u.name = strings.TrimSpace(name)
	u.email = email
	u.UpdateTimestamp()

	return nil
}

// IsSocialLoginUser ソーシャルログインユーザーかどうかを判定
func (u User) IsSocialLoginUser() bool {
	return u.auth0UserID.IsSocialLogin()
}

// GetDisplayName 表示用の名前を取得
func (u User) GetDisplayName() string {
	if u.name != "" {
		return u.name
	}
	// 名前が設定されていない場合はメールアドレスのローカル部分を使用
	return u.email.GetLocalPart()
}


// validateUserName ユーザー名のバリデーション
func validateUserName(name string) error {
	name = strings.TrimSpace(name)

	// 空文字チェック
	if name == "" {
		return common.NewValidationError("name", name, "user name is required")
	}

	// 長さチェック
	if len(name) > 100 {
		return common.NewValidationError("name", name, "user name must be 100 characters or less")
	}

	// 禁止文字チェック（基本的な制御文字のみ）
	for _, char := range name {
		if char < 32 || char == 127 { // ASCII制御文字
			return common.NewValidationError("name", name, "user name contains invalid characters")
		}
	}

	return nil
}