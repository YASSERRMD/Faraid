package solver

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// sc builds a labeled scenario.
func sc(name string, shares map[string]rational.Fraction) Scenario {
	return Scenario{Name: name, Shares: shares}
}

func TestReserveKhuntha(t *testing.T) {
	// A son with an intersex heir: as male they split equally, as female the
	// son takes two thirds. Each is guaranteed the minimum; the difference is
	// reserved until the sex is clarified.
	res := Reserve([]Scenario{
		sc("male", map[string]rational.Fraction{"son": rational.New(1, 2), "khuntha": rational.New(1, 2)}),
		sc("female", map[string]rational.Fraction{"son": rational.New(2, 3), "khuntha": rational.New(1, 3)}),
	})
	if !res.Guaranteed["son"].Equal(rational.New(1, 2)) {
		t.Errorf("son guaranteed = %s, want 1/2", res.Guaranteed["son"])
	}
	if !res.Guaranteed["khuntha"].Equal(rational.New(1, 3)) {
		t.Errorf("khuntha guaranteed = %s, want 1/3", res.Guaranteed["khuntha"])
	}
	if !res.Reserved.Equal(rational.New(1, 6)) {
		t.Errorf("reserved = %s, want 1/6", res.Reserved)
	}
}

func TestReserveMafqud(t *testing.T) {
	// A daughter with a missing son. If alive he takes the residue; if dead she
	// inherits the whole by radd. She is guaranteed only her with-son share.
	res := Reserve([]Scenario{
		sc("alive", map[string]rational.Fraction{"daughter": rational.New(1, 3), "son": rational.New(2, 3)}),
		sc("dead", map[string]rational.Fraction{"daughter": rational.One()}),
	})
	if !res.Guaranteed["daughter"].Equal(rational.New(1, 3)) {
		t.Errorf("daughter guaranteed = %s, want 1/3", res.Guaranteed["daughter"])
	}
	if !res.Guaranteed["son"].IsZero() {
		t.Errorf("missing son guaranteed = %s, want 0", res.Guaranteed["son"])
	}
	if !res.Reserved.Equal(rational.New(2, 3)) {
		t.Errorf("reserved = %s, want 2/3", res.Reserved)
	}
}

func TestReserveHaml(t *testing.T) {
	// Wife, father, and an unborn child. Across son, daughter, and stillborn the
	// wife is guaranteed 1/8 and the father 1/6, with 17/24 reserved until birth.
	res := Reserve([]Scenario{
		sc("son", map[string]rational.Fraction{
			"wife": rational.New(1, 8), "father": rational.New(1, 6), "son": rational.New(17, 24),
		}),
		sc("daughter", map[string]rational.Fraction{
			"wife": rational.New(1, 8), "father": rational.New(3, 8), "daughter": rational.New(1, 2),
		}),
		sc("stillborn", map[string]rational.Fraction{
			"wife": rational.New(1, 4), "father": rational.New(3, 4),
		}),
	})
	if !res.Guaranteed["wife"].Equal(rational.New(1, 8)) {
		t.Errorf("wife guaranteed = %s, want 1/8", res.Guaranteed["wife"])
	}
	if !res.Guaranteed["father"].Equal(rational.New(1, 6)) {
		t.Errorf("father guaranteed = %s, want 1/6", res.Guaranteed["father"])
	}
	if !res.Reserved.Equal(rational.New(17, 24)) {
		t.Errorf("reserved = %s, want 17/24", res.Reserved)
	}
}

func TestReserveEmpty(t *testing.T) {
	res := Reserve(nil)
	if len(res.Guaranteed) != 0 {
		t.Errorf("no scenarios should guarantee nothing, got %v", res.Guaranteed)
	}
	if !res.Reserved.IsWhole() {
		t.Errorf("with no scenarios the whole estate is reserved, got %s", res.Reserved)
	}
}
