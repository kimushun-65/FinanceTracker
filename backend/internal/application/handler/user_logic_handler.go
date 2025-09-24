package handler

import (
	"context"

	"financetracker/internal/application/dto"
	"financetracker/internal/application/service"
	"financetracker/pkg/logger"
)

// UserLogicHandler ユーザー関連のビジネスロジックハンドラー
// Lambda と HTTP の両方から利用可能
type UserLogicHandler struct {
	userService *service.UserService
	logger      *logger.Logger
}

// NewUserLogicHandler 新しいUserLogicHandlerを作成
func NewUserLogicHandler(userService *service.UserService, logger *logger.Logger) *UserLogicHandler {
	return &UserLogicHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetCurrentUser 現在のユーザー情報を取得
func (h *UserLogicHandler) GetCurrentUser(ctx context.Context, auth0ID string) (*dto.UserResponse, error) {
	// ユーザー情報を取得
	user, err := h.userService.GetUserByAuth0ID(ctx, auth0ID)
	if err != nil {
		h.logger.Error("ユーザー取得エラー: " + err.Error() + " auth0Id: " + auth0ID)
		return nil, err
	}

	// ユーザーが見つからない場合は作成
	if user == nil {
		h.logger.Info("新規ユーザー作成 auth0Id: " + auth0ID)

		// TODO: 実際の実装では、Auth0のトークンから情報を取得する必要があります
		createReq := &dto.CreateUserFromAuth0Request{
			Auth0ID: auth0ID,
			Email:   auth0ID + "@example.com", // 仮の実装
			Name:    "New User",               // 仮の実装
		}

		user, err = h.userService.CreateUserFromAuth0(ctx, createReq)
		if err != nil {
			h.logger.Error("ユーザー作成エラー: " + err.Error() + " auth0Id: " + auth0ID)
			return nil, err
		}
	}

	return user, nil
}

// UpdateCurrentUser 現在のユーザー情報を更新
func (h *UserLogicHandler) UpdateCurrentUser(ctx context.Context, auth0ID string, updateReq *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// ユーザー情報を更新
	user, err := h.userService.UpdateUserByAuth0ID(ctx, auth0ID, updateReq)
	if err != nil {
		h.logger.Error("ユーザー更新エラー: " + err.Error() + " auth0Id: " + auth0ID)
		return nil, err
	}

	if user == nil {
		h.logger.Error("ユーザーが見つかりません auth0Id: " + auth0ID)
		return nil, err
	}

	return user, nil
}
