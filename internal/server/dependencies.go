package server

import (
	"context"
	"log"
	"os"
	"rag-poc/internal/database"
	"strconv"
)

func getServerPort() int {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		log.Fatalf("Invalid or missing PORT: %v", err)
	}
	return port
}

func initializeDatabase() database.Service {
	db := database.New()
	ctx := context.Background()
	if err := db.DB().Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}
