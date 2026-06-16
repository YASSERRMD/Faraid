package solver

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func TestMadhahibPresets(t *testing.T) {
	if len(Madhahib()) != 4 {
		t.Fatalf("expected four schools, got %d", len(Madhahib()))
	}
	// The grandfather divergence: only Hanafi excludes the siblings.
	if Hanafi.GrandfatherView != JaddAbuHanifa {
		t.Error("Hanafi should use the Abu Hanifa grandfather view")
	}
	for _, m := range []Madhhab{Maliki, Shafii, Hanbali} {
		if m.GrandfatherView != JaddZayd {
			t.Errorf("%s should use the Zayd grandfather view", m.Name)
		}
	}
	// The mushtaraka divergence: only Maliki and Shafi'i share.
	if Maliki.MushtarakaView != MushtarakaShare || Shafii.MushtarakaView != MushtarakaShare {
		t.Error("Maliki and Shafi'i should share in mushtaraka")
	}
	if Hanafi.MushtarakaView != MushtarakaNoShare || Hanbali.MushtarakaView != MushtarakaNoShare {
		t.Error("Hanafi and Hanbali should not share in mushtaraka")
	}
	// Distant kindred: Hanafi and Hanbali let them inherit.
	if Hanafi.DhawuArham != DhawuArhamInherit || Hanbali.DhawuArham != DhawuArhamInherit {
		t.Error("Hanafi and Hanbali should let distant kindred inherit")
	}
	if Maliki.DhawuArham != DhawuArhamExcluded {
		t.Error("Maliki should exclude distant kindred")
	}
}

func TestMadhhabByName(t *testing.T) {
	if m, ok := MadhhabByName("Shafii"); !ok || m.Name != "Shafii" {
		t.Errorf("lookup Shafii failed: %v %v", m, ok)
	}
	if _, ok := MadhhabByName("Jafari"); ok {
		t.Error("unknown school should not be found")
	}
}

func TestMadhhabDrivesGrandfatherResult(t *testing.T) {
	// The same case resolves differently under the two grandfather views,
	// driven only by the profile data.
	h := heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 1)

	zayd := GrandfatherWithSiblings(rational.One(), h, Maliki.GrandfatherView)
	if zayd.SiblingsExcluded || !zayd.SiblingShares[heir.FullBrother].Equal(rational.New(1, 2)) {
		t.Errorf("Zayd: brother should share, got %v excluded=%v", zayd.SiblingShares, zayd.SiblingsExcluded)
	}

	hanafi := GrandfatherWithSiblings(rational.One(), h, Hanafi.GrandfatherView)
	if !hanafi.SiblingsExcluded {
		t.Error("Hanafi: the grandfather should exclude the brother")
	}
}

func TestDhawuArhamRouteString(t *testing.T) {
	if DhawuArhamInherit.String() != "inherit" || DhawuArhamExcluded.String() != "excluded" {
		t.Error("route labels wrong")
	}
}
