package repository

import (
	"context"

	"financetracker/internal/domain/common/repository"
	"financetracker/internal/domain/user/entity"
	userValue "financetracker/internal/domain/user/value"
)

// UserRepository ユーザーリポジトリのインターフェース
type UserRepository interface {
	repository.BaseRepository[*entity.User]

	// FindByAuth0UserID Auth0ユーザーIDでユーザーを検索
	FindByAuth0UserID(ctx context.Context, auth0UserID userValue.Auth0ID) (*entity.User, error)

	// ExistsByAuth0UserID Auth0ユーザーIDでユーザーが存在するかチェック
	ExistsByAuth0UserID(ctx context.Context, auth0UserID userValue.Auth0ID) (bool, error)

	// FindAll 全ユーザーを取得
	FindAll(ctx context.Context) ([]*entity.User, error)
}
