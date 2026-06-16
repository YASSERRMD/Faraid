// Package derivation emits the structured, step-by-step audit trail for every
// result.
//
// Each solver stage appends a typed step recording its inputs, the rule
// reference, and the output fraction, so a result can be fully explained and
// independently verified. The emitter is introduced in Phase 28.
package derivation
