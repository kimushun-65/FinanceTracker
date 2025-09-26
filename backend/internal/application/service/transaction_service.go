package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"financetracker/internal/application/dto"
	"financetracker/internal/domain/common/value"
	transactionDomain "financetracker/internal/domain/transaction/entity"
	transactionRepo "financetracker/internal/domain/transaction/repository"
	transactionValue "financetracker/internal/domain/transaction/value"
	userValue "financetracker/internal/domain/user/value"
	"financetracker/pkg/errors"
	"financetracker/pkg/logger"
)

// TransactionService 取引管理サービス
type TransactionService struct {
	transactionRepo transactionRepo.TransactionRepository
	logger          *logger.Logger
}

// NewTransactionService 新しいTransactionServiceを作成
func NewTransactionService(
	transactionRepo transactionRepo.TransactionRepository,
	logger *logger.Logger,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		logger:          logger,
	}
}

// CreateTransaction 取引を作成
func (s *TransactionService) CreateTransaction(ctx context.Context, userID uuid.UUID, req *dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// 取引タイプを作成
	transactionType, err := transactionValue.NewTransactionType(req.TransactionType)
	if err != nil {
		s.logger.Error("取引タイプ作成エラー",
			zap.Error(err),
			zap.String("type", req.TransactionType))
		return nil, errors.NewValidationError("無効な取引タイプです")
	}

	// 金額を作成（デフォルトでJPY）
	money, err := value.NewMoney(req.Amount.IntPart(), "JPY")
	if err != nil {
		s.logger.Error("金額作成エラー",
			zap.Error(err),
			zap.String("amount", req.Amount.String()))
		return nil, errors.NewValidationError("無効な金額です")
	}

	// 説明を作成
	description, err := transactionValue.NewDescription(req.Description)
	if err != nil {
		s.logger.Error("説明作成エラー",
			zap.Error(err),
			zap.String("description", req.Description))
		return nil, errors.NewValidationError("無効な説明です")
	}

	// 取引日を設定（デフォルトは現在時刻）
	transactionDate := time.Now()
	if req.TransactionDate != nil {
		transactionDate = *req.TransactionDate
	}

	// 取引を作成
	transaction, err := transactionDomain.NewTransaction(
		domainUserID.Value(),
		req.AccountID,
		req.CategoryID,
		*money,
		transactionType,
		transactionDate,
		description,
	)
	if err != nil {
		s.logger.Error("取引作成エラー", zap.Error(err))
		return nil, errors.NewValidationError("取引の作成に失敗しました")
	}

	// リポジトリに保存
	if err := s.transactionRepo.Save(ctx, transaction); err != nil {
		s.logger.Error("取引保存エラー", zap.Error(err))
		return nil, errors.NewInternalError("取引の保存に失敗しました", err)
	}

	// DTOに変換して返却
	return dto.TransactionFromDomain(&transaction), nil
}

// GetTransaction 取引情報を取得
func (s *TransactionService) GetTransaction(ctx context.Context, userID, transactionID uuid.UUID) (*dto.TransactionResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから取引を取得
	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		s.logger.Error("取引取得エラー",
			zap.Error(err),
			zap.String("transactionID", transactionID.String()))
		return nil, errors.NewInternalError("取引情報の取得に失敗しました", err)
	}

	if transaction == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("取引が見つかりません: %s", transactionID))
	}

	// ユーザーの取引であることを確認
	if transaction.UserID() != domainUserID.Value() {
		return nil, errors.NewForbiddenError("この取引へのアクセス権限がありません")
	}

	// DTOに変換
	return dto.TransactionFromDomain(transaction), nil
}

