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
	accountService  *AccountService
	categoryService *CategoryService
	logger          *logger.Logger
}

// NewTransactionService 新しいTransactionServiceを作成
func NewTransactionService(
	transactionRepo transactionRepo.TransactionRepository,
	accountService *AccountService,
	categoryService *CategoryService,
	logger *logger.Logger,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountService:  accountService,
		categoryService: categoryService,
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

	// 口座情報を取得して口座タイプを確認
	account, err := s.accountService.GetAccount(ctx, userID, req.AccountID)
	if err != nil {
		s.logger.Error("口座取得エラー",
			zap.Error(err),
			zap.String("accountID", req.AccountID.String()))

		// 取引の保存は成功したが口座取得が失敗した場合、取引を削除してロールバック
		if deleteErr := s.transactionRepo.Delete(ctx, transaction.GetID()); deleteErr != nil {
			s.logger.Error("取引ロールバック失敗", zap.Error(deleteErr))
		}
		return nil, errors.NewInternalError("口座情報の取得に失敗しました", err)
	}

	// 口座残高を更新
	amount := decimal.NewFromInt(money.Amount())
	isDeposit := transactionType.IsIncome()

	// クレジットカード口座の場合も通常の論理を使用
	// クレジットカードの場合:
	// - 支出(expense): isDeposit = false → 残高減少（債務増加、よりマイナスに）
	// - 収入(income/payment): isDeposit = true → 残高増加（債務減少、よりプラスに）
	isCreditCard := account.AccountType == "credit_card"

	_, err = s.accountService.UpdateBalance(ctx, userID, req.AccountID, amount, isDeposit)
	if err != nil {
		s.logger.Error("残高更新エラー",
			zap.Error(err),
			zap.String("accountID", req.AccountID.String()),
			zap.String("amount", amount.String()),
			zap.Bool("isDeposit", isDeposit),
			zap.Bool("isCreditCard", isCreditCard))

		// 取引の保存は成功したが残高更新が失敗した場合、取引を削除してロールバック
		if deleteErr := s.transactionRepo.Delete(ctx, transaction.GetID()); deleteErr != nil {
			s.logger.Error("取引ロールバック失敗", zap.Error(deleteErr))
		}
		return nil, errors.NewInternalError("口座残高の更新に失敗しました", err)
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

	// リポジトリから取引一覧を取得（新しいフィルタ付きメソッドを使用）
	offset := (params.Page - 1) * params.PerPage
	result, err := s.transactionRepo.FindByUserIDWithFilters(
		ctx,
		domainUserID.Value(),
		params.CategoryID,
		params.DateFrom,
		params.DateTo,
		params.PerPage,
		offset,
	)
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

	// 残高更新のため、変更前の情報を保存
	oldAmount := transaction.Amount()
	oldAmountDecimal := decimal.NewFromInt(oldAmount.Amount())
	oldType := transaction.Type()
	accountID := transaction.AccountID()

	// 更新フィールドの適用
	if req.CategoryID != nil {
		if err := transaction.UpdateCategory(*req.CategoryID); err != nil {
			s.logger.Error("カテゴリー変更エラー", zap.Error(err))
			return nil, errors.NewValidationError("カテゴリーの変更に失敗しました")
		}
	}

	var newAmountDecimal decimal.Decimal
	amountChanged := false
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
		newAmountDecimal = decimal.NewFromInt(money.Amount())
		amountChanged = true
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

	// 金額が変更された場合は口座残高を更新
	if amountChanged {
		// 口座情報を取得してクレジットカードかどうか確認
		account, err := s.accountService.GetAccount(ctx, userID, accountID)
		if err != nil {
			s.logger.Error("口座取得エラー", zap.Error(err))
			return nil, errors.NewInternalError("口座情報の取得に失敗しました", err)
		}
		isCreditCard := account.AccountType == "credit_card"

		// 古い金額の影響を取り消す（逆操作）
		oldIsDeposit := oldType.IsIncome()
		// クレジットカードの場合も通常の論理を使用（逆操作で取り消し）
		_, err = s.accountService.UpdateBalance(ctx, userID, accountID, oldAmountDecimal, !oldIsDeposit)
		if err != nil {
			s.logger.Error("旧残高取り消しエラー",
				zap.Error(err),
				zap.String("accountID", accountID.String()),
				zap.String("oldAmount", oldAmountDecimal.String()),
				zap.Bool("reverseOperation", !oldIsDeposit),
				zap.Bool("isCreditCard", isCreditCard))
			return nil, errors.NewInternalError("口座残高の更新に失敗しました", err)
		}

		// 新しい金額を適用
		newIsDeposit := transaction.Type().IsIncome()
		// クレジットカードの場合も通常の論理を使用
		_, err = s.accountService.UpdateBalance(ctx, userID, accountID, newAmountDecimal, newIsDeposit)
		if err != nil {
			s.logger.Error("新残高適用エラー",
				zap.Error(err),
				zap.String("accountID", accountID.String()),
				zap.String("newAmount", newAmountDecimal.String()),
				zap.Bool("isDeposit", newIsDeposit),
				zap.Bool("isCreditCard", isCreditCard))

			// 新しい金額の適用に失敗した場合、古い金額を再適用
			oldIsDepositForRevert := oldType.IsIncome()
			if _, revertErr := s.accountService.UpdateBalance(ctx, userID, accountID, oldAmountDecimal, oldIsDepositForRevert); revertErr != nil {
				s.logger.Error("残高復元エラー", zap.Error(revertErr))
			}
			return nil, errors.NewInternalError("口座残高の更新に失敗しました", err)
		}
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

	// 残高更新のため、削除前の情報を保存
	amount := transaction.Amount()
	amountDecimal := decimal.NewFromInt(amount.Amount())
	transactionType := transaction.Type()
	accountID := transaction.AccountID()

	// リポジトリから削除
	if err := s.transactionRepo.Delete(ctx, transactionID); err != nil {
		s.logger.Error("取引削除エラー",
			zap.Error(err),
			zap.String("transactionID", transactionID.String()))
		return errors.NewInternalError("取引の削除に失敗しました", err)
	}

	// 口座情報を取得してクレジットカードかどうか確認
	account, err := s.accountService.GetAccount(ctx, userID, accountID)
	if err != nil {
		s.logger.Error("口座取得エラー", zap.Error(err))
		return errors.NewInternalError("口座情報の取得に失敗しました", err)
	}
	isCreditCard := account.AccountType == "credit_card"

	// 削除された取引の影響を口座残高から取り消す（逆操作）
	isDeposit := transactionType.IsIncome()
	// クレジットカードの場合も通常の論理を使用（逆操作で取り消し）
	_, err = s.accountService.UpdateBalance(ctx, userID, accountID, amountDecimal, !isDeposit)
	if err != nil {
		s.logger.Error("削除時残高更新エラー",
			zap.Error(err),
			zap.String("accountID", accountID.String()),
			zap.String("amount", amountDecimal.String()),
			zap.Bool("reverseOperation", !isDeposit),
			zap.Bool("isCreditCard", isCreditCard))
		return errors.NewInternalError("口座残高の更新に失敗しました", err)
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

// GetCategorySummary カテゴリー別サマリーを取得
func (s *TransactionService) GetCategorySummary(
	ctx context.Context,
	userID uuid.UUID,
	startDate, endDate time.Time,
	transactionType string, // "expense", "income", "all"
) (*dto.CategorySummaryResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリからトランザクションを取得
	transactions, err := s.transactionRepo.FindByUserIDAndDateRange(
		ctx, domainUserID.Value(), startDate, endDate,
	)
	if err != nil {
		s.logger.Error("カテゴリー別サマリー取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()),
			zap.Time("startDate", startDate),
			zap.Time("endDate", endDate))
		return nil, errors.NewInternalError("カテゴリー別サマリーの取得に失敗しました", err)
	}

	// タイプでフィルタ
	var filteredTransactions []transactionDomain.Transaction
	if transactionType != "all" {
		for _, tx := range transactions {
			if transactionType == "expense" && !tx.Type().IsIncome() {
				filteredTransactions = append(filteredTransactions, tx)
			} else if transactionType == "income" && tx.Type().IsIncome() {
				filteredTransactions = append(filteredTransactions, tx)
			}
		}
	} else {
		filteredTransactions = transactions
	}

	// カテゴリー別に集計
	categoryMap := make(map[uuid.UUID]*dto.CategorySummaryDetail)
	totalAmount := decimal.Zero

	for _, tx := range filteredTransactions {
		categoryID := tx.CategoryID()
		amount := decimal.NewFromInt(tx.Amount().Amount())

		if _, exists := categoryMap[categoryID]; !exists {
			// カテゴリー情報を取得
			category, err := s.categoryService.GetCategory(ctx, userID, categoryID)
			if err != nil {
				s.logger.Warn("カテゴリー情報取得エラー（スキップします）",
					zap.Error(err),
					zap.String("categoryID", categoryID.String()))
				continue
			}

			categoryMap[categoryID] = &dto.CategorySummaryDetail{
				CategoryID:       categoryID,
				CategoryName:     category.Name,
				CategoryIcon:     category.Icon,
				TotalAmount:      decimal.Zero,
				TransactionCount: 0,
				Percentage:       decimal.Zero,
			}
		}

		categoryMap[categoryID].TotalAmount = categoryMap[categoryID].TotalAmount.Add(amount)
		categoryMap[categoryID].TransactionCount++
		totalAmount = totalAmount.Add(amount)
	}

	// パーセンテージを計算
	for _, detail := range categoryMap {
		if !totalAmount.IsZero() {
			detail.Percentage = detail.TotalAmount.Div(totalAmount).Mul(decimal.NewFromInt(100))
		}
	}

	// 配列に変換（金額の大きい順にソート）
	byCategory := make([]dto.CategorySummaryDetail, 0, len(categoryMap))
	for _, detail := range categoryMap {
		byCategory = append(byCategory, *detail)
	}

	// 金額の大きい順にソート
	sort.Slice(byCategory, func(i, j int) bool {
		return byCategory[i].TotalAmount.GreaterThan(byCategory[j].TotalAmount)
	})

	return &dto.CategorySummaryResponse{
		Period: dto.PeriodInfo{
			From: startDate,
			To:   endDate,
		},
		TotalAmount: totalAmount,
		ByCategory:  byCategory,
	}, nil
}
