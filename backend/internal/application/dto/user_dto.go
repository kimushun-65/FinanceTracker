package dto

import (
	"time"

	"github.com/google/uuid"

	userDomain "financetracker/internal/domain/user/entity"
)

// UserResponse ユーザー情報レスポンス
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Auth0ID   string    `json:"auth0_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserRequest ユーザー更新リクエスト
type UpdateUserRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=1,max=100"`
	Email *string `json:"email" binding:"omitempty,email"`
}

// CreateUserRequest ユーザー作成リクエスト（内部利用）
type CreateUserRequest struct {
	Auth0ID string `json:"auth0_id" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Name    string `json:"name" binding:"required,min=1,max=100"`
}

// UserFromDomain ドメインエンティティからDTOへの変換
func UserFromDomain(user *userDomain.User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:        user.ID,
		Auth0ID:   user.Auth0UserID().String(),
		Email:     user.Email().String(),
		Name:      user.Name(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
