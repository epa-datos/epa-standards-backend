// Package ports declares the interfaces (contracts) that connect the
// business logic (internal/pkg/service) to the infrastructure
// (internal/infrastructure/repositories and .../api).
//
// Rule of thumb: services depend on interfaces defined here, never on a
// concrete repository package. That's what makes it possible to swap
// Postgres for Firestore, or a real implementation for a mockery mock in
// tests, without touching the service code.
package ports

import (
	"context"

	"github.com/epa-datos/epa-standards-backend/internal/pkg/entity"
)

// ExampleRepository is the persistence contract for the Example entity.
// Implement it under internal/infrastructure/repositories/<engine> (see the
// postgres and firestore packages for two working implementations).
type ExampleRepository interface {
	GetByID(ctx context.Context, id string) (*entity.Example, error)
	GetByPage(ctx context.Context, offset, limit int64) ([]*entity.Example, error)
	GetCount(ctx context.Context) (int, error)
	Create(ctx context.Context, example *entity.Example) error
	Update(ctx context.Context, example *entity.Example) error
	Delete(ctx context.Context, id string) error
}

// ExampleService is the business-logic contract consumed by the HTTP layer
// (internal/infrastructure/api/example). Implemented by
// internal/pkg/service/example.
type ExampleService interface {
	GetByID(ctx context.Context, id string) (*entity.Example, error)
	List(ctx context.Context, offset, limit int64) (*entity.ExamplesResponse, error)
	Create(ctx context.Context, example *entity.Example) error
	Update(ctx context.Context, id string, example *entity.Example) error
	Delete(ctx context.Context, id string) error
}
