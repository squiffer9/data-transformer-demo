package main

import (
	"context"
	"data-transformer-demo/internal/cache"
	"data-transformer-demo/internal/db"
	"data-transformer-demo/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func transformHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Start timing
	start := time.Now()

	var input service.TransformRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Basic validation
	if input.Country == "" {
		http.Error(w, "Country is required", http.StatusBadRequest)
		return
	}

	if len(input.Data) == 0 {
		http.Error(w, "Data array is required", http.StatusBadRequest)
		return
	}

	if len(input.Data) > 5000 {
		http.Error(w, "Request exceeds maximum allowed entries (5000)", http.StatusBadRequest)
		return
	}

	// Process transformation
	result := service.Transform(input)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Log timing
	duration := time.Since(start)
	log.Printf("Request processed in %v - Country: %s, Input pairs: %d, Output pairs: %d",
		duration, input.Country, len(input.Data), len(result.Data))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// Initialize and load cache
	cache := cache.GetInstance()
	if err := cache.LoadData(); err != nil {
		log.Fatalf("Failed to load initial cache data: %v", err)
	}

	// Start cache refresh loop (every 5 minutes)
	cache.StartRefreshLoop(5 * time.Minute)

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Register handlers
	http.HandleFunc("/transform", transformHandler)
	http.HandleFunc("/health", healthHandler)

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port 8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Graceful shutdown
	log.Printf("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}
}
