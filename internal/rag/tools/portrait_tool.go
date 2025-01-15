package tools

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// GetAutoportrait reads the self_portrait.txt file and formats the output.
func GetAutoportrait() string {
	// Define the path to the self_portrait.txt file
	filePath := filepath.Join("internal", "rag", "tools", "self_portrait.txt")

	// Read the contents of the file
	fileContent, err := os.ReadFile(filePath) // Updated to os.ReadFile
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
	}
	// Format the output
	return fmt.Sprintf("%s", string(fileContent))
}
