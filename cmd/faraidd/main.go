// Command faraidd is the HTTP server entrypoint for Faraid, an Islamic
// inheritance (ilm al-faraid) calculation system.
//
// The deterministic legal engine lives under internal/core and never performs
// I/O or talks to an LLM. This entrypoint loads and validates configuration,
// selects the backing store, and serves the HTTP API.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YASSERRMD/Faraid/internal/api"
	"github.com/YASSERRMD/Faraid/internal/config"
	"github.com/YASSERRMD/Faraid/internal/store"
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

	ctx := context.Background()

	var st store.Store
	if cfg.DatabaseURL != "" {
		pg, err := store.NewPostgres(ctx, cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "faraidd: store: %v\n", err)
			os.Exit(1)
		}
		defer pg.Close()
		if err := pg.Migrate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "faraidd: migrate: %v\n", err)
			os.Exit(1)
		}
		st = pg
		logger.Info("using postgres store")
	} else {
		st = store.NewMemory()
		logger.Info("using in-memory store (cases are not persisted across restarts)")
	}

	srv := api.NewServerWithStore(st)

	httpSrv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      srv.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("http server listening", "addr", cfg.HTTPAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logger.Info("shutdown signal received", "signal", sig)
	case err := <-errCh:
		logger.Error("server error", "err", err)
	}

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutCtx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
	}
	logger.Info("faraidd stopped")
}
