package derivation

import (
	"fmt"
	"strings"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Stage names a phase of the computation in the audit trail.
type Stage string

const (
	StageEstate      Stage = "estate"
	StageSpecialCase Stage = "special-case"
	StageBlocking    Stage = "blocking"
	StageFixedShare  Stage = "fixed-share"
	StageResiduary   Stage = "residuary"
	StageAwl         Stage = "awl"
	StageRadd        Stage = "radd"
	StageAsl         Stage = "asl"
	StageResult      Stage = "result"
)

// Step is one entry in the derivation. Relation is RelationInvalid when the
// step does not concern a specific heir; Fraction is zero when the step carries
// no fraction.
type Step struct {
	Stage     Stage
	Relation  heir.Relation
	Detail    string
	Reference string
	Fraction  rational.Fraction
}

// String renders a step as a single human-readable line.
func (s Step) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "[%s]", s.Stage)
	if s.Relation.Valid() {
		fmt.Fprintf(&b, " %s", s.Relation)
	}
	if s.Detail != "" {
		fmt.Fprintf(&b, ": %s", s.Detail)
	}
	if !s.Fraction.IsZero() {
		fmt.Fprintf(&b, " = %s", s.Fraction)
	}
	if s.Reference != "" {
		fmt.Fprintf(&b, " (%s)", s.Reference)
	}
	return b.String()
}

// Derivation is the ordered audit trail of a result.
type Derivation struct {
	Steps []Step
}

// New returns an empty derivation.
func New() *Derivation {
	return &Derivation{}
}

// Add appends a step to the trail.
func (d *Derivation) Add(s Step) {
	d.Steps = append(d.Steps, s)
}

// String renders the whole derivation, one step per line.
func (d *Derivation) String() string {
	lines := make([]string, len(d.Steps))
	for i, s := range d.Steps {
		lines[i] = s.String()
	}
	return strings.Join(lines, "\n")
}
