package utils

import (
	"mime/multipart"
	"strings"
)

func GetFileTypeFromMIME(file *multipart.FileHeader) (string, error) {
	fileName := file.Filename
	// get the file extension if it exists ( pdf , excel , docs ,go, etc... )
	// if the file extension does not exist, return an empty string
	fileExtension := fileName[strings.LastIndex(fileName, "."):]
	return fileExtension, nil

}
