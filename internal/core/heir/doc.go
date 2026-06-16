// Package heir models recognized heirs and the family relationship graph.
//
// The Relation type enumerates every heir slot recognized by the four Sunni
// schools as a fixed share heir (ashab al-furud) or residuary heir (asaba),
// grouped into spouse, ascendant, descendant, sibling, and collateral
// categories. Each relation carries static metadata: a label, its legal sex,
// its category, and the maximum count its slot may hold.
//
// The Heirs collection records how many individuals occupy each slot and
// reports presence. Validate rejects inputs that cannot describe a real family
// (negative counts, too many wives, both a husband and a wife) while leaving
// legal rulings such as blocking to the solver.
package heir