// GetTransactionsByUser ユーザーの取引一覧を取得
func (s *TransactionService) GetTransactionsByUser(ctx context.Context, userID uuid.UUID, params *dto.TransactionSearchParams) (*dto.TransactionListResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから取引一覧を取得
	var result []transactionDomain.Transaction
	var err error

	// 日付フィルタが指定されている場合は日付範囲で検索
	if params.DateFrom != nil || params.DateTo != nil {
		// デフォルト日付を設定（フィルタが片方だけの場合）
		startDate := time.Time{}
		endDate := time.Now()

		if params.DateFrom != nil {
			startDate = *params.DateFrom
		}
		if params.DateTo != nil {
			endDate = *params.DateTo
		}

		result, err = s.transactionRepo.FindByUserIDAndDateRange(ctx, domainUserID.Value(), startDate, endDate)

		// 手動でページネーションを適用
		if err == nil {
			offset := (params.Page - 1) * params.PerPage
			if offset >= len(result) {
				result = []transactionDomain.Transaction{}
			} else {
				end := offset + params.PerPage
				if end > len(result) {
					end = len(result)
				}
				result = result[offset:end]
			}
		}
	} else {
		// 日付フィルタがない場合は通常の検索
		offset := (params.Page - 1) * params.PerPage
		result, err = s.transactionRepo.FindByUserID(ctx, domainUserID.Value(), params.PerPage, offset)
	}
	if err != nil {
		s.logger.Error("取引一覧取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()))
		return nil, errors.NewInternalError("取引一覧の取得に失敗しました", err)
	}

	// 収入と支出の合計を計算
	totalIncome := decimal.Zero
	totalExpense := decimal.Zero
	var transactionPtrs []*transactionDomain.Transaction
	for i := range result {
		transactionPtrs = append(transactionPtrs, &result[i])
		amount := result[i].Amount()
		if result[i].Type().IsIncome() {
			totalIncome = totalIncome.Add(decimal.NewFromInt(amount.Amount()))
		} else {
			totalExpense = totalExpense.Add(decimal.NewFromInt(amount.Amount()))
		}
	}

	// 総件数を取得
	totalCount, err := s.transactionRepo.CountByUserID(ctx, domainUserID.Value())
	if err != nil {
		totalCount = int64(len(result))
	}

	// DTOに変換
	return &dto.TransactionListResponse{
		Transactions: dto.TransactionsFromDomain(transactionPtrs),
		TotalCount:   totalCount,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Page:         params.Page,
		PerPage:      params.PerPage,
	}, nil
}

// GetTransactionsByAccount 口座別の取引一覧を取得
func (s *TransactionService) GetTransactionsByAccount(ctx context.Context, userID, accountID uuid.UUID, params *dto.TransactionSearchParams) (*dto.TransactionListResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから取引一覧を取得
	offset := (params.Page - 1) * params.PerPage
	result, err := s.transactionRepo.FindByAccountID(ctx, accountID, params.PerPage, offset)
	if err != nil {
		s.logger.Error("口座別取引一覧取得エラー",
			zap.Error(err),
			zap.String("accountID", accountID.String()))
		return nil, errors.NewInternalError("取引一覧の取得に失敗しました", err)
	}

	// ユーザーの口座に属する取引のみをフィルタリング
	var filteredTransactions []*transactionDomain.Transaction
	totalIncome := decimal.Zero
	totalExpense := decimal.Zero
	for i := range result {
		if result[i].UserID() == domainUserID.Value() {
			filteredTransactions = append(filteredTransactions, &result[i])
			amount := result[i].Amount()
			if result[i].Type().IsIncome() {
				totalIncome = totalIncome.Add(decimal.NewFromInt(amount.Amount()))
			} else {
				totalExpense = totalExpense.Add(decimal.NewFromInt(amount.Amount()))
			}
		}
	}

	// DTOに変換
	return &dto.TransactionListResponse{
		Transactions: dto.TransactionsFromDomain(filteredTransactions),
		TotalCount:   int64(len(filteredTransactions)),
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Page:         params.Page,
		PerPage:      params.PerPage,
	}, nil
}

// UpdateTransaction 取引情報を更新
func (s *TransactionService) UpdateTransaction(ctx context.Context, userID, transactionID uuid.UUID, req *dto.UpdateTransactionRequest) (*dto.TransactionResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから取引を取得
	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		s.logger.Error("取引取得エラー",
			zap.Error(err),
			zap.String("transactionID", transactionID.String()))
		return nil, errors.NewInternalError("取引情報の取得に失敗しました", err)
	}

	if transaction == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("取引が見つかりません: %s", transactionID))
	}

	// ユーザーの取引であることを確認
	if transaction.UserID() != domainUserID.Value() {
		return nil, errors.NewForbiddenError("この取引へのアクセス権限がありません")
	}

	// 更新フィールドの適用
	if req.CategoryID != nil {
		// カテゴリー変更の実装が必要な場合は、ドメインエンティティにメソッドを追加
		s.logger.Warn("カテゴリー変更は現在サポートされていません",
			zap.String("transactionID", transactionID.String()))
		return nil, errors.NewValidationError("カテゴリーの変更は現在サポートされていません")
	}

	if req.Amount != nil {
		currentAmount := transaction.Amount()
		money, err := value.NewMoney(req.Amount.IntPart(), currentAmount.Currency())
		if err != nil {
			s.logger.Error("金額作成エラー",
				zap.Error(err),
				zap.String("amount", req.Amount.String()))
			return nil, errors.NewValidationError("無効な金額です")
		}
		if err := transaction.UpdateAmount(*money); err != nil {
			s.logger.Error("金額変更エラー", zap.Error(err))
			return nil, errors.NewValidationError("金額の変更に失敗しました")
		}
	}

	if req.Description != nil {
		description, err := transactionValue.NewDescription(*req.Description)
		if err != nil {
			s.logger.Error("説明作成エラー",
				zap.Error(err),
				zap.String("description", *req.Description))
			return nil, errors.NewValidationError("無効な説明です")
		}
		transaction.UpdateDescription(description)
	}

	if req.TransactionDate != nil {
		if err := transaction.UpdateDate(*req.TransactionDate); err != nil {
			s.logger.Error("取引日変更エラー", zap.Error(err))
			return nil, errors.NewValidationError("取引日の変更に失敗しました")
		}
	}

	// リポジトリで更新
	if err := s.transactionRepo.Save(ctx, *transaction); err != nil {
		s.logger.Error("取引更新エラー",
			zap.Error(err),
			zap.String("transactionID", transactionID.String()))
		return nil, errors.NewInternalError("取引情報の更新に失敗しました", err)
	}

	// DTOに変換して返却
	return dto.TransactionFromDomain(transaction), nil
}

