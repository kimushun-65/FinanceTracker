package dto

// UserInfo 認証プロバイダーから取得したユーザー情報
type UserInfo struct {
	Sub           string
	Name          string
	Email         string
	EmailVerified bool
	Picture       string
}

// TokenClaims トークンのクレーム情報
type TokenClaims struct {
	Subject string
	Roles   []string
}
