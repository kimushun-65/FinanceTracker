package scheduler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"financetracker/internal/application/service"
	userRepo "financetracker/internal/domain/user/repository"
	"financetracker/pkg/logger"
)

// AssetSnapshotScheduler 資産スナップショット定期作成スケジューラー
type AssetSnapshotScheduler struct {
	cron         *cron.Cron
	assetService *service.AssetService
	userRepo     userRepo.UserRepository
	logger       *logger.Logger
}

// NewAssetSnapshotScheduler 新しいスケジューラーを作成
func NewAssetSnapshotScheduler(
	assetService *service.AssetService,
	userRepo userRepo.UserRepository,
	logger *logger.Logger,
) *AssetSnapshotScheduler {
	// 日本時間（JST）のタイムゾーンを使用
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		logger.Warn("Failed to load Asia/Tokyo timezone, using UTC",
			zap.Error(err))
		location = time.UTC
	}

	c := cron.New(cron.WithLocation(location))

	return &AssetSnapshotScheduler{
		cron:         c,
		assetService: assetService,
		userRepo:     userRepo,
		logger:       logger,
	}
}

// Start スケジューラーを開始
func (s *AssetSnapshotScheduler) Start() error {
	// 毎日深夜0時に実行
	_, err := s.cron.AddFunc("0 0 * * *", s.createDailySnapshots)
	if err != nil {
		s.logger.Error("Failed to add cron job for daily asset snapshots",
			zap.Error(err))
		return err
	}

	s.cron.Start()
	s.logger.Info("Asset snapshot scheduler started (runs daily at 00:00 JST)")

	return nil
}

// Stop スケジューラーを停止
func (s *AssetSnapshotScheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.logger.Info("Asset snapshot scheduler stopped")
}

// createDailySnapshots 全ユーザーの日次スナップショットを作成
func (s *AssetSnapshotScheduler) createDailySnapshots() {
	s.logger.Info("Starting daily asset snapshot creation")

	ctx := context.Background()

	// 全ユーザーを取得
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		s.logger.Error("Failed to get users for snapshot creation",
			zap.Error(err))
		return
	}

	successCount := 0
	errorCount := 0

	// 各ユーザーのスナップショットを作成
	for _, user := range users {
		err := s.assetService.CreateDailySnapshot(ctx, user.ID)
		if err != nil {
			s.logger.Error("Failed to create snapshot for user",
				zap.String("user_id", user.ID.String()),
				zap.Error(err))
			errorCount++
		} else {
			successCount++
		}
	}

	s.logger.Info("Daily asset snapshot creation completed",
		zap.Int("success", successCount),
		zap.Int("errors", errorCount),
		zap.Int("total", len(users)))
}
