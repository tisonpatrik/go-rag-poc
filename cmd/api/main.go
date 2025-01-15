package main

import (
	"log"
	"rag-poc/internal/server"
)

func main() {
	// Initialize the server
	srv := server.NewServer()

	// Run the server
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
