package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"financetracker/internal/application/dto"
	budgetDomain "financetracker/internal/domain/budget/entity"
	budgetRepo "financetracker/internal/domain/budget/repository"
	budgetValue "financetracker/internal/domain/budget/value"
	"financetracker/internal/domain/common/value"
	userValue "financetracker/internal/domain/user/value"
	"financetracker/pkg/errors"
	"financetracker/pkg/logger"
)

// BudgetService 予算管理サービス
type BudgetService struct {
	budgetRepo budgetRepo.BudgetRepository
	logger     *logger.Logger
}

// NewBudgetService 新しいBudgetServiceを作成
func NewBudgetService(
	budgetRepo budgetRepo.BudgetRepository,
	logger *logger.Logger,
) *BudgetService {
	return &BudgetService{
		budgetRepo: budgetRepo,
		logger:     logger,
	}
}

// CreateBudget 予算を作成
func (s *BudgetService) CreateBudget(ctx context.Context, userID uuid.UUID, req *dto.CreateBudgetRequest) (*dto.BudgetResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// 予算名は不要（エンティティに名前フィールドがないため削除）

	// 金額を作成（デフォルトでJPY）
	money, err := value.NewMoney(req.Amount.IntPart(), "JPY")
	if err != nil {
		s.logger.Error("金額作成エラー",
			zap.Error(err),
			zap.String("amount", req.Amount.String()))
		return nil, errors.NewValidationError("無効な金額です")
	}

	// 予算期間タイプを作成
	periodType, err := budgetValue.NewPeriodType(req.Period)
	if err != nil {
		s.logger.Error("予算期間タイプ作成エラー",
			zap.Error(err),
			zap.String("period", req.Period))
		return nil, errors.NewValidationError("無効な予算期間タイプです")
	}

	// 日付の妥当性チェック
	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		return nil, errors.NewValidationError("終了日は開始日より後である必要があります")
	}

	// 予算を作成
	budget, err := budgetDomain.NewBudget(
		domainUserID.Value(),
		req.CategoryID,
		*money,
		periodType,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		s.logger.Error("予算作成エラー", zap.Error(err))
		return nil, errors.NewValidationError("予算の作成に失敗しました")
	}

	// リポジトリに保存
	if err := s.budgetRepo.Save(ctx, budget); err != nil {
		s.logger.Error("予算保存エラー", zap.Error(err))
		return nil, errors.NewInternalError("予算の保存に失敗しました", err)
	}

	// DTOに変換して返却
	return dto.BudgetFromDomain(&budget), nil
}

// GetBudget 予算情報を取得
func (s *BudgetService) GetBudget(ctx context.Context, userID, budgetID uuid.UUID) (*dto.BudgetResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから予算を取得
	budget, err := s.budgetRepo.FindByID(ctx, budgetID)
	if err != nil {
		s.logger.Error("予算取得エラー",
			zap.Error(err),
			zap.String("budgetID", budgetID.String()))
		return nil, errors.NewInternalError("予算情報の取得に失敗しました", err)
	}

	if budget == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("予算が見つかりません: %s", budgetID))
	}

	// ユーザーの予算であることを確認
	if budget.UserID() != domainUserID.Value() {
		return nil, errors.NewForbiddenError("この予算へのアクセス権限がありません")
	}

	// DTOに変換
	return dto.BudgetFromDomain(budget), nil
}

