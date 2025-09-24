package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"

	"financetracker/internal/application/dto"
	"financetracker/internal/application/service"
	"financetracker/internal/di"
	"financetracker/pkg/logger"
)

// Handler Lambda関数のハンドラー
type Handler struct {
	userService *service.UserService
	logger      *logger.Logger
}

// NewHandler 新しいHandlerを作成
func NewHandler(userService *service.UserService, logger *logger.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}

// Handle Lambdaイベントを処理
func (h *Handler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// リクエストIDをコンテキストから取得（API Gatewayで設定される）
	requestID := request.RequestContext.RequestID
	h.logger.Info("リクエスト受信",
		zap.String("method", request.HTTPMethod),
		zap.String("path", request.Path),
		zap.String("requestId", requestID))

	// Auth0ユーザーIDをAuthorizerコンテキストから取得
	auth0ID, ok := request.RequestContext.Authorizer["userId"].(string)
	if !ok {
		h.logger.Warn("認証情報が見つかりません")
		return errorResponse(401, "Unauthorized"), nil
	}

	switch request.HTTPMethod {
	case "GET":
		if request.Path == "/v1/users/me" || request.Resource == "/v1/users/me" {
			return h.handleGetCurrentUser(ctx, auth0ID)
		}
		return errorResponse(404, "Not Found"), nil
		
	case "PUT":
		if request.Path == "/v1/users/me" || request.Resource == "/v1/users/me" {
			return h.handleUpdateCurrentUser(ctx, auth0ID, request.Body)
		}
		return errorResponse(404, "Not Found"), nil
		
	default:
		return errorResponse(405, "Method Not Allowed"), nil
	}
}

// handleGetCurrentUser 現在のユーザー情報を取得
func (h *Handler) handleGetCurrentUser(ctx context.Context, auth0ID string) (events.APIGatewayProxyResponse, error) {
	user, err := h.userService.GetUserByAuth0ID(ctx, auth0ID)
	if err != nil {
		h.logger.Error("ユーザー取得エラー",
			zap.Error(err),
			zap.String("auth0Id", auth0ID))
		return errorResponse(500, "Failed to get user"), nil
	}

	// ユーザーが見つからない場合は作成
	if user == nil {
		h.logger.Info("新規ユーザー作成",
			zap.String("auth0Id", auth0ID))
		
		// Auth0から取得した情報を使用してユーザーを作成
		// 実際の実装では、Auth0のトークンから情報を取得する必要があります
		createReq := &dto.CreateUserFromAuth0Request{
			Auth0ID: auth0ID,
			Email:   fmt.Sprintf("%s@example.com", auth0ID), // 仮の実装
			Name:    "New User", // 仮の実装
		}
		
		user, err = h.userService.CreateUserFromAuth0(ctx, createReq)
		if err != nil {
			h.logger.Error("ユーザー作成エラー",
				zap.Error(err),
				zap.String("auth0Id", auth0ID))
			return errorResponse(500, "Failed to create user"), nil
		}
	}

	return successResponse(user), nil
}

// handleUpdateCurrentUser 現在のユーザー情報を更新
func (h *Handler) handleUpdateCurrentUser(ctx context.Context, auth0ID string, body string) (events.APIGatewayProxyResponse, error) {
	// リクエストボディをパース
	var updateReq dto.UpdateUserRequest
	if err := json.Unmarshal([]byte(body), &updateReq); err != nil {
		h.logger.Warn("不正なリクエストボディ",
			zap.Error(err),
			zap.String("body", body))
		return errorResponse(400, "Invalid request body"), nil
	}

	// ユーザー情報を更新
	user, err := h.userService.UpdateUserByAuth0ID(ctx, auth0ID, &updateReq)
	if err != nil {
		h.logger.Error("ユーザー更新エラー",
			zap.Error(err),
			zap.String("auth0Id", auth0ID))
		return errorResponse(500, "Failed to update user"), nil
	}

	if user == nil {
		return errorResponse(404, "User not found"), nil
	}

	return successResponse(user), nil
}

// successResponse 成功レスポンスを作成
func successResponse(data interface{}) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]interface{}{
		"success": true,
		"data":    data,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}

// errorResponse エラーレスポンスを作成
func errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]interface{}{
		"success": false,
		"error":   message,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}

func main() {
	// DIコンテナを初期化
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("DIコンテナの初期化に失敗しました: %v", err)
	}
	defer func() {
		if err := container.Close(); err != nil {
			log.Printf("DIコンテナのクローズに失敗しました: %v", err)
		}
	}()

	// ハンドラーを直接作成
	handler := NewHandler(container.UserService, container.Logger)

	// Lambda関数を開始
	lambda.Start(handler.Handle)
}