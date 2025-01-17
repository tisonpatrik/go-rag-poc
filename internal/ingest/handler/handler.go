package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rag-poc/internal/repository"
	"rag-poc/internal/utils"
	"time"
)

type IngestHandler struct {
	Queries repository.Queries
}

func NewHandler(queries repository.Queries) *IngestHandler {
	return &IngestHandler{
		Queries: queries,
	}
}

func (h *IngestHandler) HandleIngestionOfDocument(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form, allowing files up to 10 MB
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid form data"})
		return
	}

	// Process the file using the existing uploadFile logic
	ctx := context.Background()
	response, err := h.uploadFile(ctx, r)
	if err != nil {
		log.Printf("Error processing file: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}

	// Respond with success
	utils.SendJSONResponse(w, http.StatusCreated, map[string]string{"response": response})
}

func (h *IngestHandler) uploadFile(ctx context.Context, r *http.Request) (string, error) {
	fmt.Println("Processing uploaded file")

	// Retrieve the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving the file: %v", err)
		return "", fmt.Errorf("failed to retrieve file: %w", err)
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Validate that the file is a .txt file
	if err := utils.ValidateTxtFile(handler.Filename); err != nil {
		log.Printf("Invalid file type: %v", err)
		return "", fmt.Errorf("invalid file type: %w", err)
	}

	// Read the file content
	content, err := utils.ReadFile(file)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Sanitize the file content
	sanitizedContent, err := utils.SanitizeTxtContent(content)
	if err != nil {
		log.Printf("Error sanitizing content: %v", err)
		return "", fmt.Errorf("failed to sanitize file content: %w", err)
	}

	// Prepare the InsertDocumentParams
	insertParams := repository.InsertDocumentParams{
		DocumentName: handler.Filename,
		DateTime:     time.Now(),
		OriginalLink: "",
		Content:      sanitizedContent,
	}

	// Insert into the database
	insertedDocument, err := h.Queries.InsertDocument(ctx, insertParams)
	if err != nil {
		log.Printf("Database error: %v", err)
		return "", fmt.Errorf("failed to insert document into the database: %w", err)
	}

	// Return success message
	successMessage := fmt.Sprintf("Successfully uploaded, sanitized, and saved file to database. Inserted Document ID: %s", insertedDocument.ID)
	return successMessage, nil
}