// GetBudgetsByUser ユーザーの予算一覧を取得
func (s *BudgetService) GetBudgetsByUser(ctx context.Context, userID uuid.UUID, params *dto.BudgetSearchParams) (*dto.BudgetListResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから予算一覧を取得
	budgets, err := s.budgetRepo.FindByUserID(ctx, domainUserID.Value())
	if err != nil {
		s.logger.Error("予算一覧取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()))
		return nil, errors.NewInternalError("予算一覧の取得に失敗しました", err)
	}

	// フィルタリング
	var filteredBudgets []*budgetDomain.Budget
	for i := range budgets {
		budget := &budgets[i]
		// アクティブフィルター
		if params.IsActive != nil && budget.IsActive() != *params.IsActive {
			continue
		}
		// カテゴリーIDフィルター
		if params.CategoryID != nil && budget.CategoryID() != *params.CategoryID {
			continue
		}
		// 期間フィルター
		if params.Period != nil && budget.PeriodType().String() != *params.Period {
			continue
		}
		filteredBudgets = append(filteredBudgets, budget)
	}

	// DTOに変換
	return &dto.BudgetListResponse{
		Budgets:    dto.BudgetsFromDomain(filteredBudgets),
		TotalCount: int64(len(filteredBudgets)),
	}, nil
}

// UpdateBudget 予算情報を更新
func (s *BudgetService) UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, req *dto.UpdateBudgetRequest) (*dto.BudgetResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから予算を取得
	budget, err := s.budgetRepo.FindByID(ctx, budgetID)
	if err != nil {
		s.logger.Error("予算取得エラー",
			zap.Error(err),
			zap.String("budgetID", budgetID.String()))
		return nil, errors.NewInternalError("予算情報の取得に失敗しました", err)
	}

	if budget == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("予算が見つかりません: %s", budgetID))
	}

	// ユーザーの予算であることを確認
	if budget.UserID() != domainUserID.Value() {
		return nil, errors.NewForbiddenError("この予算へのアクセス権限がありません")
	}

	// 更新フィールドの適用（名前フィールドは削除済み）

	if req.Amount != nil {
		currentAmount := budget.Amount()
		money, err := value.NewMoney(req.Amount.IntPart(), currentAmount.Currency())
		if err != nil {
			s.logger.Error("金額作成エラー",
				zap.Error(err),
				zap.String("amount", req.Amount.String()))
			return nil, errors.NewValidationError("無効な金額です")
		}
		if err := budget.UpdateAmount(*money); err != nil {
			s.logger.Error("金額更新エラー", zap.Error(err))
			return nil, errors.NewValidationError("金額の更新に失敗しました")
		}
	}

	if req.StartDate != nil || req.EndDate != nil {
		startDate := budget.StartDate()
		if req.StartDate != nil {
			startDate = *req.StartDate
		}

		endDate := budget.EndDate()
		if req.EndDate != nil {
			endDate = req.EndDate
		}

		if endDate != nil && endDate.Before(startDate) {
			return nil, errors.NewValidationError("終了日は開始日より後である必要があります")
		}

		if err := budget.UpdatePeriod(startDate, endDate); err != nil {
			s.logger.Error("期間更新エラー", zap.Error(err))
			return nil, errors.NewValidationError("期間の更新に失敗しました")
		}
	}

	if req.IsActive != nil {
		if *req.IsActive {
			budget.Activate()
		} else {
			budget.Deactivate()
		}
	}

	// リポジトリで更新
	if err := s.budgetRepo.Save(ctx, *budget); err != nil {
		s.logger.Error("予算更新エラー",
			zap.Error(err),
			zap.String("budgetID", budgetID.String()))
		return nil, errors.NewInternalError("予算情報の更新に失敗しました", err)
	}

	// DTOに変換して返却
	return dto.BudgetFromDomain(budget), nil
}

// DeleteBudget 予算を削除
func (s *BudgetService) DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから予算を取得
	budget, err := s.budgetRepo.FindByID(ctx, budgetID)
	if err != nil {
		s.logger.Error("予算取得エラー",
			zap.Error(err),
			zap.String("budgetID", budgetID.String()))
		return errors.NewInternalError("予算情報の取得に失敗しました", err)
	}

	if budget == nil {
		return errors.NewNotFoundError(fmt.Sprintf("予算が見つかりません: %s", budgetID))
	}

	// ユーザーの予算であることを確認
	if budget.UserID() != domainUserID.Value() {
		return errors.NewForbiddenError("この予算へのアクセス権限がありません")
	}

	// リポジトリから削除
	if err := s.budgetRepo.Delete(ctx, budgetID); err != nil {
		s.logger.Error("予算削除エラー",
			zap.Error(err),
			zap.String("budgetID", budgetID.String()))
		return errors.NewInternalError("予算の削除に失敗しました", err)
	}

	return nil
}

// GetActiveBudgetsByPeriod 期間内のアクティブな予算一覧を取得
func (s *BudgetService) GetActiveBudgetsByPeriod(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*dto.BudgetListResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリからユーザーの全予算を取得してフィルタリング
	allBudgets, err := s.budgetRepo.FindByUserID(ctx, domainUserID.Value())
	if err != nil {
		s.logger.Error("予算取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()))
		return nil, errors.NewInternalError("予算の取得に失敗しました", err)
	}

	// 期間内の予算をフィルタリング
	var budgets []budgetDomain.Budget
	for _, budget := range allBudgets {
		if budget.IsValidForDate(startDate) || budget.IsValidForDate(endDate) {
			budgets = append(budgets, budget)
		}
	}

	// アクティブな予算のみをフィルタリング
	var activeBudgets []*budgetDomain.Budget
	for i := range budgets {
		budget := &budgets[i]
		if budget.IsActive() {
			activeBudgets = append(activeBudgets, budget)
		}
	}

	// DTOに変換
	return &dto.BudgetListResponse{
		Budgets:    dto.BudgetsFromDomain(activeBudgets),
		TotalCount: int64(len(activeBudgets)),
	}, nil
}
