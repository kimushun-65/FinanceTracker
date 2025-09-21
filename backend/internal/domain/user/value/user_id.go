package value

import (
	"github.com/google/uuid"
	"financetracker/internal/domain/common"
)

// UserID ユーザーIDを表現する値オブジェクト
type UserID struct {
	value uuid.UUID
}

// NewUserID 新しいUserIDインスタンスを作成
func NewUserID(id uuid.UUID) *UserID {
	return &UserID{
		value: id,
	}
}

// NewUserIDFromString 文字列からUserIDインスタンスを作成
func NewUserIDFromString(idStr string) (*UserID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, common.NewValidationError("user_id", idStr, "invalid user ID format")
	}
	
	return &UserID{
		value: id,
	}, nil
}

// GenerateUserID 新しいUUIDを生成してUserIDインスタンスを作成
func GenerateUserID() *UserID {
	return &UserID{
		value: uuid.New(),
	}
}

// Value UUID値を取得
func (u UserID) Value() uuid.UUID {
	return u.value
}

// String 文字列表現
func (u UserID) String() string {
	return u.value.String()
}

// Equals UserIDの同一性判定
func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}

// IsZero ゼロ値かどうかを判定
func (u UserID) IsZero() bool {
	return u.value == uuid.Nil
}