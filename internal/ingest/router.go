package ingest

import (
	"net/http"
	"rag-poc/internal/ingest/handler"
	"rag-poc/internal/repository"
)

func RegisterRagRoutes(router *http.ServeMux, queries *repository.Queries) {
	handler := handler.NewHandler(*queries)

	router.HandleFunc("POST /upload", handler.HandleIngestionOfDocument)
}
