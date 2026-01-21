package main

import (
	"context"
	"embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"

	"github.com/jwhumphries/blog/internal/metrics"
	"github.com/jwhumphries/blog/internal/server"
	"github.com/jwhumphries/blog/version"
)

//go:embed public
var publicFS embed.FS

func main() {
	// Configure logger
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		ReportCaller:    false,
	})

	// Set log level from environment
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	default:
		logger.SetLevel(log.InfoLevel)
	}

	logger.Info("starting blog server",
		"version", version.Tag,
		"commit", version.Commit,
	)

	// Initialize the server with embedded files
	srv, err := server.New(publicFS, "public", logger)
	if err != nil {
		logger.Fatal("failed to initialize server", "error", err)
	}

	// Configure HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.Handle("GET /health", server.HealthHandler())
	mux.Handle("GET /", srv.Handler())

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      metrics.Middleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Configure metrics server
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "9101"
	}

	metricsMux := http.NewServeMux()
	metricsMux.Handle("GET /metrics", metrics.Handler())

	metricsServer := &http.Server{
		Addr:         ":" + metricsPort,
		Handler:      metricsMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start metrics server
	go func() {
		logger.Info("metrics server listening", "port", metricsPort)
		if err := metricsServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("metrics server error", "error", err)
		}
	}()

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("shutting down servers")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := metricsServer.Shutdown(ctx); err != nil {
			logger.Error("metrics server shutdown error", "error", err)
		}
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Error("http server shutdown error", "error", err)
		}
	}()

	logger.Info("server listening", "port", port)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("server error", "error", err)
	}

	logger.Info("server stopped")
}
