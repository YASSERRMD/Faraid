package estate

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func TestPreDistributeNoDeductions(t *testing.T) {
	d, err := Estate{Total: 1000}.PreDistribute()
	if err != nil {
		t.Fatal(err)
	}
	if !d.Distributable.Equal(rational.FromInt(1000)) {
		t.Errorf("distributable = %s, want 1000", d.Distributable)
	}
	if d.BequestsClamped {
		t.Error("nothing should be clamped")
	}
}

func TestPreDistributeFuneralAndDebts(t *testing.T) {
	d, err := Estate{Total: 900, Funeral: 100, Debts: 200}.PreDistribute()
	if err != nil {
		t.Fatal(err)
	}
	if !d.AfterFuneralAndDebts.Equal(rational.FromInt(600)) {
		t.Errorf("after funeral/debts = %s, want 600", d.AfterFuneralAndDebts)
	}
	if !d.Distributable.Equal(rational.FromInt(600)) {
		t.Errorf("distributable = %s, want 600", d.Distributable)
	}
}

func TestPreDistributeInsolvent(t *testing.T) {
	d, err := Estate{Total: 100, Funeral: 80, Debts: 200}.PreDistribute()
	if err != nil {
		t.Fatal(err)
	}
	if !d.AfterFuneralAndDebts.IsZero() || !d.Distributable.IsZero() {
		t.Errorf("insolvent estate should leave nothing: %+v", d)
	}
}

func TestBequestWithinCap(t *testing.T) {
	// net 600, cap 200, bequest 150 within cap.
	d, _ := Estate{Total: 600, Bequests: 150}.PreDistribute()
	if d.BequestsClamped {
		t.Error("bequest within cap should not be clamped")
	}
	if !d.EffectiveBequests.Equal(rational.FromInt(150)) {
		t.Errorf("effective bequest = %s, want 150", d.EffectiveBequests)
	}
	if !d.Distributable.Equal(rational.FromInt(450)) {
		t.Errorf("distributable = %s, want 450", d.Distributable)
	}
}

func TestBequestOverCapNoConsent(t *testing.T) {
	// net 600, cap 200, bequest 300 over cap, no consent -> clamp to 200.
	d, _ := Estate{Total: 600, Bequests: 300}.PreDistribute()
	if !d.BequestsClamped {
		t.Error("bequest over cap without consent should be clamped")
	}
	if !d.EffectiveBequests.Equal(rational.FromInt(200)) {
		t.Errorf("effective bequest = %s, want 200", d.EffectiveBequests)
	}
	if !d.Distributable.Equal(rational.FromInt(400)) {
		t.Errorf("distributable = %s, want 400", d.Distributable)
	}
}

func TestBequestOverCapWithConsent(t *testing.T) {
	// net 600, bequest 300 over cap but heirs consent -> not clamped.
	d, _ := Estate{Total: 600, Bequests: 300, HeirsConsentToExcessBequest: true}.PreDistribute()
	if d.BequestsClamped {
		t.Error("bequest with consent should not be clamped")
	}
	if !d.EffectiveBequests.Equal(rational.FromInt(300)) {
		t.Errorf("effective bequest = %s, want 300", d.EffectiveBequests)
	}
	if !d.Distributable.Equal(rational.FromInt(300)) {
		t.Errorf("distributable = %s, want 300", d.Distributable)
	}
}

func TestBequestExceedsEstateWithConsent(t *testing.T) {
	// net 600, bequest 1000 exceeds estate, consent -> clamp to net, nothing left.
	d, _ := Estate{Total: 600, Bequests: 1000, HeirsConsentToExcessBequest: true}.PreDistribute()
	if !d.EffectiveBequests.Equal(rational.FromInt(600)) {
		t.Errorf("effective bequest = %s, want 600", d.EffectiveBequests)
	}
	if !d.Distributable.IsZero() {
		t.Errorf("distributable = %s, want 0", d.Distributable)
	}
}

func TestBequestCapFractional(t *testing.T) {
	// net 100, cap 100/3, bequest 50 over cap, no consent.
	d, _ := Estate{Total: 100, Bequests: 50}.PreDistribute()
	if !d.BequestCap.Equal(rational.New(100, 3)) {
		t.Errorf("cap = %s, want 100/3", d.BequestCap)
	}
	if !d.EffectiveBequests.Equal(rational.New(100, 3)) {
		t.Errorf("effective bequest = %s, want 100/3", d.EffectiveBequests)
	}
	if !d.Distributable.Equal(rational.New(200, 3)) {
		t.Errorf("distributable = %s, want 200/3", d.Distributable)
	}
}

func TestPreDistributeInvalid(t *testing.T) {
	if _, err := (Estate{Total: -1}).PreDistribute(); err == nil {
		t.Error("negative total should error")
	}
}

func TestCaseValidate(t *testing.T) {
	maleCase := Case{
		DeceasedSex: heir.Male,
		Estate:      Estate{Total: 1000},
		Heirs:       heir.New().Set(heir.Wife, 1).Set(heir.Son, 2),
	}
	if err := maleCase.Validate(); err != nil {
		t.Errorf("valid male case rejected: %v", err)
	}

	femaleCase := Case{
		DeceasedSex: heir.Female,
		Estate:      Estate{Total: 1000},
		Heirs:       heir.New().Set(heir.Husband, 1).Set(heir.Daughter, 1),
	}
	if err := femaleCase.Validate(); err != nil {
		t.Errorf("valid female case rejected: %v", err)
	}

	bad := []struct {
		name string
		c    Case
	}{
		{"unknown sex", Case{DeceasedSex: heir.SexUnknown, Heirs: heir.New()}},
		{"negative estate", Case{DeceasedSex: heir.Male, Estate: Estate{Total: -5}, Heirs: heir.New()}},
		{"nil heirs", Case{DeceasedSex: heir.Male}},
		{"invalid heirs", Case{DeceasedSex: heir.Male, Heirs: heir.New().Set(heir.Wife, 9)}},
		{"male with husband", Case{DeceasedSex: heir.Male, Heirs: heir.New().Set(heir.Husband, 1)}},
		{"female with wife", Case{DeceasedSex: heir.Female, Heirs: heir.New().Set(heir.Wife, 1)}},
	}
	for _, c := range bad {
		if err := c.c.Validate(); err == nil {
			t.Errorf("%s should be invalid", c.name)
		}
	}
}
