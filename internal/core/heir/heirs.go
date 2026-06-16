package heir

import "sort"

// Heirs is the set of heirs present in a case, with the count of individuals in
// each relation slot. The zero value is not usable; construct with New.
type Heirs struct {
	counts map[Relation]int
}

// New returns an empty Heirs set.
func New() *Heirs {
	return &Heirs{counts: make(map[Relation]int)}
}

// Set records count heirs in the given relation slot, replacing any previous
// value. A count of zero removes the slot. Set returns the receiver so calls
// can be chained. Invalid relations and negative counts are accepted here and
// reported by Validate, which keeps construction and checking separate.
func (h *Heirs) Set(r Relation, count int) *Heirs {
	if count == 0 {
		delete(h.counts, r)
		return h
	}
	h.counts[r] = count
	return h
}

// Count returns the number of heirs in the given slot, or zero if absent.
func (h *Heirs) Count(r Relation) int {
	return h.counts[r]
}

// Present reports whether at least one heir occupies the slot.
func (h *Heirs) Present(r Relation) bool {
	return h.counts[r] > 0
}

// Relations returns the present relations in canonical (enum) order.
func (h *Heirs) Relations() []Relation {
	rs := make([]Relation, 0, len(h.counts))
	for r, c := range h.counts {
		if c > 0 {
			rs = append(rs, r)
		}
	}
	sort.Slice(rs, func(i, j int) bool { return rs[i] < rs[j] })
	return rs
}

// Empty reports whether no heir is present.
func (h *Heirs) Empty() bool {
	return len(h.Relations()) == 0
}
