package handler

import (
	"net/http"
	"time"

	"financetracker/internal/application/dto"
	"financetracker/internal/application/service"
	"financetracker/internal/domain/common"
	"financetracker/pkg/errors"
	"financetracker/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BudgetHandler 予算関連のHTTPハンドラー
type BudgetHandler struct {
	budgetService *service.BudgetService
	userService   *service.UserService
	logger        *logger.Logger
}

// NewBudgetHandler 新しい予算ハンドラーを作成
func NewBudgetHandler(budgetService *service.BudgetService, userService *service.UserService, logger *logger.Logger) *BudgetHandler {
	return &BudgetHandler{
		budgetService: budgetService,
		userService:   userService,
		logger:        logger,
	}
}

// List 予算一覧を取得
// @Summary 予算一覧取得
// @Description ユーザーの予算一覧を取得します
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param category_id query string false "カテゴリID"
// @Param period query string false "期間" Enums(monthly, yearly)
// @Param is_active query bool false "アクティブ状態"
// @Param order_by query string false "ソート順" default("start_date desc")
// @Success 200 {object} dto.BudgetListResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/budgets [get]
func (h *BudgetHandler) List(c *gin.Context) {
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
	var categoryID *uuid.UUID
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if parsed, err := uuid.Parse(categoryIDStr); err == nil {
			categoryID = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "無効なカテゴリIDです",
			})
			return
		}
	}

	var period *string
	if periodStr := c.Query("period"); periodStr != "" {
		if periodStr == "monthly" || periodStr == "yearly" {
			period = &periodStr
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "期間は'monthly'または'yearly'を指定してください",
			})
			return
		}
	}

	var isActive *bool
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActiveStr == "true" {
			active := true
			isActive = &active
		} else if isActiveStr == "false" {
			active := false
			isActive = &active
		}
	}

	orderBy := c.Query("order_by")
	if orderBy == "" {
		orderBy = "start_date desc"
	}

	// 検索パラメータを作成
	searchParams := &dto.BudgetSearchParams{
		UserID:     userUUID,
		CategoryID: categoryID,
		Period:     period,
		IsActive:   isActive,
		OrderBy:    orderBy,
	}

	// サービス層を呼び出し
	budgets, err := h.budgetService.GetBudgetsByUser(c.Request.Context(), userUUID, searchParams)
	if err != nil {
		h.logger.Error("Failed to get budget list: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "予算一覧の取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, budgets)
}

// Create 予算を作成
// @Summary 予算作成
// @Description 新しい予算を作成します
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateBudgetRequest true "作成する予算情報"
// @Success 201 {object} dto.BudgetResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/budgets [post]
func (h *BudgetHandler) Create(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// リクエストボディをパース
	var req dto.CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("リクエストボディのパースエラー: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエストが無効です",
		})
		return
	}

	// UserIDをUUIDに変換
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.Error("Invalid user ID format: " + userID.(string))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "内部エラーが発生しました",
		})
		return
	}

	// サービス層を呼び出し
	budget, err := h.budgetService.CreateBudget(c.Request.Context(), userUUID, &req)
	if err != nil {
		// Check if it's a ValidationError
		if _, ok := err.(*common.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if it's a ConflictError
		if _, ok := err.(*common.ConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{
				"error": "指定されたカテゴリと期間の予算は既に存在します",
			})
			return
		}

		h.logger.Error("Failed to create budget: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "予算の作成に失敗しました",
		})
		return
	}

	c.JSON(http.StatusCreated, budget)
}

// Get 予算を取得
// @Summary 予算取得
// @Description 指定されたIDの予算を取得します
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "予算ID"
// @Success 200 {object} dto.BudgetResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/budgets/{id} [get]
func (h *BudgetHandler) Get(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータから予算IDを取得
	budgetIDStr := c.Param("id")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効な予算IDです",
		})
		return
	}

	// UserIDをUUIDに変換
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.Error("Invalid user ID format: " + userID.(string))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "内部エラーが発生しました",
		})
		return
	}

	// サービス層を呼び出し
	budget, err := h.budgetService.GetBudget(c.Request.Context(), userUUID, budgetID)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "予算が見つかりません",
			})
			return
		}

		// Check if it's a ForbiddenError
		if errors.GetStatusCode(err) == http.StatusForbidden {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "この予算にアクセスする権限がありません",
			})
			return
		}

		h.logger.Error("Failed to get budget: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "予算の取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, budget)
}

// Update 予算を更新
// @Summary 予算更新
// @Description 指定されたIDの予算を更新します
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "予算ID"
// @Param request body dto.UpdateBudgetRequest true "更新する予算情報"
// @Success 200 {object} dto.BudgetResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/budgets/{id} [put]
func (h *BudgetHandler) Update(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータから予算IDを取得
	budgetIDStr := c.Param("id")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効な予算IDです",
		})
		return
	}

	// リクエストボディをパース
	var req dto.UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("リクエストボディのパースエラー: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエストが無効です",
		})
		return
	}

	// UserIDをUUIDに変換
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.Error("Invalid user ID format: " + userID.(string))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "内部エラーが発生しました",
		})
		return
	}

	// サービス層を呼び出し
	budget, err := h.budgetService.UpdateBudget(c.Request.Context(), userUUID, budgetID, &req)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "予算が見つかりません",
			})
			return
		}

		// Check if it's a ForbiddenError
		if errors.GetStatusCode(err) == http.StatusForbidden {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "この予算にアクセスする権限がありません",
			})
			return
		}

		// Check if it's a ValidationError
		if _, ok := err.(*common.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.logger.Error("Failed to update budget: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "予算の更新に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, budget)
}

// Delete 予算を削除
// @Summary 予算削除
// @Description 指定されたIDの予算を削除します
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "予算ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/budgets/{id} [delete]
func (h *BudgetHandler) Delete(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータから予算IDを取得
	budgetIDStr := c.Param("id")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効な予算IDです",
		})
		return
	}

	// UserIDをUUIDに変換
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.Error("Invalid user ID format: " + userID.(string))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "内部エラーが発生しました",
		})
		return
	}

	// サービス層を呼び出し
	err = h.budgetService.DeleteBudget(c.Request.Context(), userUUID, budgetID)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "予算が見つかりません",
			})
			return
		}

		// Check if it's a ForbiddenError
		if errors.GetStatusCode(err) == http.StatusForbidden {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "この予算にアクセスする権限がありません",
			})
			return
		}

		h.logger.Error("Failed to delete budget: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "予算の削除に失敗しました",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCurrent 現在の期間の予算を取得
// @Summary 現在の期間の予算取得
// @Description 現在の日付に対してアクティブな予算を取得します
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.BudgetListResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/budgets/current [get]
func (h *BudgetHandler) GetCurrent(c *gin.Context) {
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

	// 現在の日付を基準に期間を設定
	now := time.Now()
	// 月初から月末まで
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	// サービス層を呼び出し
	budgets, err := h.budgetService.GetActiveBudgetsByPeriod(c.Request.Context(), userUUID, startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get current budgets: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "現在の予算の取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, budgets)
}
