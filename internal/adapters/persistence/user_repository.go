package persistence

import (
	"context"
	"sync"

	"epa-api/internal/domain"
	"epa-api/pkg/logger"
)

// UserRepository implements domain.UserRepository using in-memory storage
// In production, this would use a real database (PostgreSQL, BigQuery, etc.)
type UserRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
	log   *logger.Logger
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(log *logger.Logger) *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
		log:   log,
	}
}

// Save stores a user
func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user
	r.log.Printf("User saved: %s", user.ID)
	return nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

// Delete removes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.users, id)
	r.log.Printf("User deleted: %s", id)
	return nil
}

// GetAll returns all users (for testing)
func (r *UserRepository) GetAll() []*domain.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}

// Clear removes all users (for testing)
func (r *UserRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = make(map[string]*domain.User)
}
