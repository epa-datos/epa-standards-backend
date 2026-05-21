package domain

import (
	"context"
	"errors"
	"time"
)

// User domain entity
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserInput represents the input for creating a user
type CreateUserInput struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=2"`
}

// UserRepository defines the contract for user persistence
type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Delete(ctx context.Context, id string) error
}

// DomainErrors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidName       = errors.New("name is required and must be at least 2 characters")
)

// Validate validates a user
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if len(u.Name) < 2 {
		return ErrInvalidName
	}
	return nil
}
