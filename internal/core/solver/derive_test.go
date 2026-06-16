package solver

import (
	"strings"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/derivation"
	"github.com/YASSERRMD/Faraid/internal/core/heir"
)

// stages returns the ordered list of stage names in a derivation.
func stages(d *derivation.Derivation) []derivation.Stage {
	out := make([]derivation.Stage, len(d.Steps))
	for i, s := range d.Steps {
		out[i] = s.Stage
	}
	return out
}

func hasStage(d *derivation.Derivation, want derivation.Stage) bool {
	for _, s := range d.Steps {
		if s.Stage == want {
			return true
		}
	}
	return false
}

func TestDeriveNormalCase(t *testing.T) {
	// Husband, son, daughter: estate, fixed share, residuary, asl, result steps.
	r := solveCase(t, heir.Female, 2400, heir.New().Set(heir.Husband, 1).Set(heir.Son, 1).Set(heir.Daughter, 1), Hanafi)
	d := r.Derivation
	if d == nil {
		t.Fatal("derivation should be present")
	}
	for _, want := range []derivation.Stage{
		derivation.StageEstate, derivation.StageFixedShare, derivation.StageResiduary,
		derivation.StageAsl, derivation.StageResult,
	} {
		if !hasStage(d, want) {
			t.Errorf("derivation missing stage %q; stages=%v", want, stages(d))
		}
	}
	// The husband's fixed-share step carries its Quranic reference.
	out := d.String()
	if !strings.Contains(out, "Quran 4:12") {
		t.Errorf("derivation missing spouse reference:\n%s", out)
	}
	// The first step is the estate, the result steps come last.
	if d.Steps[0].Stage != derivation.StageEstate {
		t.Errorf("first step = %q, want estate", d.Steps[0].Stage)
	}
}

func TestDeriveBlockingAndSpecialCase(t *testing.T) {
	// Blocking step appears for an excluded heir.
	blocked := solveCase(t, heir.Female, 0, heir.New().Set(heir.Husband, 1).Set(heir.Son, 1).Set(heir.FullBrother, 1), Hanafi)
	if !hasStage(blocked.Derivation, derivation.StageBlocking) {
		t.Error("expected a blocking step")
	}

	// A special case is named and skips the per-heir fixed-share reconstruction.
	akd := solveCase(t, heir.Female, 0, heir.New().
		Set(heir.Husband, 1).Set(heir.Mother, 1).Set(heir.PaternalGrandfather, 1).Set(heir.FullSister, 1), Shafii)
	if !hasStage(akd.Derivation, derivation.StageSpecialCase) {
		t.Error("expected a special-case step")
	}
	if hasStage(akd.Derivation, derivation.StageFixedShare) {
		t.Error("special case should not reconstruct fixed-share steps")
	}
	if !strings.Contains(akd.Derivation.String(), "akdariyyah") {
		t.Errorf("special-case derivation should name akdariyyah:\n%s", akd.Derivation.String())
	}
}

func TestDeriveAwlRaddAndTreasury(t *testing.T) {
	awl := solveCase(t, heir.Female, 0, heir.New().Set(heir.Husband, 1).Set(heir.FullSister, 2), Maliki)
	if !hasStage(awl.Derivation, derivation.StageAwl) {
		t.Error("expected an awl step")
	}
	radd := solveCase(t, heir.Male, 0, heir.New().Set(heir.Mother, 1).Set(heir.Daughter, 1), Hanbali)
	if !hasStage(radd.Derivation, derivation.StageRadd) {
		t.Error("expected a radd step")
	}
	treasury := solveCase(t, heir.Female, 0, heir.New().Set(heir.Husband, 1), Hanafi)
	if !strings.Contains(treasury.Derivation.String(), "treasury") {
		t.Errorf("expected a treasury step:\n%s", treasury.Derivation.String())
	}
}
