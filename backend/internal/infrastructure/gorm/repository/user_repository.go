package repository

import (
	"context"
	"errors"
	"time"

	"financetracker/internal/domain/common"
	"financetracker/internal/domain/common/value"
	"financetracker/internal/domain/user/entity"
	userValue "financetracker/internal/domain/user/value"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserModel ユーザーのデータベースモデル
type UserModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Auth0UserID   string    `gorm:"column:auth0_id;type:varchar(255);unique;not null"`
	Email         string    `gorm:"type:varchar(255);unique;not null"`
	Name          string    `gorm:"type:varchar(100);not null"`
	EmailVerified bool      `gorm:"type:boolean;default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName Userモデルのテーブル名を指定
func (UserModel) TableName() string {
	return "users"
}

// UserRepository repository.UserRepositoryインターフェースの実装
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 新しいUserRepositoryインスタンスを作成
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Save ユーザーエンティティをデータベースに保存
func (r *UserRepository) Save(ctx context.Context, user *entity.User) error {
	model := r.toModel(user)

	result := r.db.WithContext(ctx).Save(&model)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindByID IDでユーザーを検索
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var model UserModel
	result := r.db.WithContext(ctx).First(&model, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.NewNotFoundError("user", id)
		}
		return nil, result.Error
	}

	return r.toDomain(&model)
}

// FindByAuth0UserID Auth0ユーザーIDでユーザーを検索
func (r *UserRepository) FindByAuth0UserID(ctx context.Context, auth0UserID userValue.Auth0ID) (*entity.User, error) {
	var model UserModel
	result := r.db.WithContext(ctx).Where("auth0_id = ?", auth0UserID.Value()).First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.NewNotFoundErrorWithMessage("user", uuid.Nil, "user not found with Auth0 ID")
		}
		return nil, result.Error
	}

	return r.toDomain(&model)
}

// ExistsByAuth0UserID 指定されたAuth0 IDのユーザーが存在するか確認
func (r *UserRepository) ExistsByAuth0UserID(ctx context.Context, auth0UserID userValue.Auth0ID) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&UserModel{}).Where("auth0_id = ?", auth0UserID.Value()).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

// Delete IDでユーザーを削除
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&UserModel{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.NewNotFoundError("user", id)
	}

	return nil
}

// Exists ユーザーが存在するか確認
func (r *UserRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", id).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

// toModel ドメインエンティティをデータベースモデルに変換
func (r *UserRepository) toModel(user *entity.User) *UserModel {
	return &UserModel{
		ID:            user.GetID(),
		Auth0UserID:   user.Auth0UserID().Value(),
		Email:         user.Email().Value(),
		Name:          user.Name(),
		EmailVerified: user.IsEmailVerified(),
		CreatedAt:     user.GetCreatedAt(),
		UpdatedAt:     user.GetUpdatedAt(),
	}
}

// toDomain データベースモデルをドメインエンティティに変換
func (r *UserRepository) toDomain(model *UserModel) (*entity.User, error) {
	auth0ID, err := userValue.NewAuth0ID(model.Auth0UserID)
	if err != nil {
		return nil, err
	}

	email, err := value.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	baseEntity := common.ReconstructBaseEntity(
		model.ID,
		model.CreatedAt,
		model.UpdatedAt,
	)

	return entity.ReconstructUser(
		baseEntity,
		*auth0ID,
		*email,
		model.Name,
		model.EmailVerified,
	), nil
}

// FindAll 全ユーザーを取得
func (r *UserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	var models []UserModel
	result := r.db.WithContext(ctx).Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	users := make([]*entity.User, 0, len(models))
	for _, model := range models {
		user, err := r.toDomain(&model)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
