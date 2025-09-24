package handler

import (
	"net/http"
	"strconv"
	"time"

	"financetracker/internal/application/dto"
	"financetracker/internal/application/service"
	"financetracker/internal/domain/common"
	"financetracker/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TransactionHandler トランザクション関連のHTTPハンドラー
type TransactionHandler struct {
	transactionService *service.TransactionService
	logger             *logger.Logger
}

// NewTransactionHandler 新しいトランザクションハンドラーを作成
func NewTransactionHandler(transactionService *service.TransactionService, logger *logger.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		logger:             logger,
	}
}

// List トランザクション一覧を取得
// @Summary トランザクション一覧取得
// @Description ユーザーのトランザクション一覧を取得します（日付範囲とカテゴリでフィルタリング可能）
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param from query string false "開始日 (YYYY-MM-DD形式)"
// @Param to query string false "終了日 (YYYY-MM-DD形式)"
// @Param category_id query string false "カテゴリID"
// @Param limit query int false "取得件数制限" default(100)
// @Param offset query int false "取得開始位置" default(0)
// @Success 200 {object} dto.TransactionListResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) List(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
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

	// クエリパラメータを解析
	var fromDate, toDate *time.Time
	var categoryID *uuid.UUID

	// 日付フィルタの解析
	if fromStr := c.Query("from"); fromStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromStr); err == nil {
			fromDate = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "無効な開始日の形式です (YYYY-MM-DD)",
			})
			return
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if parsed, err := time.Parse("2006-01-02", toStr); err == nil {
			// 終了日は23:59:59に設定
			endOfDay := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 999999999, parsed.Location())
			toDate = &endOfDay
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "無効な終了日の形式です (YYYY-MM-DD)",
			})
			return
		}
	}

	// カテゴリフィルタの解析
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

	// ページネーション解析
	perPage := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			perPage = parsed
		}
	}

	page := 1
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			page = (parsed / perPage) + 1
		}
	}

	// 検索パラメータを作成
	searchParams := &dto.TransactionSearchParams{
		DateFrom:   fromDate,
		DateTo:     toDate,
		CategoryID: categoryID,
		Page:       page,
		PerPage:    perPage,
	}

	// サービス層を呼び出し
	transactions, err := h.transactionService.GetTransactionsByUser(
		c.Request.Context(), userUUID, searchParams,
	)
	if err != nil {
		h.logger.Error("Failed to get transaction list: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "トランザクション一覧の取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// Create トランザクションを作成
// @Summary トランザクション作成
// @Description 新しいトランザクションを作成します
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateTransactionRequest true "作成するトランザクション情報"
// @Success 201 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/transactions [post]
func (h *TransactionHandler) Create(c *gin.Context) {
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
	var req dto.CreateTransactionRequest
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
	transaction, err := h.transactionService.CreateTransaction(c.Request.Context(), userUUID, &req)
	if err != nil {
		// Check if it's a ValidationError
		if _, ok := err.(*common.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.logger.Error("Failed to create transaction: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "トランザクションの作成に失敗しました",
		})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// Get トランザクションを取得
// @Summary トランザクション取得
// @Description 指定されたIDのトランザクションを取得します
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "トランザクションID"
// @Success 200 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/transactions/{id} [get]
func (h *TransactionHandler) Get(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータからトランザクションIDを取得
	transactionIDStr := c.Param("id")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効なトランザクションIDです",
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
	transaction, err := h.transactionService.GetTransaction(c.Request.Context(), userUUID, transactionID)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "トランザクションが見つかりません",
			})
			return
		}

		h.logger.Error("Failed to get transaction: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "トランザクションの取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// Update トランザクションを更新
// @Summary トランザクション更新
// @Description 指定されたIDのトランザクションを更新します
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "トランザクションID"
// @Param request body dto.UpdateTransactionRequest true "更新するトランザクション情報"
// @Success 200 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/transactions/{id} [put]
func (h *TransactionHandler) Update(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータからトランザクションIDを取得
	transactionIDStr := c.Param("id")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効なトランザクションIDです",
		})
		return
	}

	// リクエストボディをパース
	var req dto.UpdateTransactionRequest
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
	transaction, err := h.transactionService.UpdateTransaction(c.Request.Context(), userUUID, transactionID, &req)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "トランザクションが見つかりません",
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

		h.logger.Error("Failed to update transaction: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "トランザクションの更新に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// Delete トランザクションを削除
// @Summary トランザクション削除
// @Description 指定されたIDのトランザクションを削除します
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "トランザクションID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/transactions/{id} [delete]
func (h *TransactionHandler) Delete(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータからトランザクションIDを取得
	transactionIDStr := c.Param("id")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効なトランザクションIDです",
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
	err = h.transactionService.DeleteTransaction(c.Request.Context(), userUUID, transactionID)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "トランザクションが見つかりません",
			})
			return
		}

		h.logger.Error("Failed to delete transaction: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "トランザクションの削除に失敗しました",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// MonthlySummary 月次サマリーを取得
// @Summary 月次サマリー取得
// @Description 指定された年月のトランザクションサマリーを取得します
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param year query int true "年 (YYYY)"
// @Param month query int true "月 (1-12)"
// @Success 200 {object} dto.MonthlySummaryResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/transactions/summary/monthly [get]
func (h *TransactionHandler) MonthlySummary(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// クエリパラメータから年月を取得
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "年と月を指定してください",
		})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > 2100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効な年です",
		})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効な月です",
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
	summary, err := h.transactionService.GetMonthlyTransactionSummary(c.Request.Context(), userUUID, year, month)
	if err != nil {
		h.logger.Error("Failed to get monthly summary: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "月次サマリーの取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, summary)
}
