package main

import (
	"fmt"
	"os"
	"path/filepath"

	"dottie-modus/src/db" // Adjust the import path based on your module name
)

func main() {
	filePath := os.Getenv("ACOG_DATA_FILE")
	if filePath == "" {
		filePath = filepath.Join("data", "acog_guidelines.json")
	}
	if err := db.LoadACOGData(filePath); err != nil {
		fmt.Printf("Error populating database: %v\n", err)
	} else {
		fmt.Println("Database population completed!")
	}
}
