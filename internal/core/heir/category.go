package heir

import "sort"

// Sex returns the legal sex of the relation, or SexUnknown when r is not a
// recognized relation.
func (r Relation) Sex() Sex {
	return descriptors[r].sex
}

// Category returns the structural category of the relation. The result is only
// meaningful when r is valid; callers that may hold an invalid relation should
// check Valid first.
func (r Relation) Category() Category {
	return descriptors[r].category
}

// MaxCount returns the maximum number of heirs allowed in this slot, or 0 when
// the count is unbounded.
func (r Relation) MaxCount() int {
	return descriptors[r].maxCount
}

// IsSpouse reports whether r is a spouse relation.
func (r Relation) IsSpouse() bool { return r.Valid() && r.Category() == CategorySpouse }

// IsAscendant reports whether r is an ascendant relation.
func (r Relation) IsAscendant() bool { return r.Valid() && r.Category() == CategoryAscendant }

// IsDescendant reports whether r is a descendant relation.
func (r Relation) IsDescendant() bool { return r.Valid() && r.Category() == CategoryDescendant }

// IsSibling reports whether r is a sibling relation.
func (r Relation) IsSibling() bool { return r.Valid() && r.Category() == CategorySibling }

// IsCollateral reports whether r is a collateral relation.
func (r Relation) IsCollateral() bool { return r.Valid() && r.Category() == CategoryCollateral }

// AllRelations returns every recognized relation in canonical (enum) order.
func AllRelations() []Relation {
	rs := make([]Relation, 0, len(descriptors))
	for r := range descriptors {
		rs = append(rs, r)
	}
	sort.Slice(rs, func(i, j int) bool { return rs[i] < rs[j] })
	return rs
}

// ParseRelation returns the relation whose label matches name (the same label
// String returns), or ok false when no relation matches.
func ParseRelation(name string) (Relation, bool) {
	for r, d := range descriptors {
		if d.name == name {
			return r, true
		}
	}
	return RelationInvalid, false
}
