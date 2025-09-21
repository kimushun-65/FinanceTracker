package entity

import (
	"strings"
	"time"

	"financetracker/internal/domain/common"
	"financetracker/internal/domain/common/value"
	userValue "financetracker/internal/domain/user/value"
)

// AccountMovement 口座の入出金履歴エンティティ
type AccountMovement struct {
	common.BaseEntity
	userID     userValue.UserID
	accountID  common.BaseEntity // Accountエンティティへの参照用
	amount     value.Money
	occurredAt time.Time
	note       string
}

// NewAccountMovement 新しいAccountMovementエンティティを作成
func NewAccountMovement(
	userID userValue.UserID,
	accountID common.BaseEntity,
	amount value.Money,
	occurredAt time.Time,
	note string,
) (*AccountMovement, error) {
	if err := validateMovementNote(note); err != nil {
		return nil, err
	}

	return &AccountMovement{
		BaseEntity: common.NewBaseEntity(),
		userID:     userID,
		accountID:  accountID,
		amount:     amount,
		occurredAt: occurredAt,
		note:       strings.TrimSpace(note),
	}, nil
}

// ReconstructAccountMovement 既存のデータからAccountMovementエンティティを再構築
func ReconstructAccountMovement(
	baseEntity common.BaseEntity,
	userID userValue.UserID,
	accountID common.BaseEntity,
	amount value.Money,
	occurredAt time.Time,
	note string,
) *AccountMovement {
	return &AccountMovement{
		BaseEntity: baseEntity,
		userID:     userID,
		accountID:  accountID,
		amount:     amount,
		occurredAt: occurredAt,
		note:       note,
	}
}

// UserID ユーザーIDを取得
func (am AccountMovement) UserID() userValue.UserID {
	return am.userID
}

// AccountID 口座IDを取得
func (am AccountMovement) AccountID() common.BaseEntity {
	return am.accountID
}

// Amount 金額を取得
func (am AccountMovement) Amount() value.Money {
	return am.amount
}

// OccurredAt 発生日時を取得
func (am AccountMovement) OccurredAt() time.Time {
	return am.occurredAt
}

// Note メモを取得
func (am AccountMovement) Note() string {
	return am.note
}

// IsDeposit 入金かどうかを判定
func (am AccountMovement) IsDeposit() bool {
	return am.amount.IsPositive()
}

// IsWithdrawal 出金かどうかを判定
func (am AccountMovement) IsWithdrawal() bool {
	return am.amount.IsNegative()
}

// GetAbsoluteAmount 絶対値の金額を取得
func (am AccountMovement) GetAbsoluteAmount() value.Money {
	return *am.amount.Abs()
}

// GetMovementType 移動タイプを取得
func (am AccountMovement) GetMovementType() MovementType {
	if am.IsDeposit() {
		return MovementTypeDeposit
	}
	return MovementTypeWithdrawal
}

// UpdateNote メモを更新
func (am *AccountMovement) UpdateNote(newNote string) error {
	if err := validateMovementNote(newNote); err != nil {
		return err
	}

	am.note = strings.TrimSpace(newNote)
	am.UpdateTimestamp()
	return nil
}

// CanModify 変更可能かどうかを判定（一度記録された履歴は基本的に変更不可）
func (am AccountMovement) CanModify() bool {
	// 作成から24時間以内のみ変更可能
	return time.Since(am.GetCreatedAt()) <= 24*time.Hour
}

// GetDisplayInfo 表示用の移動情報を取得
func (am AccountMovement) GetDisplayInfo() MovementDisplayInfo {
	return MovementDisplayInfo{
		Amount:        am.GetAbsoluteAmount().Format(),
		Type:          am.GetMovementType(),
		OccurredAt:    am.occurredAt.Format("2006-01-02 15:04"),
		Note:          am.note,
		CanModify:     am.CanModify(),
	}
}

// MovementType 移動タイプ
type MovementType string

const (
	MovementTypeDeposit    MovementType = "deposit"    // 入金
	MovementTypeWithdrawal MovementType = "withdrawal" // 出金
)

// MovementDisplayInfo 表示用の移動情報
type MovementDisplayInfo struct {
	Amount     string       `json:"amount"`
	Type       MovementType `json:"type"`
	OccurredAt string       `json:"occurred_at"`
	Note       string       `json:"note"`
	CanModify  bool         `json:"can_modify"`
}

// validateMovementNote 移動メモのバリデーション
func validateMovementNote(note string) error {
	note = strings.TrimSpace(note)

	// 長さチェック（255文字以内）
	if len([]rune(note)) > 255 {
		return common.NewValidationError("note", note, "movement note must be 255 characters or less")
	}

	// 制御文字チェック
	for _, char := range note {
		if char < 32 || char == 127 { // ASCII制御文字
			return common.NewValidationError("note", note, "movement note contains invalid characters")
		}
	}

	return nil
}