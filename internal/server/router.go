package server

import (
	"net/http"
	"rag-poc/internal/database"
	"rag-poc/internal/ingest"
	"rag-poc/internal/openaiclient"
	"rag-poc/internal/rag"
	"rag-poc/internal/repository"
)

func registerRoutes(router *http.ServeMux, db database.Service, client *openaiclient.Service, queries *repository.Queries) {
	registerServerRoutes(router, db)
	rag.RegisterRagRoutes(router, client)
	ingest.RegisterRagRoutes(router, queries)
}

func registerServerRoutes(router *http.ServeMux, db database.Service) {
	serverHandler := &ServerHandler{db: db}

	router.HandleFunc("GET /hello", serverHandler.helloWorldHandler)
	router.HandleFunc("GET /health_db", serverHandler.healthDBHandler)
	router.HandleFunc("GET /stats_db", serverHandler.statsDBHandler)
	router.HandleFunc("GET /health", serverHandler.healthHandler)
}
