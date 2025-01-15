package server

import (
	"net/http"
	"rag-poc/internal/builder"
	"rag-poc/internal/repository"
)

func registerRoutes(router *http.ServeMux, queries *repository.Queries) {
	registerServerRoutes(router)
	builder.RegisterBuilderRoutes(router, queries)
}

func registerServerRoutes(router *http.ServeMux) {

	serverHandler := &ServerHandler{db: initializeDatabase()}

	router.HandleFunc("/hello", serverHandler.helloWorldHandler)
	router.HandleFunc("/health_db", serverHandler.healthDBHandler)
	router.HandleFunc("/stats_db", serverHandler.statsDBHandler)
	router.HandleFunc("/health", serverHandler.healthHandler)
}
