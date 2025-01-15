package rag

import (
	"net/http"
	"rag-poc/internal/openaiclient"
)

func RegisterRagRoutes(router *http.ServeMux, client *openaiclient.Service) {
	handler := NewHandler(*client)

	router.HandleFunc("/rag/react_agent", handler.reActHandler)
}
