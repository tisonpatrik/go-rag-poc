package server

import (
	"net/http"
	"rag-poc/internal/builder"
	"rag-poc/internal/database"
	"rag-poc/internal/repository"
)

func registerRoutes(router *http.ServeMux, db database.Service, queries *repository.Queries) {
	registerServerRoutes(router, db)
	builder.RegisterBuilderRoutes(router, queries)
}

func registerServerRoutes(router *http.ServeMux, db database.Service) {
	serverHandler := &ServerHandler{db: db}

	router.HandleFunc("/hello", serverHandler.helloWorldHandler)
	router.HandleFunc("/health_db", serverHandler.healthDBHandler)
	router.HandleFunc("/stats_db", serverHandler.statsDBHandler)
	router.HandleFunc("/health", serverHandler.healthHandler)
}
