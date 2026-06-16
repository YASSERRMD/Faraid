package config

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestNewLoggerJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := newLogger(&Config{LogLevel: "info", LogFormat: "json"}, &buf)
	logger.Info("hello", "key", "value")

	var record map[string]any
	if err := json.Unmarshal(buf.Bytes(), &record); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
	if record["msg"] != "hello" {
		t.Errorf("msg = %v, want hello", record["msg"])
	}
	if record["level"] != "INFO" {
		t.Errorf("level = %v, want INFO", record["level"])
	}
	if record["key"] != "value" {
		t.Errorf("key = %v, want value", record["key"])
	}
}

func TestNewLoggerText(t *testing.T) {
	var buf bytes.Buffer
	logger := newLogger(&Config{LogLevel: "info", LogFormat: "text"}, &buf)
	logger.Info("hello")

	out := buf.String()
	if !strings.Contains(out, "msg=hello") || !strings.Contains(out, "level=INFO") {
		t.Errorf("text output missing expected fields: %q", out)
	}
}

func TestNewLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := newLogger(&Config{LogLevel: "error", LogFormat: "json"}, &buf)
	logger.Info("should be filtered")
	if buf.Len() != 0 {
		t.Errorf("info log should be suppressed at error level, got: %s", buf.String())
	}
}

func TestNewLogger(t *testing.T) {
	if NewLogger(&Config{LogLevel: "info", LogFormat: "json"}) == nil {
		t.Error("NewLogger returned nil")
	}
}

func TestParseLevel(t *testing.T) {
	cases := map[string]slog.Level{
		"debug":   slog.LevelDebug,
		"info":    slog.LevelInfo,
		"warn":    slog.LevelWarn,
		"error":   slog.LevelError,
		"unknown": slog.LevelInfo,
	}
	for name, want := range cases {
		if got := parseLevel(name); got != want {
			t.Errorf("parseLevel(%q) = %v, want %v", name, got, want)
		}
	}
}
