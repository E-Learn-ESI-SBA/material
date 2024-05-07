package utils

import (
	"context"
	"github.com/permitio/permit-golang/pkg/enforcement"
	"github.com/permitio/permit-golang/pkg/permit"
	"log"
	"madaurus/dev/material/app/shared/iam"
	"mime/multipart"
	"strings"
	"time"
)

func GetFileTypeFromMIME(file *multipart.FileHeader) (string, error) {
	fileName := file.Filename
	// get the file extension if it exists ( pdf , excel , docs ,go, etc... )
	// if the file extension does not exist, return an empty string
	fileExtension := fileName[strings.LastIndex(fileName, "."):]
	return fileExtension, nil

}

type ResourceKeys struct {
	RType string
	keys  []string
}

func (r *ResourceKeys) GetResourceKey(key string, rtype string) {
	if r.RType == rtype {
		r.keys = append(r.keys, key)
	}
}
func GetAllowedResources(actionName string, resourceType string, userKey string, permitApi *permit.Client) []string {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	requestContext := map[string]string{
		"source": resourceType,
	}
	user := enforcement.UserBuilder(userKey).Build()
	action := enforcement.Action(actionName)
	resources, errR := permitApi.Api.ResourceInstances.ListDetailed(ctx, 1, 100, iam.TENANT, resourceType, "")
	if errR != nil {
		// handle error
		log.Printf("Error While Checking the Permission: %v\n", errR)
		return []string{}
	}
	var singleResourceBuidler enforcement.ResourceI
	var resourcesBuilders []enforcement.ResourceI
	for _, resource := range *resources {
		log.Println("Resource Id : ", resource.Key)
		singleResourceBuidler = enforcement.ResourceBuilder(resourceType).WithKey(resource.Key)
		resourcesBuilders = append(resourcesBuilders, singleResourceBuidler)

	}
	log.Printf("Length of Resources: %v\n", len(resourcesBuilders))
	rolesS, err := permitApi.FilterObjects(user, action, requestContext, resourcesBuilders...)
	if err != nil {
		// handle error
		log.Printf("Error While Checking the Permission: %v\n", err)
		return []string{}
	}
	filterR := ResourceKeys{RType: resourceType}

	for _, role := range rolesS {
		log.Println("Role: ", role.GetID())
		log.Println("Role Type: ", role.GetType())
		filterR.GetResourceKey(role.GetID(), role.GetType())
	}
	return filterR.keys
	// i want get all resources alowred by user

}
