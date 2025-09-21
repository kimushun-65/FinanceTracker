// Package entity 資産予測関連のエンティティを定義
package entity

import (
	"time"

	"financetracker/internal/domain/asset/value"
	"financetracker/internal/domain/common"
	commonValue "financetracker/internal/domain/common/value"

	"github.com/google/uuid"
)

// AssetForecast 資産予測エンティティ
type AssetForecast struct {
	common.BaseEntity
	userID          uuid.UUID
	horizonMonths   int
	forecastDate    time.Time
	predictedAssets commonValue.Money
	assumptions     value.Assumptions
	confidence      float64
}

// NewAssetForecast 新しい資産予測を作成
func NewAssetForecast(
	userID uuid.UUID,
	horizonMonths int,
	forecastDate time.Time,
	predictedAssets commonValue.Money,
	assumptions value.Assumptions,
	confidence float64,
) (AssetForecast, error) {
	if err := validateAssetForecastParams(
		userID, horizonMonths, forecastDate, confidence,
	); err != nil {
		return AssetForecast{}, err
	}

	baseEntity := common.NewBaseEntity()

	return AssetForecast{
		BaseEntity:      baseEntity,
		userID:          userID,
		horizonMonths:   horizonMonths,
		forecastDate:    forecastDate,
		predictedAssets: predictedAssets,
		assumptions:     assumptions,
		confidence:      confidence,
	}, nil
}

// validateAssetForecastParams 資産予測パラメータの妥当性を検証
func validateAssetForecastParams(
	userID uuid.UUID,
	horizonMonths int,
	forecastDate time.Time,
	confidence float64,
) error {
	if userID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"ユーザーIDが必要です",
		)
	}

	if horizonMonths < 1 || horizonMonths > 60 {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"予測期間は1〜60ヶ月の範囲である必要があります",
		)
	}

	if forecastDate.IsZero() {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"予測基準日が必要です",
		)
	}

	if confidence < 0 || confidence > 1 {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"信頼区間は0以上1以下である必要があります",
		)
	}

	return nil
}

// UserID ユーザーIDを取得
func (a AssetForecast) UserID() uuid.UUID {
	return a.userID
}

// HorizonMonths 予測期間（月数）を取得
func (a AssetForecast) HorizonMonths() int {
	return a.horizonMonths
}

// ForecastDate 予測基準日を取得
func (a AssetForecast) ForecastDate() time.Time {
	return a.forecastDate
}

// PredictedAssets 予測資産額を取得
func (a AssetForecast) PredictedAssets() commonValue.Money {
	return a.predictedAssets
}

// Assumptions 前提条件を取得
func (a AssetForecast) Assumptions() value.Assumptions {
	return a.assumptions
}

// Confidence 信頼度を取得
func (a AssetForecast) Confidence() float64 {
	return a.confidence
}

// GetTargetDate 予測対象日を計算
func (a AssetForecast) GetTargetDate() time.Time {
	return a.forecastDate.AddDate(0, a.horizonMonths, 0)
}

// IsHighConfidence 高信頼度の予測かどうかを判定
func (a AssetForecast) IsHighConfidence() bool {
	return a.confidence >= 0.8
}

// IsValidForCurrentDate 現在日付に対して有効な予測かどうかを判定
func (a AssetForecast) IsValidForCurrentDate() bool {
	// 予測基準日から1ヶ月以内であれば有効とする
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	return a.CreatedAt.After(oneMonthAgo)
}

// HasSufficientData 十分なデータがあるかチェック
func (a AssetForecast) HasSufficientData(dataMonths int) bool {
	// 単一の予測手法を想定し、最低6ヶ月のデータが必要
	const requiredMonths = 6
	return dataMonths >= requiredMonths
}

// CalculateConfidenceInterval 信頼区間を計算
func (a AssetForecast) CalculateConfidenceInterval() (lower, upper *commonValue.Money, err error) {
	// 簡易的な計算：予測額の±(1-confidence)*50%
	variance := (1 - a.confidence) * 0.5

	lowerAmount := int64(float64(a.predictedAssets.Amount()) * (1 - variance))
	upperAmount := int64(float64(a.predictedAssets.Amount()) * (1 + variance))

	lower, err = commonValue.NewMoney(lowerAmount, a.predictedAssets.Currency())
	if err != nil {
		return nil, nil, err
	}

	upper, err = commonValue.NewMoney(upperAmount, a.predictedAssets.Currency())
	if err != nil {
		return nil, nil, err
	}

	return lower, upper, nil
}
