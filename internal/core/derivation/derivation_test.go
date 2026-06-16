package derivation

import (
	"strings"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func TestStepString(t *testing.T) {
	s := Step{
		Stage:     StageFixedShare,
		Relation:  heir.Husband,
		Detail:    "an inheriting descendant is present",
		Reference: "Quran 4:12",
		Fraction:  rational.New(1, 4),
	}
	got := s.String()
	for _, want := range []string{"[fixed-share]", "husband", "an inheriting descendant", "= 1/4", "(Quran 4:12)"} {
		if !strings.Contains(got, want) {
			t.Errorf("step string %q missing %q", got, want)
		}
	}
}

func TestStepStringMinimal(t *testing.T) {
	// A step with no heir and no fraction omits those parts.
	s := Step{Stage: StageEstate, Detail: "distributable estate"}
	got := s.String()
	if got != "[estate]: distributable estate" {
		t.Errorf("minimal step = %q", got)
	}
}

func TestDerivationString(t *testing.T) {
	d := New()
	d.Add(Step{Stage: StageEstate, Detail: "estate ready"})
	d.Add(Step{Stage: StageResult, Relation: heir.Son, Fraction: rational.New(1, 2)})
	out := d.String()
	if len(d.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(d.Steps))
	}
	if !strings.Contains(out, "estate ready") || !strings.Contains(out, "[result] son = 1/2") {
		t.Errorf("derivation string wrong:\n%s", out)
	}
	if strings.Count(out, "\n") != 1 {
		t.Errorf("expected 2 lines, got %q", out)
	}
}
