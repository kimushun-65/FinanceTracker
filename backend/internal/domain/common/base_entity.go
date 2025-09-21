package common

import (
	"time"

	"github.com/google/uuid"
)

// BaseEntity 全てのドメインエンティティの共通フィールドを表現
type BaseEntity struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewBaseEntity UUIDを生成し、現在時刻を設定した新しいBaseEntityを作成
func NewBaseEntity() BaseEntity {
	now := time.Now()
	return BaseEntity{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateTimestamp UpdatedAtフィールドを現在時刻に更新
func (e *BaseEntity) UpdateTimestamp() {
	e.UpdatedAt = time.Now()
}

// GetID エンティティのIDを取得
func (e BaseEntity) GetID() uuid.UUID {
	return e.ID
}

// GetCreatedAt エンティティの作成時刻を取得
func (e BaseEntity) GetCreatedAt() time.Time {
	return e.CreatedAt
}

// GetUpdatedAt エンティティの最終更新時刻を取得
func (e BaseEntity) GetUpdatedAt() time.Time {
	return e.UpdatedAt
}

// Equals 二つのエンティティが同じIDを持つかチェック
func (e BaseEntity) Equals(other BaseEntity) bool {
	return e.ID == other.ID
}