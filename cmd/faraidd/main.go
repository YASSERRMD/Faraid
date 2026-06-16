// Command faraidd is the HTTP server entrypoint for Faraid, an Islamic
// inheritance (ilm al-faraid) calculation system.
//
// The deterministic legal engine lives under internal/core and never performs
// I/O or talks to an LLM. This entrypoint loads and validates configuration,
// sets up structured logging, and (in a later phase) serves the HTTP API.
package main

import (
	"fmt"
	"os"

	"github.com/YASSERRMD/Faraid/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "faraidd: %v\n", err)
		os.Exit(1)
	}

	logger := config.NewLogger(cfg)
	logger.Info("faraidd starting",
		"env", cfg.Env,
		"http_addr", cfg.HTTPAddr,
		"log_level", cfg.LogLevel,
		"log_format", cfg.LogFormat,
	)

	// The HTTP server is wired in a later phase. For now the binary validates
	// its configuration, emits a structured startup log, and exits cleanly.
	logger.Info("faraidd configuration loaded; http server wiring arrives in a later phase")
}
