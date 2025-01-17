package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"unicode/utf8"
)

// parseJSON decodes JSON input into the specified structure
func ParseJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

// sendJSONResponse sends a JSON response to the client
func SendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// ReadFile reads the content of an uploaded file line by line.
func ReadFile(file multipart.File) (string, error) {
	var content strings.Builder
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", errors.New("error reading file content")
		}
		content.WriteString(line)
		if err == io.EOF {
			break
		}
	}
	return content.String(), nil
}

// ValidateTxtFile checks if the uploaded file is a text file based on its extension.
func ValidateTxtFile(filename string) error {
	if !strings.HasSuffix(filename, ".txt") {
		return errors.New("only .txt files are allowed")
	}
	return nil
}

// SanitizeTxtContent ensures that the file content is UTF-8 compliant.
func SanitizeTxtContent(content string) (string, error) {
	if !utf8.ValidString(content) {
		return "", errors.New("file contains invalid UTF-8 characters")
	}
	return content, nil
}

// GetFileName extracts the file name from a given path
func GetFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
