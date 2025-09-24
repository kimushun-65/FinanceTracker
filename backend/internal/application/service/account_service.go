package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"financetracker/internal/application/dto"
	accountDomain "financetracker/internal/domain/account/entity"
	accountRepo "financetracker/internal/domain/account/repository"
	accountValue "financetracker/internal/domain/account/value"
	"financetracker/internal/domain/common/value"
	userValue "financetracker/internal/domain/user/value"
	"financetracker/pkg/errors"
	"financetracker/pkg/logger"
)

// AccountService 口座管理サービス
type AccountService struct {
	accountRepo accountRepo.AccountRepository
	logger      *logger.Logger
}

// NewAccountService 新しいAccountServiceを作成
func NewAccountService(
	accountRepo accountRepo.AccountRepository,
	logger *logger.Logger,
) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		logger:      logger,
	}
}

// CreateAccount 口座を作成
func (s *AccountService) CreateAccount(ctx context.Context, userID uuid.UUID, req *dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// 口座名を作成
	accountName, err := accountValue.NewAccountName(req.Name)
	if err != nil {
		s.logger.Error("口座名作成エラー",
			zap.Error(err),
			zap.String("name", req.Name))
		return nil, errors.NewValidationError("無効な口座名です")
	}

	// 口座タイプを作成
	accountType, err := accountValue.NewAccountType(req.AccountType)
	if err != nil {
		s.logger.Error("口座タイプ作成エラー",
			zap.Error(err),
			zap.String("type", req.AccountType))
		return nil, errors.NewValidationError("無効な口座タイプです")
	}

	// 初期残高を設定（デフォルトは0）
	initialBalance := int64(0)
	if req.InitialBalance != nil {
		initialBalance = req.InitialBalance.IntPart()
	}

	// 金額を作成
	money, err := value.NewMoney(initialBalance, req.Currency)
	if err != nil {
		s.logger.Error("金額作成エラー",
			zap.Error(err),
			zap.Int64("amount", initialBalance),
			zap.String("currency", req.Currency))
		return nil, errors.NewValidationError("無効な金額です")
	}

	// 口座を作成
	account := accountDomain.NewAccount(*domainUserID, *accountName, *accountType, *money)

	// リポジトリに保存
	if err := s.accountRepo.Save(ctx, account); err != nil {
		s.logger.Error("口座保存エラー", zap.Error(err))
		return nil, errors.NewInternalError("口座の保存に失敗しました", err)
	}

	// DTOに変換して返却
	return dto.AccountFromDomain(account), nil
}

// GetAccount 口座情報を取得
func (s *AccountService) GetAccount(ctx context.Context, userID, accountID uuid.UUID) (*dto.AccountResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから口座を取得
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		s.logger.Error("口座取得エラー",
			zap.Error(err),
			zap.String("accountID", accountID.String()))
		return nil, errors.NewInternalError("口座情報の取得に失敗しました", err)
	}

	if account == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("口座が見つかりません: %s", accountID))
	}

	// ユーザーの口座であることを確認
	if account.UserID().Value() != domainUserID.Value() {
		return nil, errors.NewForbiddenError("この口座へのアクセス権限がありません")
	}

	// DTOに変換
	return dto.AccountFromDomain(account), nil
}

