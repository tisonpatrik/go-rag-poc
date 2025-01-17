package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rag-poc/internal/database"
	"rag-poc/internal/openaiclient"
	"rag-poc/internal/repository"
	"rag-poc/internal/server/middleware"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port         string
	db           database.Service
	queries      repository.Queries
	openAIClient openaiclient.Service
	httpSrv      *http.Server
}

// NewServer initializes a new server instance.
func NewServer() *Server {
	port := os.Getenv("PORT")
	db := database.New()
	openAI := openaiclient.New()

	// Initialize queries with the database connection
	queries := repository.New(db.DB())

	router := http.NewServeMux()
	registerRoutes(router, db, &openAI, queries)

	server := &Server{
		port:         port,
		db:           db,
		queries:      *queries,
		openAIClient: openAI,
	}

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.CORS,
	)

	// Initialize the HTTP server
	server.httpSrv = &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      stack(router),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

// Run starts the HTTP server and handles graceful shutdown.
func (s *Server) Run() error {
	done := make(chan bool, 1)

	go s.gracefulShutdown(done)

	log.Printf("Starting server on port %s", s.port)
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	<-done
	log.Println("Server shut down gracefully.")
	return nil
}

// gracefulShutdown handles server shutdown gracefully.
func (s *Server) gracefulShutdown(done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		log.Printf("Forced shutdown error: %v", err)
	}

	if err := s.db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	done <- true
}
