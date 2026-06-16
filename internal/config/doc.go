// Package config loads and validates runtime configuration from the
// environment at boot, failing fast when a value is invalid or a required
// value is missing. It also builds the application's structured slog.Logger
// from that configuration.
package config
