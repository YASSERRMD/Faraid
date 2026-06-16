package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres is a Store backed by PostgreSQL through a pgx connection pool.
type Postgres struct {
	pool *pgxpool.Pool
}

// NewPostgres opens a connection pool to the given DSN.
func NewPostgres(ctx context.Context, dsn string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}

// Close releases the connection pool.
func (p *Postgres) Close() {
	p.pool.Close()
}

// nullableJSON returns nil for empty JSON so it is stored as SQL NULL.
func nullableJSON(j json.RawMessage) []byte {
	if len(j) == 0 {
		return nil
	}
	return j
}

// SaveCase inserts a new case and returns it with its generated id.
func (p *Postgres) SaveCase(ctx context.Context, name string, input, result json.RawMessage) (SavedCase, error) {
	sc := SavedCase{
		ID:        newID(),
		Name:      name,
		Input:     input,
		Result:    result,
		CreatedAt: time.Now().UTC(),
	}
	_, err := p.pool.Exec(ctx,
		`INSERT INTO cases (id, name, input, result, created_at) VALUES ($1, $2, $3, $4, $5)`,
		sc.ID, sc.Name, nullableJSON(sc.Input), nullableJSON(sc.Result), sc.CreatedAt)
	if err != nil {
		return SavedCase{}, err
	}
	return sc, nil
}

// GetCase returns the case with the given id, or ErrNotFound.
func (p *Postgres) GetCase(ctx context.Context, id string) (SavedCase, error) {
	var (
		sc            SavedCase
		input, result []byte
	)
	err := p.pool.QueryRow(ctx,
		`SELECT id, name, input, result, created_at FROM cases WHERE id = $1`, id).
		Scan(&sc.ID, &sc.Name, &input, &result, &sc.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return SavedCase{}, ErrNotFound
	}
	if err != nil {
		return SavedCase{}, err
	}
	sc.Input, sc.Result = input, result
	return sc, nil
}

// ListCases returns all cases, newest first.
func (p *Postgres) ListCases(ctx context.Context) ([]SavedCase, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT id, name, input, result, created_at FROM cases ORDER BY created_at DESC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SavedCase
	for rows.Next() {
		var (
			sc            SavedCase
			input, result []byte
		)
		if err := rows.Scan(&sc.ID, &sc.Name, &input, &result, &sc.CreatedAt); err != nil {
			return nil, err
		}
		sc.Input, sc.Result = input, result
		out = append(out, sc)
	}
	return out, rows.Err()
}

// DeleteCase removes the case with the given id, or returns ErrNotFound.
func (p *Postgres) DeleteCase(ctx context.Context, id string) error {
	tag, err := p.pool.Exec(ctx, `DELETE FROM cases WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
