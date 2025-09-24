package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"financetracker/internal/application/dto"
	"financetracker/internal/domain/common/value"
	userDomain "financetracker/internal/domain/user/entity"
	userRepo "financetracker/internal/domain/user/repository"
	userValue "financetracker/internal/domain/user/value"
	"financetracker/pkg/errors"
	"financetracker/pkg/logger"
)

// UserService ユーザー管理サービス
type UserService struct {
	userRepo userRepo.UserRepository
	logger   *logger.Logger
}

// NewUserService 新しいUserServiceを作成
func NewUserService(
	userRepo userRepo.UserRepository,
	logger *logger.Logger,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetUser ユーザー情報を取得
func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	// ユーザーIDからドメインのUserID値オブジェクトを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリからユーザーを取得
	user, err := s.userRepo.FindByID(ctx, domainUserID.Value())
	if err != nil {
		s.logger.Error("ユーザー取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()))
		return nil, errors.NewInternalError("ユーザー情報の取得に失敗しました", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("ユーザーが見つかりません: %s", userID))
	}

	// DTOに変換
	return dto.UserFromDomain(user), nil
}

// UpdateUser ユーザー情報を更新
func (s *UserService) UpdateUser(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// ユーザーIDからドメインのUserID値オブジェクトを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリからユーザーを取得
	user, err := s.userRepo.FindByID(ctx, domainUserID.Value())
	if err != nil {
		s.logger.Error("ユーザー取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()))
		return nil, errors.NewInternalError("ユーザー情報の取得に失敗しました", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("ユーザーが見つかりません: %s", userID))
	}

	// 更新フィールドの適用
	if req.Name != nil {
		name := *req.Name
		if req.Email != nil {
			// メールアドレスも更新する場合
			email, err := value.NewEmail(*req.Email)
			if err != nil {
				s.logger.Error("メールアドレス作成エラー",
					zap.Error(err),
					zap.String("email", *req.Email))
				return nil, errors.NewValidationError("無効なメールアドレスです")
			}
			if err := user.UpdateProfile(name, *email); err != nil {
				s.logger.Error("プロファイル更新エラー", zap.Error(err))
				return nil, errors.NewValidationError("プロファイルの更新に失敗しました")
			}
		} else {
			// 名前のみ更新
			if err := user.UpdateProfile(name, user.Email()); err != nil {
				s.logger.Error("プロファイル更新エラー", zap.Error(err))
				return nil, errors.NewValidationError("プロファイルの更新に失敗しました")
			}
		}
	} else if req.Email != nil {
		// メールアドレスのみ更新
		email, err := value.NewEmail(*req.Email)
		if err != nil {
			s.logger.Error("メールアドレス作成エラー",
				zap.Error(err),
				zap.String("email", *req.Email))
			return nil, errors.NewValidationError("無効なメールアドレスです")
		}
		if err := user.UpdateProfile(user.Name(), *email); err != nil {
			s.logger.Error("プロファイル更新エラー", zap.Error(err))
			return nil, errors.NewValidationError("プロファイルの更新に失敗しました")
		}
	}

	// リポジトリで更新（SaveはUpdateも兼ねる）
	if err := s.userRepo.Save(ctx, user); err != nil {
		s.logger.Error("ユーザー更新エラー",
			zap.Error(err),
			zap.String("userID", userID.String()))
		return nil, errors.NewInternalError("ユーザー情報の更新に失敗しました", err)
	}

	// DTOに変換して返却
	return dto.UserFromDomain(user), nil
}

// GetUserByAuth0ID Auth0IDでユーザーを取得
func (s *UserService) GetUserByAuth0ID(ctx context.Context, auth0ID string) (*dto.UserResponse, error) {
	// Auth0IDの値オブジェクトを作成
	auth0IDObj, err := userValue.NewAuth0ID(auth0ID)
	if err != nil {
		s.logger.Error("Auth0IDの作成に失敗しました",
			zap.Error(err),
			zap.String("auth0Id", auth0ID))
		return nil, errors.NewValidationError("無効なAuth0IDです")
	}

	// ユーザーを取得
	user, err := s.userRepo.FindByAuth0UserID(ctx, *auth0IDObj)
	if err != nil {
		s.logger.Error("ユーザー取得エラー",
			zap.Error(err),
			zap.String("auth0Id", auth0ID))
		return nil, errors.NewInternalError("ユーザーの取得に失敗しました", err)
	}

	if user == nil {
		return nil, nil
	}

	// DTOに変換して返却
	return dto.UserFromDomain(user), nil
}

// UpdateUserByAuth0ID Auth0IDでユーザー情報を更新
func (s *UserService) UpdateUserByAuth0ID(ctx context.Context, auth0ID string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Auth0IDの値オブジェクトを作成
	auth0IDObj, err := userValue.NewAuth0ID(auth0ID)
	if err != nil {
		s.logger.Error("Auth0IDの作成に失敗しました",
			zap.Error(err),
			zap.String("auth0Id", auth0ID))
		return nil, errors.NewValidationError("無効なAuth0IDです")
	}

	// ユーザーを取得
	user, err := s.userRepo.FindByAuth0UserID(ctx, *auth0IDObj)
	if err != nil {
		s.logger.Error("ユーザー取得エラー",
			zap.Error(err),
			zap.String("auth0Id", auth0ID))
		return nil, errors.NewInternalError("ユーザーの取得に失敗しました", err)
	}

	if user == nil {
		return nil, nil
	}

	// 更新フィールドの適用
	if req.Name != nil {
		name := *req.Name
		if req.Email != nil {
			// メールアドレスも更新する場合
			email, err := value.NewEmail(*req.Email)
			if err != nil {
				s.logger.Error("メールアドレス作成エラー",
					zap.Error(err),
					zap.String("email", *req.Email))
				return nil, errors.NewValidationError("無効なメールアドレスです")
			}
			if err := user.UpdateProfile(name, *email); err != nil {
				s.logger.Error("プロファイル更新エラー", zap.Error(err))
				return nil, errors.NewValidationError("プロファイルの更新に失敗しました")
			}
		} else {
			// 名前のみ更新
			if err := user.UpdateProfile(name, user.Email()); err != nil {
				s.logger.Error("プロファイル更新エラー", zap.Error(err))
				return nil, errors.NewValidationError("プロファイルの更新に失敗しました")
			}
		}
	} else if req.Email != nil {
		// メールアドレスのみ更新
		email, err := value.NewEmail(*req.Email)
		if err != nil {
			s.logger.Error("メールアドレス作成エラー",
				zap.Error(err),
				zap.String("email", *req.Email))
			return nil, errors.NewValidationError("無効なメールアドレスです")
		}
		if err := user.UpdateProfile(user.Name(), *email); err != nil {
			s.logger.Error("プロファイル更新エラー", zap.Error(err))
			return nil, errors.NewValidationError("プロファイルの更新に失敗しました")
		}
	}

	// リポジトリに保存
	if err := s.userRepo.Save(ctx, user); err != nil {
		s.logger.Error("ユーザー保存エラー",
			zap.Error(err),
			zap.String("userId", user.ID.String()))
		return nil, errors.NewInternalError("ユーザーの更新に失敗しました", err)
	}

	// DTOに変換して返却
	return dto.UserFromDomain(user), nil
}

// CreateUserFromAuth0 Auth0からの情報でユーザーを作成（内部利用）
func (s *UserService) CreateUserFromAuth0(ctx context.Context, req *dto.CreateUserFromAuth0Request) (*dto.UserResponse, error) {
	// Auth0IDを作成
	auth0ID, err := userValue.NewAuth0ID(req.Auth0ID)
	if err != nil {
		s.logger.Error("Auth0ID作成エラー",
			zap.Error(err),
			zap.String("auth0ID", req.Auth0ID))
		return nil, errors.NewValidationError("無効なAuth0IDです")
	}
	
	// メールアドレスを作成
	email, err := value.NewEmail(req.Email)
	if err != nil {
		s.logger.Error("メールアドレス作成エラー",
			zap.Error(err),
			zap.String("email", req.Email))
		return nil, errors.NewValidationError("無効なメールアドレスです")
	}

	// 既存ユーザーの確認
	existingUser, err := s.userRepo.FindByAuth0UserID(ctx, *auth0ID)
	if err != nil {
		s.logger.Error("既存ユーザー確認エラー",
			zap.Error(err),
			zap.String("auth0ID", req.Auth0ID))
		return nil, errors.NewInternalError("ユーザーの確認に失敗しました", err)
	}

	if existingUser != nil {
		// 既存ユーザーが存在する場合は、そのユーザーを返す
		return dto.UserFromDomain(existingUser), nil
	}

	// 新規ユーザーを作成
	user, err := userDomain.NewUser(*auth0ID, *email, req.Name)
	if err != nil {
		s.logger.Error("ユーザー作成エラー", zap.Error(err))
		return nil, errors.NewValidationError("ユーザーの作成に失敗しました")
	}

	// リポジトリに保存
	if err := s.userRepo.Save(ctx, user); err != nil {
		s.logger.Error("ユーザー保存エラー", zap.Error(err))
		return nil, errors.NewInternalError("ユーザーの保存に失敗しました", err)
	}

	// DTOに変換して返却
	return dto.UserFromDomain(user), nil
}
