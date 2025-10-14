package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matsuboshi/league-matrix-app/internal/handler"
)

const port = "8080"

func main() {
	matrixHandler := handler.NewMatrixHandler()

	http.HandleFunc("/", matrixHandler.ListMatrixOperations)
	http.HandleFunc("/matrix", matrixHandler.ListMatrixOperations)
	http.HandleFunc("/matrix/", matrixHandler.ProcessMatrix)
	http.HandleFunc("/health", matrixHandler.HealthCheck)

	// Configure HTTP server with timeouts
	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 5 * time.Second,  // Maximum time to read request headers (prevents slow header attacks)
		ReadTimeout:       7 * time.Second,  // Maximum duration for reading the entire request
		WriteTimeout:      30 * time.Second, // Maximum duration before timing out writes
		IdleTimeout:       60 * time.Second, // Maximum time to wait for next request with keep-alive
	}

	slog.Info("starting HTTP server",
		"port", port,
		"address", "http://localhost:"+port,
		"read_timeout", server.ReadTimeout,
		"write_timeout", server.WriteTimeout)

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to start", "error", err, "port", port)
			os.Exit(1)
		}
	}()

	// Setup signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	// Listen for SIGINT (Ctrl+C) and SIGTERM (Docker/K8s stop)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal is received
	sig := <-quit
	slog.Info("shutdown signal received", "signal", sig.String())

	// Create context with timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	slog.Info("gracefully shutting down server", "timeout", "30s")
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}
