package rules

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Grandmother shares are established by the Sunnah at one sixth. The maternal
// grandmother is excluded by the mother; the paternal grandmother is excluded
// by the mother and by the father. Those exclusions are applied by the blocking
// stage, so the rule here records only the prescribed one sixth.
//
// When more than one grandmother inherits at the same level they share a single
// one sixth rather than one sixth each. That pooling of an equal furud across a
// group is applied during share assembly, not here.
func init() {
	maternal := FixedShareRule{
		Share:     rational.New(1, 6),
		Condition: "grandmother present and not blocked",
		Reference: "Sunnah (grandmother one sixth)",
		When:      func(Context) bool { return true },
	}
	register(heir.MaternalGrandmother, maternal)
	register(heir.PaternalGrandmother, maternal)
}
