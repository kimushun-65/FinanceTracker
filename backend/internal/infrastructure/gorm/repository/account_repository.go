package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	accountDomain "financetracker/internal/domain/account/entity"
	accountRepo "financetracker/internal/domain/account/repository"
	accountValue "financetracker/internal/domain/account/value"
	"financetracker/internal/domain/common"
	"financetracker/internal/domain/common/repository"
	"financetracker/internal/domain/common/value"
	userValue "financetracker/internal/domain/user/value"
	"financetracker/internal/infrastructure/gorm/model"
)

// AccountRepository 口座リポジトリの実装
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository 新しいAccountRepositoryを作成
func NewAccountRepository(db *gorm.DB) accountRepo.AccountRepository {
	return &AccountRepository{db: db}
}

// Save 口座を保存（新規作成または更新）
func (r *AccountRepository) Save(ctx context.Context, account *accountDomain.Account) error {
	accountModel := r.toModel(account)

	// IDが存在する場合は更新、存在しない場合は作成
	result := r.db.WithContext(ctx).Save(accountModel)
	if result.Error != nil {
		return fmt.Errorf("口座の保存に失敗しました: %w", result.Error)
	}

	return nil
}

// FindByID IDで口座を取得
func (r *AccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*accountDomain.Account, error) {
	var accountModel model.Account
	result := r.db.WithContext(ctx).First(&accountModel, "id = ?", id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("口座の取得に失敗しました: %w", result.Error)
	}

	return r.toDomain(&accountModel)
}

// Delete 口座を削除
func (r *AccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Account{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("口座の削除に失敗しました: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("削除対象の口座が見つかりません: %s", id)
	}

	return nil
}

// Exists 口座が存在するかチェック
func (r *AccountRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Account{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("口座の存在確認に失敗しました: %w", result.Error)
	}

	return count > 0, nil
}

// FindByUserID ユーザーIDで口座を検索
func (r *AccountRepository) FindByUserID(ctx context.Context, userID uuid.UUID, pagination *repository.Pagination, sorts ...*repository.Sort) (*repository.PagedResult[*accountDomain.Account], error) {
	var accountModels []model.Account
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// ソート設定
	if len(sorts) > 0 {
		for _, sort := range sorts {
			if sort != nil {
				orderClause := fmt.Sprintf("%s %s", sort.Field, sort.Order)
				query = query.Order(orderClause)
			}
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// ページネーション設定
	if pagination != nil {
		// 総件数を取得
		var totalCount int64
		countQuery := r.db.WithContext(ctx).Model(&model.Account{}).Where("user_id = ?", userID)
		if err := countQuery.Count(&totalCount).Error; err != nil {
			return nil, fmt.Errorf("口座数の取得に失敗しました: %w", err)
		}

		// ページネーション適用
		query = query.Limit(pagination.PageSize).Offset(pagination.Offset)

		// データ取得
		if err := query.Find(&accountModels).Error; err != nil {
			return nil, fmt.Errorf("口座一覧の取得に失敗しました: %w", err)
		}

		// ドメインモデルに変換
		accounts := make([]*accountDomain.Account, len(accountModels))
		for i, accountModel := range accountModels {
			account, err := r.toDomain(&accountModel)
			if err != nil {
				return nil, err
			}
			accounts[i] = account
		}

		return repository.NewPagedResult(accounts, totalCount, pagination), nil
	}

	// ページネーションなしの場合
	if err := query.Find(&accountModels).Error; err != nil {
		return nil, fmt.Errorf("口座一覧の取得に失敗しました: %w", err)
	}

	// ドメインモデルに変換
	accounts := make([]*accountDomain.Account, len(accountModels))
	for i, accountModel := range accountModels {
		account, err := r.toDomain(&accountModel)
		if err != nil {
			return nil, err
		}
		accounts[i] = account
	}

	return &repository.PagedResult[*accountDomain.Account]{
		Items:      accounts,
		TotalCount: int64(len(accounts)),
		Page:       1,
		PageSize:   len(accounts),
		TotalPages: 1,
		HasNext:    false,
		HasPrev:    false,
	}, nil
}

// FindByUserIDAndID ユーザーIDとエンティティIDで検索
func (r *AccountRepository) FindByUserIDAndID(ctx context.Context, userID, entityID uuid.UUID) (*accountDomain.Account, error) {
	var accountModel model.Account
	result := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, entityID).First(&accountModel)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("口座の取得に失敗しました: %w", result.Error)
	}

	return r.toDomain(&accountModel)
}

// DeleteByUserIDAndID ユーザーIDとエンティティIDで削除
func (r *AccountRepository) DeleteByUserIDAndID(ctx context.Context, userID, entityID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, entityID).Delete(&model.Account{})
	if result.Error != nil {
		return fmt.Errorf("口座の削除に失敗しました: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("削除対象の口座が見つかりません")
	}

	return nil
}

// CountByUserID ユーザーのエンティティ数をカウント
func (r *AccountRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Account{}).Where("user_id = ?", userID).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("口座数の取得に失敗しました: %w", result.Error)
	}

	return count, nil
}

