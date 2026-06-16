// Package rules holds the declarative legal data of the engine: the fixed
// share table (ashab al-furud), and in later phases the blocking lattice
// (hajb) and per-madhhab rule variants.
//
// Each fixed-share heir has an ordered list of FixedShareRule entries. Given a
// Context describing the heirs effectively present, FixedShare returns the
// first rule whose condition holds, together with its prescribed fraction and
// a source reference for the audit trail. Rules that differ by school are
// encoded as data rather than as hardcoded branches so the engine stays
// madhhab-aware.
package rules
