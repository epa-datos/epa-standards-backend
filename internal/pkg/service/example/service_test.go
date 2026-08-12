package example

import (
	"context"
	"testing"
	"time"

	"github.com/epa-datos/epa-standards-backend/internal/pkg/entity"
	"github.com/epa-datos/epa-standards-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// This file is the reference example for unit-testing a service: mock the
// port (ports.ExampleRepository) with the mockery-generated mock, program
// its expectations with .On(...), and assert on the service's behavior.
// See docs/TESTING.md for the full write-up.

func TestService_GetByID(t *testing.T) {
	repo := mocks.NewExampleRepository(t)
	svc := NewService(repo)

	want := &entity.Example{ID: "1", Name: "Sample"}
	repo.On("GetByID", mock.Anything, "1").Return(want, nil)

	got, err := svc.GetByID(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		example   *entity.Example
		mockSetup func(repo *mocks.ExampleRepository)
		wantErr   error
	}{
		{
			name:    "valid example is created",
			example: &entity.Example{Name: "Sample"},
			mockSetup: func(repo *mocks.ExampleRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Example")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "empty name is rejected before touching the repository",
			example:   &entity.Example{Name: "   "},
			mockSetup: func(repo *mocks.ExampleRepository) {},
			wantErr:   ErrNameRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewExampleRepository(t)
			tt.mockSetup(repo)
			svc := NewService(repo)

			err := svc.Create(context.Background(), tt.example)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.False(t, tt.example.CreatedAt.IsZero(), "CreatedAt should be set by the service")
		})
	}
}

func TestService_Update_NotFound(t *testing.T) {
	repo := mocks.NewExampleRepository(t)
	svc := NewService(repo)

	repo.On("GetByID", mock.Anything, "missing").Return(nil, nil)

	err := svc.Update(context.Background(), "missing", &entity.Example{Name: "Sample"})

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestService_List(t *testing.T) {
	repo := mocks.NewExampleRepository(t)
	svc := NewService(repo)

	items := []*entity.Example{{ID: "1", Name: "Sample", CreatedAt: time.Now()}}
	repo.On("GetCount", mock.Anything).Return(1, nil)
	repo.On("GetByPage", mock.Anything, int64(0), int64(20)).Return(items, nil)

	resp, err := svc.List(context.Background(), 0, 20)

	assert.NoError(t, err)
	assert.Equal(t, 1, resp.TotalItems)
	assert.Len(t, resp.Items, 1)
}
