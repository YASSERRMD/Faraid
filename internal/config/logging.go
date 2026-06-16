package config

import (
	"io"
	"log/slog"
	"os"
)

// NewLogger builds a structured slog.Logger from the configuration, writing to
// standard output. The handler is JSON unless the configured format is text.
func NewLogger(c *Config) *slog.Logger {
	return newLogger(c, os.Stdout)
}

// newLogger builds a logger writing to w. It is separated from NewLogger so
// tests can capture the output.
func newLogger(c *Config, w io.Writer) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(c.LogLevel)}
	var h slog.Handler
	if c.LogFormat == "text" {
		h = slog.NewTextHandler(w, opts)
	} else {
		h = slog.NewJSONHandler(w, opts)
	}
	return slog.New(h)
}

// parseLevel maps a configured level name to an slog.Level, defaulting to info
// for unrecognized names.
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
