package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"financetracker/internal/domain/common"
	commonValue "financetracker/internal/domain/common/value"
	transactionDomain "financetracker/internal/domain/transaction/entity"
	transactionRepo "financetracker/internal/domain/transaction/repository"
	transactionValue "financetracker/internal/domain/transaction/value"
	"financetracker/internal/infrastructure/gorm/model"
)

// TransactionRepository 取引リポジトリの実装
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository 新しいTransactionRepositoryを作成
func NewTransactionRepository(db *gorm.DB) transactionRepo.TransactionRepository {
	return &TransactionRepository{db: db}
}

// Save 取引を保存（新規作成または更新）
func (r *TransactionRepository) Save(ctx context.Context, transaction transactionDomain.Transaction) error {
	transactionModel := r.toModel(&transaction)

	// IDが存在する場合は更新、存在しない場合は作成
	result := r.db.WithContext(ctx).Save(transactionModel)
	if result.Error != nil {
		return fmt.Errorf("取引の保存に失敗しました: %w", result.Error)
	}

	return nil
}

// FindByID IDで取引を取得
func (r *TransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*transactionDomain.Transaction, error) {
	var transactionModel model.Transaction
	result := r.db.WithContext(ctx).First(&transactionModel, "id = ?", id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("取引の取得に失敗しました: %w", result.Error)
	}

	return r.toDomain(&transactionModel)
}

// FindByUserID ユーザーIDで取引を取得
func (r *TransactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]transactionDomain.Transaction, error) {
	var transactionModels []model.Transaction
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("date DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&transactionModels)
	if result.Error != nil {
		return nil, fmt.Errorf("取引一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	transactions := make([]transactionDomain.Transaction, len(transactionModels))
	for i, transactionModel := range transactionModels {
		transaction, err := r.toDomain(&transactionModel)
		if err != nil {
			return nil, err
		}
		transactions[i] = *transaction
	}

	return transactions, nil
}

// FindByAccountID アカウントIDで取引を取得
func (r *TransactionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID, limit int, offset int) ([]transactionDomain.Transaction, error) {
	var transactionModels []model.Transaction
	query := r.db.WithContext(ctx).Where("account_id = ?", accountID).Order("date DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&transactionModels)
	if result.Error != nil {
		return nil, fmt.Errorf("取引一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	transactions := make([]transactionDomain.Transaction, len(transactionModels))
	for i, transactionModel := range transactionModels {
		transaction, err := r.toDomain(&transactionModel)
		if err != nil {
			return nil, err
		}
		transactions[i] = *transaction
	}

	return transactions, nil
}

// FindByUserIDAndDateRange ユーザーIDと日付範囲で取引を取得
func (r *TransactionRepository) FindByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]transactionDomain.Transaction, error) {
	var transactionModels []model.Transaction
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Order("date DESC, created_at DESC").
		Find(&transactionModels)

	if result.Error != nil {
		return nil, fmt.Errorf("取引一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	transactions := make([]transactionDomain.Transaction, len(transactionModels))
	for i, transactionModel := range transactionModels {
		transaction, err := r.toDomain(&transactionModel)
		if err != nil {
			return nil, err
		}
		transactions[i] = *transaction
	}

	return transactions, nil
}

// Update 取引を更新
func (r *TransactionRepository) Update(ctx context.Context, transaction transactionDomain.Transaction) error {
	transactionModel := r.toModel(&transaction)

	result := r.db.WithContext(ctx).Model(&model.Transaction{}).
		Where("id = ?", transaction.ID).
		Updates(transactionModel)

	if result.Error != nil {
		return fmt.Errorf("取引の更新に失敗しました: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("更新対象の取引が見つかりません: %s", transaction.ID)
	}

	return nil
}

// Delete 取引を削除
func (r *TransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Transaction{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("取引の削除に失敗しました: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("削除対象の取引が見つかりません: %s", id)
	}

	return nil
}

// CountByUserID ユーザーの取引数を取得
func (r *TransactionRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Transaction{}).Where("user_id = ?", userID).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("取引数の取得に失敗しました: %w", result.Error)
	}

	return count, nil
}

// toModel ドメインモデルからGORMモデルへの変換
func (r *TransactionRepository) toModel(transaction *transactionDomain.Transaction) *model.Transaction {
	amount := transaction.Amount()
	description := transaction.Description().String()
	var descriptionPtr *string
	if description != "" {
		descriptionPtr = &description
	}

	return &model.Transaction{
		Base: model.Base{
			ID:        transaction.ID,
			CreatedAt: transaction.CreatedAt,
			UpdatedAt: transaction.UpdatedAt,
		},
		UserID:     transaction.UserID(),
		AccountID:  transaction.AccountID(),
		CategoryID: transaction.CategoryID(),
		Amount:     decimal.NewFromInt(amount.Amount()),
		Type: func() model.TransactionType {
			// ドメイン（小文字）からDB（大文字）への変換
			switch transaction.Type().String() {
			case "income":
				return model.TransactionTypeIncome
			case "expense":
				return model.TransactionTypeExpense
			default:
				return model.TransactionType(transaction.Type().String())
			}
		}(),
		Date:        transaction.Date(),
		Description: descriptionPtr,
	}
}

// toDomain GORMモデルからドメインモデルへの変換
func (r *TransactionRepository) toDomain(transactionModel *model.Transaction) (*transactionDomain.Transaction, error) {
	// 金額
	money, err := commonValue.NewMoney(transactionModel.Amount.IntPart(), "JPY")
	if err != nil {
		return nil, fmt.Errorf("金額の作成に失敗しました: %w", err)
	}

	// 取引タイプ (DBは大文字、ドメインは小文字を期待)
	typeStr := string(transactionModel.Type)
	switch typeStr {
	case "INCOME":
		typeStr = "income"
	case "EXPENSE":
		typeStr = "expense"
	}
	txType, err := transactionValue.NewTransactionType(typeStr)
	if err != nil {
		return nil, fmt.Errorf("取引タイプの作成に失敗しました: %w", err)
	}

	// 説明
	description := ""
	if transactionModel.Description != nil {
		description = *transactionModel.Description
	}
	desc, err := transactionValue.NewDescription(description)
	if err != nil {
		return nil, fmt.Errorf("説明の作成に失敗しました: %w", err)
	}

	// ドメインエンティティを作成
	transaction, err := transactionDomain.NewTransaction(
		transactionModel.UserID,
		transactionModel.AccountID,
		transactionModel.CategoryID,
		*money,
		txType,
		transactionModel.Date,
		desc,
	)
	if err != nil {
		return nil, fmt.Errorf("取引エンティティの作成に失敗しました: %w", err)
	}

	// BaseEntityを設定
	transaction.BaseEntity = common.BaseEntity{
		ID:        transactionModel.ID,
		CreatedAt: transactionModel.CreatedAt,
		UpdatedAt: transactionModel.UpdatedAt,
	}

	return &transaction, nil
}
