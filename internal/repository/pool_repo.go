package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/mgmacri/pool-maintenance-app/internal/domain"
)

var (
	ErrPoolNotFound = errors.New("pool not found")
)

// PoolRepository defines access to pool data.
type PoolRepository interface {
	Create(p domain.Pool) error
	GetByID(id string) (domain.Pool, error)
	List() ([]domain.Pool, error)
}

// InMemoryPoolRepo is a thread-safe mock implementation.
type InMemoryPoolRepo struct {
	mu    sync.RWMutex
	store map[string]domain.Pool
}

func NewInMemoryPoolRepo() *InMemoryPoolRepo {
	return &InMemoryPoolRepo{
		store: make(map[string]domain.Pool),
	}
}

func (r *InMemoryPoolRepo) Create(p domain.Pool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.store[p.ID]; exists {
		return errors.New("pool already exists")
	}
	p.CreatedAt = time.Now().UTC()
	r.store[p.ID] = p
	return nil
}

func (r *InMemoryPoolRepo) GetByID(id string) (domain.Pool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[id]
	if !ok {
		return domain.Pool{}, ErrPoolNotFound
	}
	return p, nil
}

func (r *InMemoryPoolRepo) List() ([]domain.Pool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pools := make([]domain.Pool, 0, len(r.store))
	for _, p := range r.store {
		pools = append(pools, p)
	}
	return pools, nil
}
