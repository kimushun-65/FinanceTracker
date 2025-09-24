package handler

import (
	"net/http"

	"financetracker/internal/application/dto"
	"financetracker/internal/application/handler"
	"financetracker/internal/application/service"
	"financetracker/internal/domain/common"
	"financetracker/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UserHandler ユーザー関連のHTTPハンドラー
type UserHandler struct {
	userService      *service.UserService
	userLogicHandler *handler.UserLogicHandler
	logger           *logger.Logger
}

// NewUserHandler 新しいユーザーハンドラーを作成
func NewUserHandler(userService *service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService:      userService,
		userLogicHandler: handler.NewUserLogicHandler(userService, logger),
		logger:           logger,
	}
}

// GetCurrentUser 現在のユーザー情報を取得
// @Summary 現在のユーザー情報取得
// @Description 認証済みユーザーの情報を取得します
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// 認証情報からAuth0 IDを取得
	auth0ID, exists := c.Get("auth0_id")
	if !exists {
		h.logger.Error("Auth0 ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// ビジネスロジックハンドラーを使用してユーザー情報を取得
	user, err := h.userLogicHandler.GetCurrentUser(c.Request.Context(), auth0ID.(string))
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "ユーザーが見つかりません",
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

		// Any other error is internal server error
		h.logger.Error("Failed to get user: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザー情報の取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateCurrentUser 現在のユーザー情報を更新
// @Summary 現在のユーザー情報更新
// @Description 認証済みユーザーの情報を更新します
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.UpdateUserRequest true "更新するユーザー情報"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/me [put]
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	// 認証情報からAuth0 IDを取得
	auth0ID, exists := c.Get("auth0_id")
	if !exists {
		h.logger.Error("Auth0 ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証情報が見つかりません",
		})
		return
	}

	// リクエストボディをパース
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("リクエストボディのパースエラー: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエストが無効です",
		})
		return
	}

	// ビジネスロジックハンドラーを使用してユーザー情報を更新
	user, err := h.userLogicHandler.UpdateCurrentUser(c.Request.Context(), auth0ID.(string), &req)
	if err != nil {
		// Check if it's a NotFoundError
		if _, ok := err.(*common.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "ユーザーが見つかりません",
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

		// Any other error is internal server error
		h.logger.Error("Failed to update user: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザー情報の更新に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
