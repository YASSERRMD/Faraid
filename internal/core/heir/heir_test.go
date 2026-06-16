package heir

import "testing"

func TestSexString(t *testing.T) {
	cases := map[Sex]string{Male: "male", Female: "female", SexUnknown: "unknown", Sex(99): "unknown"}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("Sex(%d).String() = %q, want %q", s, got, want)
		}
	}
}

func TestCategoryString(t *testing.T) {
	cases := map[Category]string{
		CategorySpouse:     "spouse",
		CategoryAscendant:  "ascendant",
		CategoryDescendant: "descendant",
		CategorySibling:    "sibling",
		CategoryCollateral: "collateral",
		Category(99):       "unknown",
	}
	for c, want := range cases {
		if got := c.String(); got != want {
			t.Errorf("Category(%d).String() = %q, want %q", c, got, want)
		}
	}
}

func TestRelationValidAndString(t *testing.T) {
	if !Son.Valid() {
		t.Error("Son should be valid")
	}
	if RelationInvalid.Valid() {
		t.Error("RelationInvalid should not be valid")
	}
	if Relation(999).Valid() {
		t.Error("unknown relation should not be valid")
	}
	if got := Daughter.String(); got != "daughter" {
		t.Errorf("Daughter.String() = %q, want daughter", got)
	}
	if got := Relation(999).String(); got != "invalid relation" {
		t.Errorf("unknown relation String() = %q, want invalid relation", got)
	}
}

func TestRelationMetadata(t *testing.T) {
	if Husband.Sex() != Male || Wife.Sex() != Female {
		t.Error("spouse sexes wrong")
	}
	if Relation(999).Sex() != SexUnknown {
		t.Error("unknown relation sex should be SexUnknown")
	}
	if Wife.MaxCount() != 4 {
		t.Errorf("Wife.MaxCount() = %d, want 4", Wife.MaxCount())
	}
	if Father.MaxCount() != 1 {
		t.Errorf("Father.MaxCount() = %d, want 1", Father.MaxCount())
	}
	if Son.MaxCount() != 0 {
		t.Errorf("Son.MaxCount() = %d, want 0 (unbounded)", Son.MaxCount())
	}
}

func TestCategoryHelpers(t *testing.T) {
	if !Husband.IsSpouse() || Husband.IsAscendant() {
		t.Error("Husband category helpers wrong")
	}
	if !Father.IsAscendant() || !Son.IsDescendant() || !FullBrother.IsSibling() || !FullPaternalUncle.IsCollateral() {
		t.Error("category helpers wrong")
	}
	// An invalid relation belongs to no category.
	bad := Relation(999)
	if bad.IsSpouse() || bad.IsAscendant() || bad.IsDescendant() || bad.IsSibling() || bad.IsCollateral() {
		t.Error("invalid relation should not match any category")
	}
}

func TestAllRelations(t *testing.T) {
	rs := AllRelations()
	if len(rs) != 23 {
		t.Errorf("AllRelations length = %d, want 23", len(rs))
	}
	for i := range rs {
		if !rs[i].Valid() {
			t.Errorf("AllRelations contains invalid relation %d", rs[i])
		}
		if i > 0 && rs[i-1] >= rs[i] {
			t.Error("AllRelations is not sorted ascending")
		}
	}
	if rs[0] != Husband {
		t.Errorf("first relation = %v, want husband", rs[0])
	}
}

func TestParseRelation(t *testing.T) {
	for _, r := range AllRelations() {
		got, ok := ParseRelation(r.String())
		if !ok || got != r {
			t.Errorf("ParseRelation(%q) = %v, %v; want %v", r.String(), got, ok, r)
		}
	}
	if _, ok := ParseRelation("third cousin"); ok {
		t.Error("unknown name should not parse")
	}
}

func TestHeirsCollection(t *testing.T) {
	h := New()
	if !h.Empty() {
		t.Error("new Heirs should be empty")
	}
	h.Set(Son, 2).Set(Daughter, 1).Set(Husband, 1)
	if h.Count(Son) != 2 || h.Count(Daughter) != 1 {
		t.Error("counts not recorded")
	}
	if !h.Present(Son) || h.Present(Mother) {
		t.Error("presence wrong")
	}
	// Setting to zero removes the slot.
	h.Set(Husband, 0)
	if h.Present(Husband) {
		t.Error("setting count 0 should remove the slot")
	}
	// Replacement.
	h.Set(Son, 3)
	if h.Count(Son) != 3 {
		t.Error("Set should replace the count")
	}
	rs := h.Relations()
	if len(rs) != 2 || rs[0] != Son || rs[1] != Daughter {
		t.Errorf("Relations = %v, want [son daughter] in order", rs)
	}
	if h.Empty() {
		t.Error("Heirs with members should not be empty")
	}
}

func TestValidate(t *testing.T) {
	// Valid: empty.
	if err := New().Validate(); err != nil {
		t.Errorf("empty heirs should be valid: %v", err)
	}
	// Valid: a typical case.
	good := New().Set(Husband, 1).Set(Mother, 1).Set(Son, 2).Set(Daughter, 1)
	if err := good.Validate(); err != nil {
		t.Errorf("typical case should be valid: %v", err)
	}
	// Valid: father and paternal grandfather together (grandfather is blocked,
	// not impossible).
	if err := New().Set(Father, 1).Set(PaternalGrandfather, 1).Validate(); err != nil {
		t.Errorf("father plus grandfather should be valid: %v", err)
	}

	bad := []struct {
		name  string
		heirs *Heirs
	}{
		{"negative count", New().Set(Son, -1)},
		{"too many wives", New().Set(Wife, 5)},
		{"two fathers", New().Set(Father, 2)},
		{"husband and wife", New().Set(Husband, 1).Set(Wife, 1)},
		{"unknown relation", New().Set(Relation(999), 1)},
	}
	for _, c := range bad {
		if err := c.heirs.Validate(); err == nil {
			t.Errorf("%s should be invalid", c.name)
		}
	}
}
