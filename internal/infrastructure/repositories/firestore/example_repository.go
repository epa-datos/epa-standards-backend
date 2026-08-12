package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/epa-datos/epa-standards-backend/internal/pkg/entity"
	"github.com/epa-datos/epa-standards-backend/internal/pkg/ports"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ExampleRepository implements ports.ExampleRepository using a Firestore
// collection instead of a SQL table. Notice it satisfies the *exact same*
// interface as postgres.ExampleRepository — the service layer never knows
// (or cares) which one is wired in main.go.
type ExampleRepository struct {
	client *firestore.Client
}

// NewExampleRepository builds a Firestore-backed ports.ExampleRepository.
func NewExampleRepository(client *firestore.Client) ports.ExampleRepository {
	return &ExampleRepository{client: client}
}

func (r *ExampleRepository) collection() *firestore.CollectionRef {
	return r.client.Collection(examplesCollection)
}

func (r *ExampleRepository) GetByID(ctx context.Context, id string) (*entity.Example, error) {
	doc, err := r.collection().Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var e entity.Example
	if err := doc.DataTo(&e); err != nil {
		return nil, err
	}
	e.ID = doc.Ref.ID
	return &e, nil
}

func (r *ExampleRepository) GetByPage(ctx context.Context, offset, limit int64) ([]*entity.Example, error) {
	iter := r.collection().OrderBy("createdAt", firestore.Desc).Offset(int(offset)).Limit(int(limit)).Documents(ctx)
	defer iter.Stop()

	var items []*entity.Example
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var e entity.Example
		if err := doc.DataTo(&e); err != nil {
			return nil, err
		}
		e.ID = doc.Ref.ID
		items = append(items, &e)
	}
	return items, nil
}

func (r *ExampleRepository) GetCount(ctx context.Context) (int, error) {
	docs, err := r.collection().Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	return len(docs), nil
}

func (r *ExampleRepository) Create(ctx context.Context, example *entity.Example) error {
	ref := r.collection().NewDoc()
	if _, err := ref.Set(ctx, example); err != nil {
		return err
	}
	example.ID = ref.ID
	return nil
}

func (r *ExampleRepository) Update(ctx context.Context, example *entity.Example) error {
	_, err := r.collection().Doc(example.ID).Set(ctx, example)
	return err
}

func (r *ExampleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection().Doc(id).Delete(ctx)
	return err
}
