package handler

import (
	"net/http"

	"financetracker/internal/application/dto"
	"financetracker/internal/application/service"
	"financetracker/internal/domain/common"
	"financetracker/pkg/errors"
	"financetracker/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CategoryHandler カテゴリ関連のHTTPハンドラー
type CategoryHandler struct {
	categoryService *service.CategoryService
	userService     *service.UserService
	logger          *logger.Logger
}

// NewCategoryHandler 新しいカテゴリハンドラーを作成
func NewCategoryHandler(categoryService *service.CategoryService, userService *service.UserService, logger *logger.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		userService:     userService,
		logger:          logger,
	}
}

// List ユーザーカテゴリ一覧を取得
// @Summary ユーザーカテゴリ一覧取得
// @Description ユーザーのカスタムカテゴリ一覧を取得します
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param category_master_id query string false "カテゴリマスターID"
// @Param is_active query bool false "アクティブ状態"
// @Param order_by query string false "ソート順" default("created_at desc")
// @Success 200 {object} dto.CategoryListResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
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
	var categoryMasterID *uuid.UUID
	if categoryMasterIDStr := c.Query("category_master_id"); categoryMasterIDStr != "" {
		if parsed, err := uuid.Parse(categoryMasterIDStr); err == nil {
			categoryMasterID = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "無効なカテゴリマスターIDです",
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
		orderBy = "created_at desc"
	}

	// 検索パラメータを作成
	searchParams := &dto.CategorySearchParams{
		UserID:           userUUID,
		CategoryMasterID: categoryMasterID,
		IsActive:         isActive,
		OrderBy:          orderBy,
	}

	// サービス層を呼び出し
	categories, err := h.categoryService.GetCategoriesByUser(c.Request.Context(), userUUID, searchParams)
	if err != nil {
		h.logger.Error("Failed to get category list: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "カテゴリ一覧の取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Create カテゴリを作成
// @Summary カテゴリ作成
// @Description 新しいユーザーカテゴリを作成します
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateCategoryRequest true "作成するカテゴリ情報"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
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
	var req dto.CreateCategoryRequest
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
	category, err := h.categoryService.CreateCategory(c.Request.Context(), userUUID, &req)
	if err != nil {
		// Check if it's a ValidationError
		if _, ok := err.(*common.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.logger.Error("Failed to create category: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "カテゴリの作成に失敗しました",
		})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// Get カテゴリを取得
// @Summary カテゴリ取得
// @Description 指定されたIDのカテゴリを取得します
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "カテゴリID"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) Get(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータからカテゴリIDを取得
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効なカテゴリIDです",
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
	category, err := h.categoryService.GetCategory(c.Request.Context(), userUUID, categoryID)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "カテゴリが見つかりません",
			})
			return
		}

		// Check if it's a ForbiddenError
		if errors.GetStatusCode(err) == http.StatusForbidden {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "このカテゴリにアクセスする権限がありません",
			})
			return
		}

		h.logger.Error("Failed to get category: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "カテゴリの取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Update カテゴリを更新
// @Summary カテゴリ更新
// @Description 指定されたIDのカテゴリを更新します
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "カテゴリID"
// @Param request body dto.UpdateCategoryRequest true "更新するカテゴリ情報"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータからカテゴリIDを取得
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効なカテゴリIDです",
		})
		return
	}

	// リクエストボディをパース
	var req dto.UpdateCategoryRequest
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
	category, err := h.categoryService.UpdateCategory(c.Request.Context(), userUUID, categoryID, &req)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "カテゴリが見つかりません",
			})
			return
		}

		// Check if it's a ForbiddenError
		if errors.GetStatusCode(err) == http.StatusForbidden {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "このカテゴリにアクセスする権限がありません",
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

		h.logger.Error("Failed to update category: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "カテゴリの更新に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Delete カテゴリを削除（非活性化）
// @Summary カテゴリ削除
// @Description 指定されたIDのカテゴリを削除（非活性化）します
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "カテゴリID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// パスパラメータからカテゴリIDを取得
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効なカテゴリIDです",
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
	err = h.categoryService.DeleteCategory(c.Request.Context(), userUUID, categoryID)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "カテゴリが見つかりません",
			})
			return
		}

		// Check if it's a ForbiddenError
		if errors.GetStatusCode(err) == http.StatusForbidden {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "このカテゴリにアクセスする権限がありません",
			})
			return
		}

		h.logger.Error("Failed to delete category: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "カテゴリの削除に失敗しました",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListMaster カテゴリマスター一覧を取得
// @Summary カテゴリマスター一覧取得
// @Description 利用可能なカテゴリマスター（テンプレート）一覧を取得します
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param category_type query string false "カテゴリタイプ" Enums(income, expense)
// @Param order_by query string false "ソート順" default("display_order asc")
// @Success 200 {object} dto.CategoryMasterListResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/categories/master [get]
func (h *CategoryHandler) ListMaster(c *gin.Context) {
	// 認証確認（ログインユーザーのみアクセス可能）
	_, exists := c.Get("UserID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// クエリパラメータを解析
	var categoryType *string
	if categoryTypeStr := c.Query("category_type"); categoryTypeStr != "" {
		if categoryTypeStr == "income" || categoryTypeStr == "expense" {
			categoryType = &categoryTypeStr
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "カテゴリタイプは'income'または'expense'を指定してください",
			})
			return
		}
	}

	orderBy := c.Query("order_by")
	if orderBy == "" {
		orderBy = "display_order asc"
	}

	// 検索パラメータを作成
	searchParams := &dto.CategoryMasterSearchParams{
		CategoryType: categoryType,
		OrderBy:      orderBy,
	}

	// サービス層を呼び出し
	categoryMasters, err := h.categoryService.GetCategoryMasters(c.Request.Context(), searchParams)
	if err != nil {
		h.logger.Error("Failed to get category master list: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "カテゴリマスター一覧の取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, categoryMasters)
}
