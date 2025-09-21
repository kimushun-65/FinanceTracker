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

	// FindNewUsers 新規ユーザー（指定日数以内に登録）のリストを取得
	FindNewUsers(ctx context.Context, withinDays int, pagination *repository.Pagination) (*repository.PagedResult[*entity.User], error)

}


// UserStats ユーザー統計情報
type UserStats struct {
	TotalUsers        int64 `json:"total_users"`         // 総ユーザー数
	NewUsersThisMonth int64 `json:"new_users_this_month"` // 今月の新規ユーザー数
	SocialLoginUsers  int64 `json:"social_login_users"`  // ソーシャルログインユーザー数
}

