package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

// ErrNotFound is returned when a case does not exist.
var ErrNotFound = errors.New("store: case not found")

// SavedCase is a persisted case: the input problem and a snapshot of the
// computed result, with the derivation, both as opaque JSON so the store stays
// decoupled from the API and core types.
type SavedCase struct {
	ID        string
	Name      string
	Input     json.RawMessage
	Result    json.RawMessage
	CreatedAt time.Time
}

// Store persists saved cases.
type Store interface {
	SaveCase(ctx context.Context, name string, input, result json.RawMessage) (SavedCase, error)
	GetCase(ctx context.Context, id string) (SavedCase, error)
	ListCases(ctx context.Context) ([]SavedCase, error)
	DeleteCase(ctx context.Context, id string) error
}

// newID returns a random 128 bit identifier as hex.
func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
