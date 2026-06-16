package rules

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func mustNoShare(t *testing.T, r heir.Relation, h *heir.Heirs) {
	t.Helper()
	if share, _, ok := FixedShare(r, ctx(h)); ok {
		t.Errorf("%s: expected no fixed share, got %s", r, share)
	}
}

func TestDaughterShares(t *testing.T) {
	mustShare(t, heir.Daughter, heir.New().Set(heir.Daughter, 1), rational.New(1, 2))
	mustShare(t, heir.Daughter, heir.New().Set(heir.Daughter, 2), rational.New(2, 3))
	mustShare(t, heir.Daughter, heir.New().Set(heir.Daughter, 3), rational.New(2, 3))
	// With a son the daughter is residuary, not a fixed-share heir.
	mustNoShare(t, heir.Daughter, heir.New().Set(heir.Daughter, 1).Set(heir.Son, 1))
}

func TestSonsDaughterStandalone(t *testing.T) {
	mustShare(t, heir.SonsDaughter, heir.New().Set(heir.SonsDaughter, 1), rational.New(1, 2))
	mustShare(t, heir.SonsDaughter, heir.New().Set(heir.SonsDaughter, 2), rational.New(2, 3))
}

func TestSonsDaughterCompletion(t *testing.T) {
	// One daughter (1/2) plus a son's daughter: she completes the group to two
	// thirds with one sixth.
	h := heir.New().Set(heir.Daughter, 1).Set(heir.SonsDaughter, 1)
	mustShare(t, heir.Daughter, h, rational.New(1, 2))
	mustShare(t, heir.SonsDaughter, h, rational.New(1, 6))
	// Two son's daughters still take one sixth collectively.
	h2 := heir.New().Set(heir.Daughter, 1).Set(heir.SonsDaughter, 2)
	mustShare(t, heir.SonsDaughter, h2, rational.New(1, 6))
}

func TestSonsDaughterBlockedOrResiduary(t *testing.T) {
	// Two daughters consume two thirds; with no son's son the son's daughter
	// is blocked and has no fixed share.
	mustNoShare(t, heir.SonsDaughter, heir.New().Set(heir.Daughter, 2).Set(heir.SonsDaughter, 1))
	// A son blocks the son's daughter entirely.
	mustNoShare(t, heir.SonsDaughter, heir.New().Set(heir.Son, 1).Set(heir.SonsDaughter, 1))
	// A son's son makes the son's daughter residuary, so she has no fixed share.
	mustNoShare(t, heir.SonsDaughter, heir.New().Set(heir.SonsSon, 1).Set(heir.SonsDaughter, 1))
	// Two daughters with a son's son present: residuary, no fixed share.
	mustNoShare(t, heir.SonsDaughter, heir.New().Set(heir.Daughter, 2).Set(heir.SonsDaughter, 1).Set(heir.SonsSon, 1))
}
