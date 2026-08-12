// Package example contains the business logic for the Example entity. It
// implements ports.ExampleService and depends only on ports.ExampleRepository
// (never on a concrete database package), which is what allows tests to
// inject a mockery-generated mock instead of a real database.
package example

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/epa-datos/epa-standards-backend/internal/pkg/entity"
	"github.com/epa-datos/epa-standards-backend/internal/pkg/ports"
)

// Sentinel errors. Handlers translate these into HTTP status codes
// (see internal/infrastructure/api/example/handlers.go).
var (
	ErrNameRequired = errors.New("name is required")
	ErrNotFound     = errors.New("example not found")
)

type service struct {
	repo ports.ExampleRepository
}

// NewService builds an ports.ExampleService backed by the given repository.
// Swap the postgres/firestore implementation here, in main.go, without
// touching this file.
func NewService(repo ports.ExampleRepository) ports.ExampleService {
	return &service{repo: repo}
}

// GetByID returns a single example by id.
func (s *service) GetByID(ctx context.Context, id string) (*entity.Example, error) {
	return s.repo.GetByID(ctx, id)
}

// List returns a page of examples together with the total count.
func (s *service) List(ctx context.Context, offset, limit int64) (*entity.ExamplesResponse, error) {
	total, err := s.repo.GetCount(ctx)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetByPage(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return &entity.ExamplesResponse{TotalItems: total, Items: items}, nil
}

// Create validates and persists a new example. This is where you'd add
// business rules (uniqueness checks, defaulting, side effects, ...).
func (s *service) Create(ctx context.Context, example *entity.Example) error {
	if strings.TrimSpace(example.Name) == "" {
		return ErrNameRequired
	}

	now := time.Now().UTC()
	example.CreatedAt = now
	example.UpdatedAt = now

	return s.repo.Create(ctx, example)
}

// Update validates and persists changes to an existing example.
func (s *service) Update(ctx context.Context, id string, example *entity.Example) error {
	if strings.TrimSpace(example.Name) == "" {
		return ErrNameRequired
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}

	example.ID = id
	example.CreatedAt = existing.CreatedAt
	example.UpdatedAt = time.Now().UTC()

	return s.repo.Update(ctx, example)
}

// Delete removes an example by id.
func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
