package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"financetracker/internal/application/dto"
	accountRepo "financetracker/internal/domain/account/repository"
	assetDomain "financetracker/internal/domain/asset/entity"
	assetRepo "financetracker/internal/domain/asset/repository"
	assetValue "financetracker/internal/domain/asset/value"
	commonValue "financetracker/internal/domain/common/value"
	"financetracker/pkg/errors"
	"financetracker/pkg/logger"
)

// AssetService 資産サービス
type AssetService struct {
	assetSnapshotRepo assetRepo.AssetSnapshotRepository
	accountRepo       accountRepo.AccountRepository
	logger            *logger.Logger
}

// NewAssetService 新しいAssetServiceを作成
func NewAssetService(
	assetSnapshotRepo assetRepo.AssetSnapshotRepository,
	accountRepo accountRepo.AccountRepository,
	logger *logger.Logger,
) *AssetService {
	return &AssetService{
		assetSnapshotRepo: assetSnapshotRepo,
		accountRepo:       accountRepo,
		logger:            logger,
	}
}

// GetAssetSnapshots 資産スナップショット一覧を取得
func (s *AssetService) GetAssetSnapshots(
	ctx context.Context,
	userID uuid.UUID,
	startDate, endDate time.Time,
) (*dto.AssetSnapshotListResponse, error) {
	// リポジトリからスナップショットを取得
	snapshots, err := s.assetSnapshotRepo.FindByUserIDAndDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		s.logger.Error("資産スナップショット取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()),
			zap.Time("startDate", startDate),
			zap.Time("endDate", endDate))
		return nil, errors.NewInternalError("資産スナップショットの取得に失敗しました", err)
	}

	// DTOに変換
	response := &dto.AssetSnapshotListResponse{
		Snapshots:  dto.AssetSnapshotsFromDomain(snapshots),
		TotalCount: len(snapshots),
	}

	return response, nil
}

