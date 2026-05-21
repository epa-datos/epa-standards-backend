package usecases

import (
	"context"
	"fmt"
	"time"

	"epa-api/internal/domain"
)

// CreateUserUsecase creates a new user
type CreateUserUsecase struct {
	userRepo domain.UserRepository
}

// NewCreateUserUsecase creates a new CreateUserUsecase
func NewCreateUserUsecase(userRepo domain.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{
		userRepo: userRepo,
	}
}

// Execute creates a new user with the given input
func (uc *CreateUserUsecase) Execute(ctx context.Context, input domain.CreateUserInput) (*domain.User, error) {
	// 1. Validate input
	if input.Email == "" {
		return nil, domain.ErrInvalidEmail
	}
	if len(input.Name) < 2 {
		return nil, domain.ErrInvalidName
	}

	// 2. Check if user already exists
	existing, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err == nil && existing != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// 3. Create user entity
	now := time.Now()
	user := &domain.User{
		ID:        fmt.Sprintf("user_%d", now.UnixNano()), // Simple ID generation
		Email:     input.Email,
		Name:      input.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 4. Validate user
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// 5. Save to repository
	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}
