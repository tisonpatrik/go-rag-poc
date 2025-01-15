package rag

import (
	"net/http"
	"rag-poc/internal/openaiclient"
	"rag-poc/internal/rag/handler"
)

func RegisterRagRoutes(router *http.ServeMux, client *openaiclient.Service) {
	handler := handler.NewHandler(*client)

	router.HandleFunc("/rag/react_agent", handler.ReActHandler)
}
