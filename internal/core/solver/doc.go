// Package solver computes inheritance shares deterministically.
//
// It resolves total exclusion (hajb hirman) and share reduction (hajb nuqsan)
// over the heirs present, then in later phases assigns fixed shares, allocates
// the residue to the residuary heirs, computes asl al-mas'ala, and applies
// awl, radd, tashih, and the classical special cases. The end-to-end Solve
// pipeline is assembled in Phase 27.
package solver
