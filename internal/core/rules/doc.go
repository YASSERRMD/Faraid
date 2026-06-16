// Package rules holds the declarative legal data of the engine: fixed share
// tables (ashab al-furud), the blocking lattice (hajb), and per-madhhab rule
// variants.
//
// Rules that differ by school are encoded as data, not hardcoded branches, so
// the engine stays madhhab-aware. These tables are introduced from Phase 7
// onward.
package rules