// DeleteTransaction 取引を削除
func (s *TransactionService) DeleteTransaction(ctx context.Context, userID, transactionID uuid.UUID) error {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから取引を取得
	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		s.logger.Error("取引取得エラー",
			zap.Error(err),
			zap.String("transactionID", transactionID.String()))
		return errors.NewInternalError("取引情報の取得に失敗しました", err)
	}

	if transaction == nil {
		return errors.NewNotFoundError(fmt.Sprintf("取引が見つかりません: %s", transactionID))
	}

	// ユーザーの取引であることを確認
	if transaction.UserID() != domainUserID.Value() {
		return errors.NewForbiddenError("この取引へのアクセス権限がありません")
	}

	// リポジトリから削除
	if err := s.transactionRepo.Delete(ctx, transactionID); err != nil {
		s.logger.Error("取引削除エラー",
			zap.Error(err),
			zap.String("transactionID", transactionID.String()))
		return errors.NewInternalError("取引の削除に失敗しました", err)
	}

	return nil
}

// GetMonthlyTransactionSummary 月次取引サマリーを取得
func (s *TransactionService) GetMonthlyTransactionSummary(ctx context.Context, userID uuid.UUID, year, month int) (*dto.MonthlyTransactionSummary, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// 月初と月末の日付を設定
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// リポジトリから月次データを取得
	transactions, err := s.transactionRepo.FindByUserIDAndDateRange(ctx, domainUserID.Value(), startDate, endDate)
	if err != nil {
		s.logger.Error("月次サマリー取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()),
			zap.Int("year", year),
			zap.Int("month", month))
		return nil, errors.NewInternalError("月次サマリーの取得に失敗しました", err)
	}

	// 収入と支出の合計を計算
	totalIncome := decimal.Zero
	totalExpense := decimal.Zero
	
	// 日別データを格納するマップ
	dailyDataMap := make(map[string]*dto.DailyTransactionData)
	
	for _, transaction := range transactions {
		amount := transaction.Amount()
		amountDecimal := decimal.NewFromInt(amount.Amount())
		
		// 月全体の合計計算
		if transaction.Type().IsIncome() {
			totalIncome = totalIncome.Add(amountDecimal)
		} else {
			totalExpense = totalExpense.Add(amountDecimal)
		}
		
		// 日別データの集計
		dateKey := transaction.Date().Format("2006-01-02")
		if dailyData, exists := dailyDataMap[dateKey]; exists {
			// 既存の日付データを更新
			if transaction.Type().IsIncome() {
				dailyData.TotalIncome = dailyData.TotalIncome.Add(amountDecimal)
			} else {
				dailyData.TotalExpense = dailyData.TotalExpense.Add(amountDecimal)
			}
			dailyData.NetAmount = dailyData.TotalIncome.Sub(dailyData.TotalExpense)
			dailyData.Count++
		} else {
			// 新しい日付データを作成
			newDailyData := &dto.DailyTransactionData{
				Date:         transaction.Date(),
				TotalIncome:  decimal.Zero,
				TotalExpense: decimal.Zero,
				Count:        1,
			}
			
			if transaction.Type().IsIncome() {
				newDailyData.TotalIncome = amountDecimal
			} else {
				newDailyData.TotalExpense = amountDecimal
			}
			newDailyData.NetAmount = newDailyData.TotalIncome.Sub(newDailyData.TotalExpense)
			
			dailyDataMap[dateKey] = newDailyData
		}
	}
	
	// マップをスライスに変換してソート
	dailyDataSlice := make([]dto.DailyTransactionData, 0, len(dailyDataMap))
	for _, dailyData := range dailyDataMap {
		dailyDataSlice = append(dailyDataSlice, *dailyData)
	}
	
	// 日付順にソート
	sort.Slice(dailyDataSlice, func(i, j int) bool {
		return dailyDataSlice[i].Date.Before(dailyDataSlice[j].Date)
	})

	// DTOに変換
	return &dto.MonthlyTransactionSummary{
		Year:         year,
		Month:        month,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		NetAmount:    totalIncome.Sub(totalExpense),
		DailyData:    dailyDataSlice,
	}, nil
}
