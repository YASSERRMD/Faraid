package store

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
)

// TestPostgresIntegration exercises the pgx store against a real database. It
// is skipped unless FARAID_TEST_DATABASE_URL points at a PostgreSQL instance,
// since no container is available in the standard test environment.
func TestPostgresIntegration(t *testing.T) {
	dsn := os.Getenv("FARAID_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set FARAID_TEST_DATABASE_URL to run the postgres integration test")
	}
	ctx := context.Background()

	pg, err := NewPostgres(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pg.Close()

	if _, err := pg.pool.Exec(ctx, "DROP TABLE IF EXISTS cases"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if _, err := pg.pool.Exec(ctx, migrationUp(t)); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { _, _ = pg.pool.Exec(ctx, "DROP TABLE IF EXISTS cases") })

	sc, err := pg.SaveCase(ctx, "pg case", json.RawMessage(`{"a":1}`), json.RawMessage(`{"r":1}`))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := pg.GetCase(ctx, sc.ID)
	if err != nil || string(got.Input) != `{"a": 1}` && string(got.Input) != `{"a":1}` {
		t.Errorf("get: %+v %v", got, err)
	}
	list, err := pg.ListCases(ctx)
	if err != nil || len(list) != 1 {
		t.Errorf("list: %d %v", len(list), err)
	}
	if err := pg.DeleteCase(ctx, sc.ID); err != nil {
		t.Errorf("delete: %v", err)
	}
	if _, err := pg.GetCase(ctx, sc.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("get after delete = %v, want ErrNotFound", err)
	}
}

// migrationUp returns the Up section of the first migration, with the goose
// markers stripped.
func migrationUp(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("../../migrations/0001_create_cases.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	s := string(data)
	if i := strings.Index(s, "-- +goose Down"); i >= 0 {
		s = s[:i]
	}
	if i := strings.Index(s, "-- +goose Up"); i >= 0 {
		s = s[i+len("-- +goose Up"):]
	}
	s = strings.ReplaceAll(s, "-- +goose StatementBegin", "")
	s = strings.ReplaceAll(s, "-- +goose StatementEnd", "")
	return s
}
