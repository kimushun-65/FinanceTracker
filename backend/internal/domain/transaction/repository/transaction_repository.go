// Package repository トランザクション関連のリポジトリインターフェースを定義
package repository

import (
	"context"
	"time"

	"financetracker/internal/domain/transaction/entity"

	"github.com/google/uuid"
)

// TransactionRepository トランザクションリポジトリインターフェース
type TransactionRepository interface {
	// Save トランザクションを保存
	Save(ctx context.Context, transaction entity.Transaction) error
	
	// FindByID IDでトランザクションを取得
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	
	// FindByUserID ユーザーIDでトランザクションを取得
	FindByUserID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]entity.Transaction, error)
	
	// FindByAccountID アカウントIDでトランザクションを取得
	FindByAccountID(ctx context.Context, accountID uuid.UUID, limit int, offset int) ([]entity.Transaction, error)
	
	// FindByUserIDAndDateRange ユーザーIDと日付範囲でトランザクションを取得
	FindByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]entity.Transaction, error)
	
	// Update トランザクションを更新
	Update(ctx context.Context, transaction entity.Transaction) error
	
	// Delete トランザクションを削除
	Delete(ctx context.Context, id uuid.UUID) error
	
	// CountByUserID ユーザーのトランザクション数を取得
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}