// FindByUserIDAndType ユーザーIDと口座タイプで検索
func (r *AccountRepository) FindByUserIDAndType(ctx context.Context, userID userValue.UserID, accountType accountValue.AccountType) ([]*accountDomain.Account, error) {
	var accountModels []model.Account
	result := r.db.WithContext(ctx).Where("user_id = ? AND account_type = ?", userID.Value(), accountType.Value()).Find(&accountModels)

	if result.Error != nil {
		return nil, fmt.Errorf("口座の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	accounts := make([]*accountDomain.Account, len(accountModels))
	for i, accountModel := range accountModels {
		account, err := r.toDomain(&accountModel)
		if err != nil {
			return nil, err
		}
		accounts[i] = account
	}

	return accounts, nil
}

// GetTotalBalanceByUserID ユーザーの総残高を取得（通貨別）
func (r *AccountRepository) GetTotalBalanceByUserID(ctx context.Context, userID userValue.UserID) (map[string]value.Money, error) {
	type BalanceResult struct {
		Currency string
		Total    int64
	}

	var results []BalanceResult
	err := r.db.WithContext(ctx).Model(&model.Account{}).
		Select("currency, SUM(balance) as total").
		Where("user_id = ?", userID.Value()).
		Group("currency").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("総残高の取得に失敗しました: %w", err)
	}

	balances := make(map[string]value.Money)
	for _, result := range results {
		money, err := value.NewMoney(result.Total, result.Currency)
		if err != nil {
			return nil, fmt.Errorf("金額の作成に失敗しました: %w", err)
		}
		balances[result.Currency] = *money
	}

	return balances, nil
}

// GetTotalAssetsByUserID ユーザーの総資産を取得
func (r *AccountRepository) GetTotalAssetsByUserID(ctx context.Context, userID userValue.UserID) (value.Money, error) {
	// 簡易実装：JPY通貨のみをサポート
	var total int64
	err := r.db.WithContext(ctx).Model(&model.Account{}).
		Where("user_id = ? AND currency = ?", userID.Value(), "JPY").
		Select("COALESCE(SUM(balance), 0)").
		Scan(&total).Error

	if err != nil {
		return value.Money{}, fmt.Errorf("総資産の取得に失敗しました: %w", err)
	}

	money, err := value.NewMoney(total, "JPY")
	if err != nil {
		return value.Money{}, fmt.Errorf("金額の作成に失敗しました: %w", err)
	}

	return *money, nil
}

// toModel ドメインモデルからGORMモデルへの変換
func (r *AccountRepository) toModel(account *accountDomain.Account) *model.Account {
	balance := account.CurrentBalance()

	return &model.Account{
		Base: model.Base{
			ID:        account.ID,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		},
		UserID:   account.UserID().Value(),
		Name:     account.Name().String(),
		Type:     model.AccountType(account.Type().Value()),
		Balance:  decimal.NewFromInt(balance.Amount()),
		Currency: balance.Currency(),
		IsActive: true,
	}
}

// toDomain GORMモデルからドメインモデルへの変換
func (r *AccountRepository) toDomain(accountModel *model.Account) (*accountDomain.Account, error) {
	// ユーザーID
	userID := userValue.NewUserID(accountModel.UserID)

	// 口座名
	accountName, err := accountValue.NewAccountName(accountModel.Name)
	if err != nil {
		return nil, fmt.Errorf("口座名の作成に失敗しました: %w", err)
	}

	// 口座タイプ（データベースのタイプをドメインタイプにマッピング）
	var domainAccountType string
	switch accountModel.Type {
	case model.AccountTypeCash:
		domainAccountType = "cash"
	case model.AccountTypeBank:
		domainAccountType = "checking"
	case model.AccountTypeInvestment:
		domainAccountType = "investment"
	case model.AccountTypeCreditCard:
		domainAccountType = "checking" // クレジットカードは当座預金として扱う
	case model.AccountTypeLoan:
		domainAccountType = "checking" // ローンは当座預金として扱う
	case model.AccountTypeOther:
		domainAccountType = "checking" // その他は当座預金として扱う
	default:
		domainAccountType = "checking" // デフォルトは当座預金
	}

	accountType, err := accountValue.NewAccountType(domainAccountType)
	if err != nil {
		return nil, fmt.Errorf("口座タイプの作成に失敗しました: %w", err)
	}

	// 残高（初期残高と現在残高を同じ値として扱う）
	money, err := value.NewMoney(accountModel.Balance.IntPart(), accountModel.Currency)
	if err != nil {
		return nil, fmt.Errorf("残高の作成に失敗しました: %w", err)
	}

	// 残高オブジェクト（初期残高と現在残高を同じ値として扱う）
	balance, err := accountValue.NewBalance(*money, *money)
	if err != nil {
		return nil, fmt.Errorf("残高オブジェクトの作成に失敗しました: %w", err)
	}

	// BaseEntity
	baseEntity := common.BaseEntity{
		ID:        accountModel.ID,
		CreatedAt: accountModel.CreatedAt,
		UpdatedAt: accountModel.UpdatedAt,
	}

	// ドメインエンティティを再構築
	account := accountDomain.ReconstructAccount(
		baseEntity,
		*userID,
		*accountName,
		*accountType,
		*balance,
	)

	return account, nil
}