// GetAccountsByUser ユーザーの口座一覧を取得
func (s *AccountService) GetAccountsByUser(ctx context.Context, userID uuid.UUID) (*dto.AccountListResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから口座一覧を取得
	result, err := s.accountRepo.FindByUserID(ctx, domainUserID.Value(), nil)
	if err != nil {
		s.logger.Error("口座一覧取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()))
		return nil, errors.NewInternalError("口座一覧の取得に失敗しました", err)
	}

	accounts := result.Items

	// 合計残高を計算
	totalBalance := decimal.Zero
	for _, account := range accounts {
		balance := account.CurrentBalance()
		// 通貨が同じ場合のみ合計に加算（簡易実装）
		if balance.Currency() == "JPY" {
			totalBalance = totalBalance.Add(decimal.NewFromInt(balance.Amount()))
		}
	}

	// DTOに変換
	return &dto.AccountListResponse{
		Accounts:     dto.AccountsFromDomain(accounts),
		TotalCount:   result.TotalCount,
		TotalBalance: totalBalance,
	}, nil
}

// UpdateAccount 口座情報を更新
func (s *AccountService) UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, req *dto.UpdateAccountRequest) (*dto.AccountResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから口座を取得
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		s.logger.Error("口座取得エラー",
			zap.Error(err),
			zap.String("accountID", accountID.String()))
		return nil, errors.NewInternalError("口座情報の取得に失敗しました", err)
	}

	if account == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("口座が見つかりません: %s", accountID))
	}

	// ユーザーの口座であることを確認
	if account.UserID().Value() != domainUserID.Value() {
		return nil, errors.NewForbiddenError("この口座へのアクセス権限がありません")
	}

	// 更新フィールドの適用
	if req.Name != nil {
		accountName, err := accountValue.NewAccountName(*req.Name)
		if err != nil {
			s.logger.Error("口座名作成エラー",
				zap.Error(err),
				zap.String("name", *req.Name))
			return nil, errors.NewValidationError("無効な口座名です")
		}
		account.UpdateName(*accountName)
	}

	if req.AccountType != nil {
		// 口座タイプの変更は通常許可しないが、必要であればドメインエンティティにメソッドを追加
		s.logger.Warn("口座タイプの変更は現在サポートされていません",
			zap.String("accountID", accountID.String()),
			zap.String("newType", *req.AccountType))
		return nil, errors.NewValidationError("口座タイプの変更は現在サポートされていません")
	}

	// リポジトリで更新
	if err := s.accountRepo.Save(ctx, account); err != nil {
		s.logger.Error("口座更新エラー",
			zap.Error(err),
			zap.String("accountID", accountID.String()))
		return nil, errors.NewInternalError("口座情報の更新に失敗しました", err)
	}

	// DTOに変換して返却
	return dto.AccountFromDomain(account), nil
}

// DeleteAccount 口座を削除
func (s *AccountService) DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから口座を取得
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		s.logger.Error("口座取得エラー",
			zap.Error(err),
			zap.String("accountID", accountID.String()))
		return errors.NewInternalError("口座情報の取得に失敗しました", err)
	}

	if account == nil {
		return errors.NewNotFoundError(fmt.Sprintf("口座が見つかりません: %s", accountID))
	}

	// ユーザーの口座であることを確認
	if account.UserID().Value() != domainUserID.Value() {
		return errors.NewForbiddenError("この口座へのアクセス権限がありません")
	}

	// 残高が0でない場合は削除不可
	balance := account.CurrentBalance()
	if balance.Amount() != 0 {
		return errors.NewValidationError("残高がある口座は削除できません")
	}

	// リポジトリから削除
	if err := s.accountRepo.Delete(ctx, accountID); err != nil {
		s.logger.Error("口座削除エラー",
			zap.Error(err),
			zap.String("accountID", accountID.String()))
		return errors.NewInternalError("口座の削除に失敗しました", err)
	}

	return nil
}

// UpdateBalance 口座残高を更新
func (s *AccountService) UpdateBalance(ctx context.Context, userID, accountID uuid.UUID, amount decimal.Decimal, isDeposit bool) (*dto.AccountResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリから口座を取得
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		s.logger.Error("口座取得エラー",
			zap.Error(err),
			zap.String("accountID", accountID.String()))
		return nil, errors.NewInternalError("口座情報の取得に失敗しました", err)
	}

	if account == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("口座が見つかりません: %s", accountID))
	}

	// ユーザーの口座であることを確認
	if account.UserID().Value() != domainUserID.Value() {
		return nil, errors.NewForbiddenError("この口座へのアクセス権限がありません")
	}

	// 金額を作成
	balance := account.CurrentBalance()
	money, err := value.NewMoney(amount.IntPart(), balance.Currency())
	if err != nil {
		s.logger.Error("金額作成エラー",
			zap.Error(err),
			zap.String("amount", amount.String()),
			zap.String("currency", balance.Currency()))
		return nil, errors.NewValidationError("無効な金額です")
	}

	// 残高を更新
	if err := account.UpdateBalance(*money, isDeposit); err != nil {
		s.logger.Error("残高更新エラー", zap.Error(err))
		return nil, errors.NewValidationError("残高の更新に失敗しました")
	}

	// リポジトリで更新
	if err := s.accountRepo.Save(ctx, account); err != nil {
		s.logger.Error("口座更新エラー",
			zap.Error(err),
			zap.String("accountID", accountID.String()))
		return nil, errors.NewInternalError("口座情報の更新に失敗しました", err)
	}

	// DTOに変換して返却
	return dto.AccountFromDomain(account), nil
}
