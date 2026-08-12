// Package postgres is a reference implementation of ports.ExampleRepository
// backed by PostgreSQL via GORM. Copy this package (client.go + db.go) as the
// starting point for any new SQL-backed repository.
package postgres

import (
	"context"

	"github.com/epa-datos/epa-standards-backend/internal/pkg/entity"
	"github.com/epa-datos/epa-standards-backend/internal/pkg/ports"
	"gorm.io/gorm"
)

// exampleRow is the GORM model for the "examples" table. Keeping it separate
// from entity.Example lets the persistence schema evolve independently from
// the API contract.
type exampleRow struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	CreatedAt   int64
	UpdatedAt   int64
}

func (exampleRow) TableName() string { return "examples" }

// ExampleRepository implements ports.ExampleRepository using PostgreSQL.
type ExampleRepository struct {
	db *gorm.DB
}

// NewExampleRepository builds a Postgres-backed ports.ExampleRepository.
// Wire it in main.go: postgres.NewExampleRepository(postgres.NewClient()).
func NewExampleRepository(db *gorm.DB) ports.ExampleRepository {
	return &ExampleRepository{db: db}
}

func (r *ExampleRepository) GetByID(ctx context.Context, id string) (*entity.Example, error) {
	var row exampleRow
	err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toEntity(&row), nil
}

func (r *ExampleRepository) GetByPage(ctx context.Context, offset, limit int64) ([]*entity.Example, error) {
	var rows []exampleRow
	if err := r.db.WithContext(ctx).Offset(int(offset)).Limit(int(limit)).Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]*entity.Example, 0, len(rows))
	for i := range rows {
		items = append(items, toEntity(&rows[i]))
	}
	return items, nil
}

func (r *ExampleRepository) GetCount(ctx context.Context) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&exampleRow{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *ExampleRepository) Create(ctx context.Context, example *entity.Example) error {
	row := fromEntity(example)
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return err
	}
	example.ID = row.ID
	return nil
}

func (r *ExampleRepository) Update(ctx context.Context, example *entity.Example) error {
	return r.db.WithContext(ctx).Save(fromEntity(example)).Error
}

func (r *ExampleRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&exampleRow{}, "id = ?", id).Error
}

func toEntity(row *exampleRow) *entity.Example {
	return &entity.Example{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
	}
}

func fromEntity(e *entity.Example) *exampleRow {
	return &exampleRow{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
	}
}
