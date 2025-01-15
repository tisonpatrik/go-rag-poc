package builder

import (
	"net/http"
	"rag-poc/internal/repository"
)

func RegisterBuilderRoutes(router *http.ServeMux, queries *repository.Queries) {
	handler := NewHandler(queries, 100) // Queue size of 100

	// Start the worker pool
	go handler.ProcessQueue(5) // 5 workers

	// Register routes
	router.HandleFunc("/builder/listen", handler.ListenHandler)
}
