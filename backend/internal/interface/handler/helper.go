package handler

import (
	"errors"

	"financetracker/internal/application/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	// ErrUnauthorized 認証エラー
	ErrUnauthorized = errors.New("unauthorized")
)

// getUserID はコンテキストからユーザーIDを取得する共通ヘルパー関数
// UserIDが設定されていない場合は、Auth0 IDから実際のユーザーを検索する
func getUserID(c *gin.Context, userService *service.UserService) (uuid.UUID, error) {
	// 認証情報からユーザーIDを取得
	userID, exists := c.Get("UserID")
	if exists {
		// UserIDをUUIDに変換
		userUUID, err := uuid.Parse(userID.(string))
		if err != nil {
			return uuid.Nil, err
		}
		return userUUID, nil
	}

	// UserIDが設定されていない場合、Auth0 IDから取得を試みる
	auth0ID, auth0Exists := c.Get("auth0_id")
	if !auth0Exists {
		return uuid.Nil, ErrUnauthorized
	}

	user, err := userService.GetUserByAuth0ID(c.Request.Context(), auth0ID.(string))
	if err != nil {
		return uuid.Nil, err
	}

	// 次回のために UserID をコンテキストに保存
	c.Set("UserID", user.ID.String())

	return user.ID, nil
}
