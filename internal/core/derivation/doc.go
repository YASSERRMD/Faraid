// Package derivation emits the structured, step-by-step audit trail for every
// result.
//
// A Derivation is an ordered list of typed Steps, each recording one stage of
// the computation: the heir it concerns, a human-readable detail, a source
// reference, and the resulting fraction. The solver builds a Derivation
// alongside each result so the outcome can be explained and independently
// verified.
package derivation
