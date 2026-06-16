// Package store persists saved cases. The Store interface has an in-memory
// implementation used in tests and for running without a database, and a
// PostgreSQL implementation backed by pgx. Schema changes are managed as goose
// migrations under migrations/. Cases keep their input problem and a snapshot
// of the computed result, including the derivation, as JSONB.
package store
