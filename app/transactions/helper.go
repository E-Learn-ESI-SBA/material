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
	projectRoot := filepath.Dir(filepath.Dir(pwd))

	// Construct the path to the storage/files directory
	storageFilesPath := filepath.Join(projectRoot, "storage", dir)

	log.Printf("Path to storage/files: %v", storageFilesPath)
	return storageFilesPath, nil
}
