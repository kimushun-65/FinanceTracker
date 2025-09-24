package handler

import (
	"net/http"

	"financetracker/internal/application/dto"
	"financetracker/internal/application/service"
	"financetracker/internal/domain/common"
	"financetracker/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AccountHandler 口座関連のHTTPハンドラー
type AccountHandler struct {
	accountService *service.AccountService
	logger         *logger.Logger
}

// NewAccountHandler 新しい口座ハンドラーを作成
func NewAccountHandler(accountService *service.AccountService, logger *logger.Logger) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		logger:         logger,
	}
}

// List 口座一覧を取得
// @Summary 口座一覧取得
// @Description ユーザーの口座一覧を取得します
// @Tags accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param limit query int false "取得件数制限" default(100)
// @Param offset query int false "取得開始位置" default(0)
// @Success 200 {object} dto.AccountListResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounts [get]
func (h *AccountHandler) List(c *gin.Context) {
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

	// サービス層を呼び出し
	accounts, err := h.accountService.GetAccountsByUser(c.Request.Context(), userUUID)
	if err != nil {
		h.logger.Error("Failed to get account list: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "口座一覧の取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// Create 口座を作成
// @Summary 口座作成
// @Description 新しい口座を作成します
// @Tags accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateAccountRequest true "作成する口座情報"
// @Success 201 {object} dto.AccountResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounts [post]
func (h *AccountHandler) Create(c *gin.Context) {
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
	var req dto.CreateAccountRequest
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
	account, err := h.accountService.CreateAccount(c.Request.Context(), userUUID, &req)
	if err != nil {
		// Check if it's a ValidationError
		if _, ok := err.(*common.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.logger.Error("Failed to create account: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "口座の作成に失敗しました",
		})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// Get 口座を取得
// @Summary 口座取得
// @Description 指定されたIDの口座を取得します
// @Tags accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "口座ID"
// @Success 200 {object} dto.AccountResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounts/{id} [get]
func (h *AccountHandler) Get(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータから口座IDを取得
	accountIDStr := c.Param("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効な口座IDです",
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
	account, err := h.accountService.GetAccount(c.Request.Context(), userUUID, accountID)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "口座が見つかりません",
			})
			return
		}

		h.logger.Error("Failed to get account: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "口座の取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, account)
}

// Update 口座を更新
// @Summary 口座更新
// @Description 指定されたIDの口座を更新します
// @Tags accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "口座ID"
// @Param request body dto.UpdateAccountRequest true "更新する口座情報"
// @Success 200 {object} dto.AccountResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounts/{id} [put]
func (h *AccountHandler) Update(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータから口座IDを取得
	accountIDStr := c.Param("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効な口座IDです",
		})
		return
	}

	// リクエストボディをパース
	var req dto.UpdateAccountRequest
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
	account, err := h.accountService.UpdateAccount(c.Request.Context(), userUUID, accountID, &req)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "口座が見つかりません",
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

		h.logger.Error("Failed to update account: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "口座の更新に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, account)
}

// Delete 口座を削除
// @Summary 口座削除
// @Description 指定されたIDの口座を削除します
// @Tags accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "口座ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounts/{id} [delete]
func (h *AccountHandler) Delete(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータから口座IDを取得
	accountIDStr := c.Param("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効な口座IDです",
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
	err = h.accountService.DeleteAccount(c.Request.Context(), userUUID, accountID)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "口座が見つかりません",
			})
			return
		}

		h.logger.Error("Failed to delete account: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "口座の削除に失敗しました",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
