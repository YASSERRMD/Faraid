package store

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestMemoryCRUD(t *testing.T) {
	ctx := context.Background()
	m := NewMemory()

	// Inject a clock so the ordering is deterministic.
	t1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Hour)
	times := []time.Time{t1, t2}
	i := 0
	m.now = func() time.Time {
		ti := times[i]
		i++
		return ti
	}

	first, err := m.SaveCase(ctx, "first", json.RawMessage(`{"a":1}`), json.RawMessage(`{"r":1}`))
	if err != nil {
		t.Fatal(err)
	}
	second, _ := m.SaveCase(ctx, "second", json.RawMessage(`{"a":2}`), nil)
	if first.ID == second.ID {
		t.Fatal("ids should be unique")
	}

	got, err := m.GetCase(ctx, first.ID)
	if err != nil || got.Name != "first" || string(got.Input) != `{"a":1}` {
		t.Errorf("get first = %+v, %v", got, err)
	}

	list, err := m.ListCases(ctx)
	if err != nil || len(list) != 2 {
		t.Fatalf("list = %d, %v", len(list), err)
	}
	if list[0].ID != second.ID || list[1].ID != first.ID {
		t.Error("list should be newest first")
	}

	if err := m.DeleteCase(ctx, first.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := m.GetCase(ctx, first.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("get after delete = %v, want ErrNotFound", err)
	}
	if err := m.DeleteCase(ctx, first.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("delete again = %v, want ErrNotFound", err)
	}
}

func TestMemoryNotFound(t *testing.T) {
	if _, err := NewMemory().GetCase(context.Background(), "nope"); !errors.Is(err, ErrNotFound) {
		t.Errorf("get unknown = %v, want ErrNotFound", err)
	}
}

func TestNewIDUnique(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 1000; i++ {
		id := newID()
		if len(id) != 32 {
			t.Fatalf("id length = %d, want 32", len(id))
		}
		if seen[id] {
			t.Fatal("duplicate id")
		}
		seen[id] = true
	}
}
