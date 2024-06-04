package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"madaurus/dev/material/app/shared"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func GetFileTypeFromMIME(file *multipart.FileHeader) (string, error) {
	fileName := file.Filename
	// get the file extension if it exists ( pdf , excel , docs ,go, etc... )
	// if the file extension does not exist, return an empty string
	fileExtension := fileName[strings.LastIndex(fileName, "."):]
	return fileExtension, nil

}

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

func GetUserPayload(c *gin.Context) (*UserDetails, error) {
	var user *UserDetails
	value, err := c.Get("user")
	if !err {
		log.Printf("Error getting user from context: ")
		return user, errors.New(shared.USER_NOT_INJECTED)
	}
	user = value.(*UserDetails)
	return user, nil
}
