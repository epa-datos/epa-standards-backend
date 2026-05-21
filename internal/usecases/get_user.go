package usecases

import (
	"context"

	"epa-api/internal/domain"
)

// GetUserUsecase retrieves a user by ID
type GetUserUsecase struct {
	userRepo domain.UserRepository
}

// NewGetUserUsecase creates a new GetUserUsecase
func NewGetUserUsecase(userRepo domain.UserRepository) *GetUserUsecase {
	return &GetUserUsecase{
		userRepo: userRepo,
	}
}

// Execute retrieves a user by ID
func (uc *GetUserUsecase) Execute(ctx context.Context, id string) (*domain.User, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}
