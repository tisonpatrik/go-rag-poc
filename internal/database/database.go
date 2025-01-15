package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	Health() string
	Stats() map[string]string
	Close() error
	DB() *pgxpool.Pool
}

type service struct {
	db *pgxpool.Pool
}

func (s *service) DB() *pgxpool.Pool {
	return s.db
}

var (
	database   = os.Getenv("RAG_DB_DATABASE")
	password   = os.Getenv("RAG_DB_PASSWORD")
	username   = os.Getenv("RAG_DB_USERNAME")
	port       = os.Getenv("RAG_DB_PORT")
	host       = os.Getenv("RAG_DB_HOST")
	dbInstance *service
)

// New initializes the database connection service.
func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Failed to parse database configuration: %v", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	dbInstance = &service{
		db: db,
	}

	log.Printf("Connected to database: %s", database)
	return dbInstance
}

// Health checks the connectivity of the database.
func (s *service) Health() string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.Ping(ctx); err != nil {
		log.Printf("Database health check failed: %v", err)
		return "down"
	}
	return "up"
}

// Stats provides key database statistics.
func (s *service) Stats() map[string]string {
	poolStats := s.db.Stat()

	stats := map[string]string{
		"idle_connections":     fmt.Sprintf("%d", poolStats.IdleConns()),
		"total_connections":    fmt.Sprintf("%d", poolStats.TotalConns()),
		"acquired_connections": fmt.Sprintf("%d", poolStats.AcquiredConns()),
	}

	if poolStats.AcquiredConns() > 40 {
		stats["load_message"] = "The database is under heavy load."
	} else {
		stats["load_message"] = "Database load is normal."
	}

	return stats
}

// Close closes the database connection pool.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	s.db.Close()
	return nil
}
