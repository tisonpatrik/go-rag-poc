package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"rag-poc/internal/database"
	"rag-poc/internal/server/middleware"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"rag-poc/internal/repository"
)

type Server struct {
	port    int
	db      database.Service
	queries repository.Queries
	httpSrv *http.Server
}

// NewServer initializes a new server instance.
func NewServer() *Server {
	port := getServerPort()
	db := initializeDatabase()

	// Initialize queries with the database connection
	queries := repository.New(db.DB())

	router := http.NewServeMux()
	registerRoutes(router, queries)

	server := &Server{
		port:    port,
		db:      db,
		queries: *queries,
	}

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.CORS,
	)

	// Initialize the HTTP server
	server.httpSrv = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      stack(router),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

// Run starts the HTTP server and handles graceful shutdown.
func (s *Server) Run() error {
	// Channel to signal shutdown completion
	done := make(chan bool, 1)

	// Start graceful shutdown in a separate goroutine
	go s.gracefulShutdown(done)

	// Start the server
	log.Printf("Starting server on port %d", s.port)
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	// Wait for the shutdown to complete
	<-done
	log.Println("Server shut down gracefully.")
	return nil
}

// gracefulShutdown handles server shutdown gracefully.
func (s *Server) gracefulShutdown(done chan bool) {
	// Listen for termination signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("Shutting down server gracefully...")

	// Shutdown with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		log.Printf("Forced shutdown error: %v", err)
	}

	// Close database connection
	if err := s.db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	done <- true
}
