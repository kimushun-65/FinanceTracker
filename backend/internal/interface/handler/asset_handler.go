package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"financetracker/internal/application/service"
	"financetracker/pkg/logger"
)

// AssetHandler 資産関連のHTTPハンドラー
type AssetHandler struct {
	assetService *service.AssetService
	userService  *service.UserService
	logger       *logger.Logger
}

// NewAssetHandler 新しい資産ハンドラーを作成
func NewAssetHandler(
	assetService *service.AssetService,
	userService *service.UserService,
	logger *logger.Logger,
) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
		userService:  userService,
		logger:       logger,
	}
}

// GetAssetSnapshots 資産スナップショット一覧を取得
// @Summary 資産スナップショット一覧取得
// @Description 期間を指定して資産スナップショットを取得します
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param from query string false "開始日 (YYYY-MM-DD)"
// @Param to query string false "終了日 (YYYY-MM-DD)"
// @Success 200 {object} dto.AssetSnapshotListResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/assets/snapshots [get]
func (h *AssetHandler) GetAssetSnapshots(c *gin.Context) {
	// ユーザーIDを取得
	userUUID, err := getUserID(c, h.userService)
	if err != nil {
		if err == ErrUnauthorized {
			h.logger.Error("User ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "認証情報が見つかりません",
			})
		} else {
			h.logger.Error("Failed to get user ID: " + err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "ユーザーが見つかりません",
			})
		}
		return
	}

	// クエリパラメータを解析
	fromStr := c.Query("from")
	toStr := c.Query("to")

	// デフォルトの日付範囲（過去1年間）
	now := time.Now()
	startDate := now.AddDate(-1, 0, 0) // 1年前
	endDate := now

	// fromパラメータが指定されている場合
	if fromStr != "" {
		parsedFrom, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			h.logger.Error("Invalid from date: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "無効な開始日です（YYYY-MM-DD形式で指定してください）",
			})
			return
		}
		startDate = parsedFrom
	}

	// toパラメータが指定されている場合
	if toStr != "" {
		parsedTo, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			h.logger.Error("Invalid to date: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "無効な終了日です（YYYY-MM-DD形式で指定してください）",
			})
			return
		}
		endDate = parsedTo
	}

	// 日付範囲の妥当性チェック
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "開始日は終了日より前である必要があります",
		})
		return
	}

	// サービス層を呼び出し
	snapshots, err := h.assetService.GetAssetSnapshots(c.Request.Context(), userUUID, startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get asset snapshots: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "資産スナップショットの取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, snapshots)
}

// GetLatestAssetSnapshot 最新の資産スナップショットを取得
// @Summary 最新の資産スナップショット取得
// @Description 最新の資産スナップショットを取得します（存在しない場合は現在の資産状況を計算）
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.AssetSnapshotResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/assets/snapshots/latest [get]
func (h *AssetHandler) GetLatestAssetSnapshot(c *gin.Context) {
	// ユーザーIDを取得
	userUUID, err := getUserID(c, h.userService)
	if err != nil {
		if err == ErrUnauthorized {
			h.logger.Error("User ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "認証情報が見つかりません",
			})
		} else {
			h.logger.Error("Failed to get user ID: " + err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "ユーザーが見つかりません",
			})
		}
		return
	}

	// サービス層を呼び出し
	snapshot, err := h.assetService.GetLatestAssetSnapshot(c.Request.Context(), userUUID)
	if err != nil {
		h.logger.Error("Failed to get latest asset snapshot: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "最新の資産スナップショットの取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, snapshot)
}

// CalculateCurrentAssetSnapshot 現在の資産スナップショットを計算
// @Summary 現在の資産スナップショット計算
// @Description トランザクション履歴から現在の資産状況を計算します
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.AssetSnapshotResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/assets/snapshots/current [get]
func (h *AssetHandler) CalculateCurrentAssetSnapshot(c *gin.Context) {
	// ユーザーIDを取得
	userUUID, err := getUserID(c, h.userService)
	if err != nil {
		if err == ErrUnauthorized {
			h.logger.Error("User ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "認証情報が見つかりません",
			})
		} else {
			h.logger.Error("Failed to get user ID: " + err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "ユーザーが見つかりません",
			})
		}
		return
	}

	// サービス層を呼び出し
	snapshot, err := h.assetService.CalculateCurrentAssetSnapshot(c.Request.Context(), userUUID, time.Now())
	if err != nil {
		h.logger.Error("Failed to calculate current asset snapshot: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "資産スナップショットの計算に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, snapshot)
}

// CreateAssetSnapshot 資産スナップショットを作成
// @Summary 資産スナップショット作成
// @Description 新しい資産スナップショットを作成します
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateAssetSnapshotRequest true "作成する資産スナップショット情報"
// @Success 201 {object} dto.AssetSnapshotResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/assets/snapshots [post]
func (h *AssetHandler) CreateAssetSnapshot(c *gin.Context) {
	// ユーザーIDを取得
	_, err := getUserID(c, h.userService)
	if err != nil {
		h.logger.Error("Failed to get user ID: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// リクエストボディをパース
	var req struct {
		SnapshotDate string `json:"snapshot_date" binding:"required"`
		Accounts     []struct {
			AccountID string  `json:"account_id" binding:"required"`
			Balance   float64 `json:"balance" binding:"required"`
		} `json:"accounts" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("リクエストボディのパースエラー: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエストが無効です",
		})
		return
	}

	// サービス層を呼び出し
	// Note: このAPIは現時点では実装を保留（必要に応じて後で実装）
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "この機能は現在実装中です",
	})
}

// HandleNotImplemented 未実装のハンドラー
func (h *AssetHandler) HandleNotImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "この機能は現在実装中です",
	})
}
