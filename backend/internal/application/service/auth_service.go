package service

import (
	"context"
	"errors"
	"fmt"

	"financetracker/internal/domain/common"
	"financetracker/internal/domain/common/value"
	"financetracker/internal/domain/user/entity"
	"financetracker/internal/domain/user/repository"
	userValue "financetracker/internal/domain/user/value"
	"financetracker/internal/infrastructure/auth0"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) SyncUser(ctx context.Context, userInfo *auth0.UserInfo, claims *auth0.TokenClaims) (*entity.User, error) {
	// Import necessary value objects
	auth0IDValue, err := userValue.NewAuth0ID(userInfo.Sub)
	if err != nil {
		return nil, fmt.Errorf("failed to create Auth0 ID: %w", err)
	}

	emailValue, err := value.NewEmail(userInfo.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create email value: %w", err)
	}

	// Check if user already exists
	existingUser, err := s.userRepo.FindByAuth0UserID(ctx, *auth0IDValue)
	if err != nil {
		var notFoundErr *common.NotFoundError
		if !errors.As(err, &notFoundErr) {
			return nil, fmt.Errorf("failed to check existing user: %w", err)
		}
	}

	if existingUser != nil {
		// Update existing user
		if err := existingUser.UpdateProfile(userInfo.Name, *emailValue); err != nil {
			return nil, fmt.Errorf("failed to update user profile: %w", err)
		}

		if err := s.userRepo.Save(ctx, existingUser); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		return existingUser, nil
	}

	// Create new user
	newUser, err := entity.NewUser(*auth0IDValue, *emailValue, userInfo.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	if err := s.userRepo.Save(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

func (s *AuthService) GetUserByAuth0ID(ctx context.Context, auth0ID string) (*entity.User, error) {
	auth0IDValue, err := userValue.NewAuth0ID(auth0ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create Auth0 ID: %w", err)
	}

	user, err := s.userRepo.FindByAuth0UserID(ctx, *auth0IDValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by Auth0 ID: %w", err)
	}
	return user, nil
}
