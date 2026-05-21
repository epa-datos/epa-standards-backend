package usecases

import (
	"context"
	"testing"

	"epa-api/internal/domain"
	"epa-api/internal/adapters/persistence"
	"epa-api/pkg/logger"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		input     domain.CreateUserInput
		wantError bool
		errType   error
	}{
		{
			name: "valid user creation",
			input: domain.CreateUserInput{
				Email: "user@example.com",
				Name:  "John Doe",
			},
			wantError: false,
		},
		{
			name: "invalid email",
			input: domain.CreateUserInput{
				Email: "",
				Name:  "John",
			},
			wantError: true,
			errType:   domain.ErrInvalidEmail,
		},
		{
			name: "invalid name",
			input: domain.CreateUserInput{
				Email: "user@example.com",
				Name:  "A",
			},
			wantError: true,
			errType:   domain.ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			log := logger.New("debug")
			repo := persistence.NewUserRepository(log)
			usecase := NewCreateUserUsecase(repo)

			// Execute
			user, err := usecase.Execute(context.Background(), tt.input)

			// Assert
			if (err != nil) != tt.wantError {
				t.Errorf("got error %v, want error %v", err, tt.wantError)
			}

			if tt.wantError && err != tt.errType {
				t.Errorf("got error %v, want error %v", err, tt.errType)
			}

			if !tt.wantError {
				if user.Email != tt.input.Email {
					t.Errorf("got email %s, want %s", user.Email, tt.input.Email)
				}
				if user.Name != tt.input.Name {
					t.Errorf("got name %s, want %s", user.Name, tt.input.Name)
				}
			}
		})
	}
}

func TestDuplicateUser(t *testing.T) {
	// Setup
	log := logger.New("debug")
	repo := persistence.NewUserRepository(log)
	usecase := NewCreateUserUsecase(repo)

	input := domain.CreateUserInput{
		Email: "user@example.com",
		Name:  "John",
	}

	// Create first user
	_, err := usecase.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to create first user: %v", err)
	}

	// Try to create duplicate
	_, err = usecase.Execute(context.Background(), input)
	if err != domain.ErrUserAlreadyExists {
		t.Errorf("got error %v, want ErrUserAlreadyExists", err)
	}
}
