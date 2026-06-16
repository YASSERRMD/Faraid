package store

import (
	"context"
	"encoding/json"
	"sort"
	"sync"
	"time"
)

// Memory is an in-memory Store, used in tests and for running without a
// database.
type Memory struct {
	mu    sync.Mutex
	cases map[string]SavedCase
	now   func() time.Time
}

// NewMemory returns an empty in-memory store.
func NewMemory() *Memory {
	return &Memory{cases: make(map[string]SavedCase), now: time.Now}
}

// SaveCase stores a new case and returns it with its generated id.
func (m *Memory) SaveCase(_ context.Context, name string, input, result json.RawMessage) (SavedCase, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sc := SavedCase{
		ID:        newID(),
		Name:      name,
		Input:     input,
		Result:    result,
		CreatedAt: m.now(),
	}
	m.cases[sc.ID] = sc
	return sc, nil
}

// GetCase returns the case with the given id, or ErrNotFound.
func (m *Memory) GetCase(_ context.Context, id string) (SavedCase, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sc, ok := m.cases[id]
	if !ok {
		return SavedCase{}, ErrNotFound
	}
	return sc, nil
}

// ListCases returns all cases, newest first.
func (m *Memory) ListCases(_ context.Context) ([]SavedCase, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]SavedCase, 0, len(m.cases))
	for _, sc := range m.cases {
		out = append(out, sc)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// DeleteCase removes the case with the given id, or returns ErrNotFound.
func (m *Memory) DeleteCase(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.cases[id]; !ok {
		return ErrNotFound
	}
	delete(m.cases, id)
	return nil
}
