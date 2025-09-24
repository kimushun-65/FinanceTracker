package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	budgetDomain "financetracker/internal/domain/budget/entity"
	budgetRepo "financetracker/internal/domain/budget/repository"
	budgetValue "financetracker/internal/domain/budget/value"
	"financetracker/internal/domain/common"
	commonValue "financetracker/internal/domain/common/value"
	"financetracker/internal/infrastructure/gorm/model"
)

// BudgetRepository 予算リポジトリの実装
type BudgetRepository struct {
	db *gorm.DB
}

// NewBudgetRepository 新しいBudgetRepositoryを作成
func NewBudgetRepository(db *gorm.DB) budgetRepo.BudgetRepository {
	return &BudgetRepository{db: db}
}

// Save 予算を保存（新規作成または更新）
func (r *BudgetRepository) Save(ctx context.Context, budget budgetDomain.Budget) error {
	budgetModel := r.toModel(&budget)

	// IDが存在する場合は更新、存在しない場合は作成
	result := r.db.WithContext(ctx).Save(budgetModel)
	if result.Error != nil {
		return fmt.Errorf("予算の保存に失敗しました: %w", result.Error)
	}

	return nil
}

// FindByID IDで予算を取得
func (r *BudgetRepository) FindByID(ctx context.Context, id uuid.UUID) (*budgetDomain.Budget, error) {
	var budgetModel model.Budget
	result := r.db.WithContext(ctx).First(&budgetModel, "id = ?", id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("予算の取得に失敗しました: %w", result.Error)
	}

	return r.toDomain(&budgetModel)
}

// FindByUserID ユーザーIDで予算を取得
func (r *BudgetRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]budgetDomain.Budget, error) {
	var budgetModels []model.Budget
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("start_date DESC").Find(&budgetModels)

	if result.Error != nil {
		return nil, fmt.Errorf("予算一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	budgets := make([]budgetDomain.Budget, len(budgetModels))
	for i, budgetModel := range budgetModels {
		budget, err := r.toDomain(&budgetModel)
		if err != nil {
			return nil, err
		}
		budgets[i] = *budget
	}

	return budgets, nil
}

// FindByUserIDAndCategoryID ユーザーIDとカテゴリIDで予算を取得
func (r *BudgetRepository) FindByUserIDAndCategoryID(ctx context.Context, userID, categoryID uuid.UUID) ([]budgetDomain.Budget, error) {
	var budgetModels []model.Budget
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND category_id = ?", userID, categoryID).
		Order("start_date DESC").
		Find(&budgetModels)

	if result.Error != nil {
		return nil, fmt.Errorf("予算一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	budgets := make([]budgetDomain.Budget, len(budgetModels))
	for i, budgetModel := range budgetModels {
		budget, err := r.toDomain(&budgetModel)
		if err != nil {
			return nil, err
		}
		budgets[i] = *budget
	}

	return budgets, nil
}

// FindActiveByUserID ユーザーのアクティブな予算を取得
func (r *BudgetRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]budgetDomain.Budget, error) {
	var budgetModels []model.Budget
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("start_date DESC").
		Find(&budgetModels)

	if result.Error != nil {
		return nil, fmt.Errorf("アクティブな予算一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	budgets := make([]budgetDomain.Budget, len(budgetModels))
	for i, budgetModel := range budgetModels {
		budget, err := r.toDomain(&budgetModel)
		if err != nil {
			return nil, err
		}
		budgets[i] = *budget
	}

	return budgets, nil
}

// FindByUserIDAndDate ユーザーIDと日付で有効な予算を取得
func (r *BudgetRepository) FindByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]budgetDomain.Budget, error) {
	var budgetModels []model.Budget
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND start_date <= ?", userID, date)

	// 終了日の条件は、NULLの場合も考慮
	query = query.Where("end_date IS NULL OR end_date >= ?", date)

	result := query.Order("start_date DESC").Find(&budgetModels)
	if result.Error != nil {
		return nil, fmt.Errorf("予算一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	budgets := make([]budgetDomain.Budget, len(budgetModels))
	for i, budgetModel := range budgetModels {
		budget, err := r.toDomain(&budgetModel)
		if err != nil {
			return nil, err
		}
		budgets[i] = *budget
	}

	return budgets, nil
}

// Update 予算を更新
func (r *BudgetRepository) Update(ctx context.Context, budget budgetDomain.Budget) error {
	budgetModel := r.toModel(&budget)

	result := r.db.WithContext(ctx).Model(&model.Budget{}).
		Where("id = ?", budget.ID).
		Updates(budgetModel)

	if result.Error != nil {
		return fmt.Errorf("予算の更新に失敗しました: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("更新対象の予算が見つかりません: %s", budget.ID)
	}

	return nil
}

// Delete 予算を削除
func (r *BudgetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Budget{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("予算の削除に失敗しました: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("削除対象の予算が見つかりません: %s", id)
	}

	return nil
}

// toModel ドメインモデルからGORMモデルへの変換
func (r *BudgetRepository) toModel(budget *budgetDomain.Budget) *model.Budget {
	amount := budget.Amount()

	return &model.Budget{
		Base: model.Base{
			ID:        budget.ID,
			CreatedAt: budget.CreatedAt,
			UpdatedAt: budget.UpdatedAt,
		},
		UserID:     budget.UserID(),
		CategoryID: budget.CategoryID(),
		Amount:     decimal.NewFromInt(amount.Amount()),
		PeriodType: model.PeriodType(budget.PeriodType().String()),
		StartDate:  budget.StartDate(),
		EndDate:    budget.EndDate(),
		IsActive:   budget.IsActive(),
	}
}

// toDomain GORMモデルからドメインモデルへの変換
func (r *BudgetRepository) toDomain(budgetModel *model.Budget) (*budgetDomain.Budget, error) {
	// 金額
	money, err := commonValue.NewMoney(budgetModel.Amount.IntPart(), "JPY")
	if err != nil {
		return nil, fmt.Errorf("金額の作成に失敗しました: %w", err)
	}

	// 期間タイプ
	periodType, err := budgetValue.NewPeriodType(string(budgetModel.PeriodType))
	if err != nil {
		return nil, fmt.Errorf("期間タイプの作成に失敗しました: %w", err)
	}

	// ドメインエンティティを作成
	budget, err := budgetDomain.NewBudget(
		budgetModel.UserID,
		budgetModel.CategoryID,
		*money,
		periodType,
		budgetModel.StartDate,
		budgetModel.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("予算エンティティの作成に失敗しました: %w", err)
	}

	// BaseEntityを設定
	budget.BaseEntity = common.BaseEntity{
		ID:        budgetModel.ID,
		CreatedAt: budgetModel.CreatedAt,
		UpdatedAt: budgetModel.UpdatedAt,
	}

	// アクティブ状態を設定
	if !budgetModel.IsActive {
		budget.Deactivate()
	}

	return &budget, nil
}
