package link

import (
	"context"
	"sync"
)

type MemoryStorage struct {
	mu    sync.RWMutex
	links map[string]*Link
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		links: make(map[string]*Link),
	}
}

func (m *MemoryStorage) Save(ctx context.Context, l *Link) error {

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.links[l.Code]; exists {
		return ErrCodeExists
	}
	m.links[l.Code] = l

	return nil
}

func (m *MemoryStorage) GetByCode(ctx context.Context, code string) (*Link, error) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	l, ok := m.links[code]

	if !ok {
		return nil, ErrNotFound
	}
	return l, nil
}

func (m *MemoryStorage) IncrementClicks(ctx context.Context, code string) error {

	m.mu.Lock()
	defer m.mu.Unlock()

	l, ok := m.links[code]

	if !ok {
		return ErrNotFound
	}
	l.Clicks++

	return nil
}
