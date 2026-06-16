package heir

// Category groups heirs by their structural place in the family graph. It is
// used by later phases to reason about classes of heirs without enumerating
// every relation.
type Category uint8

const (
	CategorySpouse Category = iota
	CategoryAscendant
	CategoryDescendant
	CategorySibling
	CategoryCollateral
)

// String returns a lowercase label for the category.
func (c Category) String() string {
	switch c {
	case CategorySpouse:
		return "spouse"
	case CategoryAscendant:
		return "ascendant"
	case CategoryDescendant:
		return "descendant"
	case CategorySibling:
		return "sibling"
	case CategoryCollateral:
		return "collateral"
	default:
		return "unknown"
	}
}

// Relation identifies a recognized heir slot in the family graph. The values
// cover the heirs recognized by the four Sunni schools as fixed share heirs
// (ashab al-furud) or residuary heirs (asaba). Lineal slots such as the son's
// son and the paternal grandfather stand for the nearest occupant of that line
// (the agnatic grandson how low soever, the grandfather how high soever).
type Relation int

const (
	RelationInvalid Relation = iota

	// Spouses.
	Husband
	Wife

	// Ascendants.
	Father
	Mother
	PaternalGrandfather
	PaternalGrandmother
	MaternalGrandmother

	// Descendants.
	Son
	Daughter
	SonsSon
	SonsDaughter

	// Siblings.
	FullBrother
	FullSister
	ConsanguineBrother
	ConsanguineSister
	UterineBrother
	UterineSister

	// Collaterals.
	FullBrothersSon
	ConsanguineBrothersSon
	FullPaternalUncle
	ConsanguinePaternalUncle
	FullPaternalUnclesSon
	ConsanguinePaternalUnclesSon
)

// unbounded marks a relation whose count has no fixed upper limit.
const unbounded = 0

// descriptor holds the static metadata of a relation.
type descriptor struct {
	name     string
	sex      Sex
	category Category
	maxCount int // 0 means unbounded
}

var descriptors = map[Relation]descriptor{
	Husband: {"husband", Male, CategorySpouse, 1},
	Wife:    {"wife", Female, CategorySpouse, 4},

	Father:              {"father", Male, CategoryAscendant, 1},
	Mother:              {"mother", Female, CategoryAscendant, 1},
	PaternalGrandfather: {"paternal grandfather", Male, CategoryAscendant, 1},
	PaternalGrandmother: {"paternal grandmother", Female, CategoryAscendant, 1},
	MaternalGrandmother: {"maternal grandmother", Female, CategoryAscendant, 1},

	Son:          {"son", Male, CategoryDescendant, unbounded},
	Daughter:     {"daughter", Female, CategoryDescendant, unbounded},
	SonsSon:      {"son's son", Male, CategoryDescendant, unbounded},
	SonsDaughter: {"son's daughter", Female, CategoryDescendant, unbounded},

	FullBrother:        {"full brother", Male, CategorySibling, unbounded},
	FullSister:         {"full sister", Female, CategorySibling, unbounded},
	ConsanguineBrother: {"consanguine brother", Male, CategorySibling, unbounded},
	ConsanguineSister:  {"consanguine sister", Female, CategorySibling, unbounded},
	UterineBrother:     {"uterine brother", Male, CategorySibling, unbounded},
	UterineSister:      {"uterine sister", Female, CategorySibling, unbounded},

	FullBrothersSon:              {"full brother's son", Male, CategoryCollateral, unbounded},
	ConsanguineBrothersSon:       {"consanguine brother's son", Male, CategoryCollateral, unbounded},
	FullPaternalUncle:            {"full paternal uncle", Male, CategoryCollateral, unbounded},
	ConsanguinePaternalUncle:     {"consanguine paternal uncle", Male, CategoryCollateral, unbounded},
	FullPaternalUnclesSon:        {"full paternal uncle's son", Male, CategoryCollateral, unbounded},
	ConsanguinePaternalUnclesSon: {"consanguine paternal uncle's son", Male, CategoryCollateral, unbounded},
}

// Valid reports whether r is a recognized heir relation.
func (r Relation) Valid() bool {
	_, ok := descriptors[r]
	return ok
}

// String returns the English label for r, or a placeholder when r is not a
// recognized relation.
func (r Relation) String() string {
	if d, ok := descriptors[r]; ok {
		return d.name
	}
	return "invalid relation"
}