// CalculateCurrentAssetSnapshot 現在の資産スナップショットを計算
// トランザクション履歴から指定日時点の資産状況を計算
func (s *AssetService) CalculateCurrentAssetSnapshot(
	ctx context.Context,
	userID uuid.UUID,
	targetDate time.Time,
) (*dto.AssetSnapshotResponse, error) {
	// 全口座を取得
	result, err := s.accountRepo.FindByUserID(ctx, userID, nil)
	if err != nil {
		s.logger.Error("口座一覧取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()))
		return nil, errors.NewInternalError("口座一覧の取得に失敗しました", err)
	}

	accounts := result.Items

	// 各口座の残高を取得
	accountBalances := make([]assetValue.AccountBalance, 0, len(accounts))
	totalAssets := decimal.Zero

	for _, account := range accounts {
		balance := account.CurrentBalance()
		balanceDecimal := decimal.NewFromInt(balance.Amount())

		// クレジットカード口座は負債なので資産から除外
		if !account.Type().IsCreditCard() {
			totalAssets = totalAssets.Add(balanceDecimal)
		}

		accountBalances = append(accountBalances, assetValue.AccountBalance{
			AccountID:   account.ID,
			AccountName: account.Name().String(),
			Balance:     balance,
		})
	}

	// DTOに変換
	accountResponses := make([]dto.AccountBalanceResponse, len(accountBalances))
	for i, ab := range accountBalances {
		accountResponses[i] = dto.AccountBalanceResponse{
			AccountID:   ab.AccountID,
			AccountName: ab.AccountName,
			Balance:     decimal.NewFromInt(ab.Balance.Amount()),
			Currency:    ab.Balance.Currency(),
		}
	}

	response := &dto.AssetSnapshotResponse{
		ID:           uuid.New(), // 計算時は新しいIDを生成
		UserID:       userID,
		SnapshotDate: targetDate,
		TotalAssets:  totalAssets,
		Accounts:     accountResponses,
		CreatedAt:    time.Now(),
	}

	return response, nil
}

// GetLatestAssetSnapshot 最新の資産スナップショットを取得
func (s *AssetService) GetLatestAssetSnapshot(
	ctx context.Context,
	userID uuid.UUID,
) (*dto.AssetSnapshotResponse, error) {
	// リポジトリから最新のスナップショットを取得
	snapshot, err := s.assetSnapshotRepo.FindLatestByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("最新資産スナップショット取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()))
		return nil, errors.NewInternalError("最新の資産スナップショットの取得に失敗しました", err)
	}

	// スナップショットが存在しない場合は現在の資産状況を計算
	if snapshot == nil {
		s.logger.Info("スナップショットが存在しないため、現在の資産状況を計算します",
			zap.String("userID", userID.String()))
		return s.CalculateCurrentAssetSnapshot(ctx, userID, time.Now())
	}

	return dto.AssetSnapshotFromDomain(snapshot), nil
}

// CreateAssetSnapshot 資産スナップショットを作成
func (s *AssetService) CreateAssetSnapshot(
	ctx context.Context,
	userID uuid.UUID,
	req *dto.CreateAssetSnapshotRequest,
) (*dto.AssetSnapshotResponse, error) {
	// すでに同じ日付のスナップショットが存在するかチェック
	exists, err := s.assetSnapshotRepo.ExistsByUserIDAndDate(ctx, userID, req.SnapshotDate)
	if err != nil {
		s.logger.Error("スナップショット存在確認エラー",
			zap.Error(err),
			zap.String("userID", userID.String()),
			zap.Time("snapshotDate", req.SnapshotDate))
		return nil, errors.NewInternalError("スナップショットの存在確認に失敗しました", err)
	}

	if exists {
		return nil, errors.NewValidationError(
			fmt.Sprintf("指定された日付のスナップショットは既に存在します: %s", req.SnapshotDate.Format("2006-01-02")),
		)
	}

	// 口座別残高を作成
	accountBalances := make([]assetValue.AccountBalance, 0, len(req.Accounts))
	totalAssets := decimal.Zero

	for _, accountReq := range req.Accounts {
		// 口座が存在するかチェック
		account, err := s.accountRepo.FindByID(ctx, accountReq.AccountID)
		if err != nil {
			s.logger.Error("口座取得エラー",
				zap.Error(err),
				zap.String("accountID", accountReq.AccountID.String()))
			return nil, errors.NewInternalError("口座の取得に失敗しました", err)
		}

		if account == nil {
			return nil, errors.NewNotFoundError(
				fmt.Sprintf("口座が見つかりません: %s", accountReq.AccountID),
			)
		}

		// ユーザーの口座であることを確認
		if account.UserID().Value() != userID {
			return nil, errors.NewForbiddenError("指定された口座へのアクセス権限がありません")
		}

		// 金額を作成
		currency := account.CurrentBalance().Currency()
		balance, err := commonValue.NewMoney(accountReq.Balance.IntPart(), currency)
		if err != nil {
			s.logger.Error("金額作成エラー",
				zap.Error(err),
				zap.String("balance", accountReq.Balance.String()),
				zap.String("currency", currency))
			return nil, errors.NewValidationError("無効な金額です")
		}

		accountBalances = append(accountBalances, assetValue.AccountBalance{
			AccountID:   accountReq.AccountID,
			AccountName: account.Name().String(),
			Balance:     *balance,
		})

		// クレジットカード口座は負債なので資産から除外
		if !account.Type().IsCreditCard() {
			totalAssets = totalAssets.Add(accountReq.Balance)
		}
	}

	// AccountBreakdownを作成
	accountBreakdown, err := assetValue.NewAccountBreakdown(accountBalances)
	if err != nil {
		s.logger.Error("口座別内訳作成エラー", zap.Error(err))
		return nil, errors.NewValidationError("口座別内訳の作成に失敗しました")
	}

	// 総資産を作成（簡易実装：JPYのみ）
	totalAssetsMoney, err := commonValue.NewMoney(totalAssets.IntPart(), "JPY")
	if err != nil {
		s.logger.Error("総資産作成エラー",
			zap.Error(err),
			zap.String("totalAssets", totalAssets.String()))
		return nil, errors.NewInternalError("総資産の作成に失敗しました", err)
	}

	// ドメインエンティティを作成
	snapshot, err := assetDomain.NewAssetSnapshot(userID, req.SnapshotDate, *totalAssetsMoney, accountBreakdown)
	if err != nil {
		s.logger.Error("スナップショット作成エラー", zap.Error(err))
		return nil, errors.NewValidationError("スナップショットの作成に失敗しました")
	}

	// リポジトリに保存
	if err := s.assetSnapshotRepo.Save(ctx, snapshot); err != nil {
		s.logger.Error("スナップショット保存エラー", zap.Error(err))
		return nil, errors.NewInternalError("スナップショットの保存に失敗しました", err)
	}

	return dto.AssetSnapshotFromDomain(&snapshot), nil
}
