package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"block-flow/internal/api"
	"block-flow/internal/config"
	"block-flow/internal/engine"
	"block-flow/internal/storage"
)

func main() {
	// Initialize simple logger
	log.Println("Starting Block-Flow Server v1.0.0")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Server will listen on %s", cfg.Server.Address)

	// Initialize storage
	storage := storage.NewFileStorage(cfg.Storage.DataDir)

	// Initialize flow engine with simple logger
	logger := &engine.SimpleLogger{}
	flowEngine := engine.New(storage, logger)

	// Load and start existing flows on startup
	ctx := context.Background()
	if err := flowEngine.LoadAndStartFlows(ctx); err != nil {
		log.Printf("Warning: Failed to load existing flows: %v", err)
	}

	// Initialize API router
	router := api.NewRouter(flowEngine, storage)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Println("Server started successfully. Press Ctrl+C to stop.")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")

	// Give the server 30 seconds to finish handling existing requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown flow engine
	flowEngine.Shutdown(ctx)

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
