package transactions

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GetStorageFile(dir string) (string, error) {

	// Get the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return "", err
	}

	// Go up two directories to get to the project root
	log.Printf("Current working directory: %v", pwd)

	// Construct the path to the storage/files directory
	storageFilesPath := filepath.Join(pwd, "storage", dir)

	log.Printf("Path to storage/files: %v", storageFilesPath)
	return storageFilesPath, nil
}
