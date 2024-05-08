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
	Keys []string
}

func (r *ResourceKeys) GetResourceKey(key string) {
	r.Keys = append(r.Keys, key)
}
func GetAllowedResources(actionName string, resourceType string, userKey string, permitApi *permit.Client) []string {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	requestContext := map[string]string{
		"source": "docs",
	}
	user := enforcement.UserBuilder(userKey).Build()
	action := enforcement.Action(actionName)
	resources, errR := permitApi.Api.ResourceInstances.ListDetailed(ctx, 1, 100, iam.TENANT, resourceType, "")
	if errR != nil {
		log.Printf("Error While Checking the Permission: %v\n", errR)
		return []string{}
	}
	var singleResourceBuilder enforcement.Resource
	var resourcesBuilders []enforcement.Resource
	for _, resource := range *resources {
		singleResourceBuilder = enforcement.ResourceBuilder(resourceType).WithKey(resource.Key).WithTenant(enforcement.DefaultTenant).WithContext(requestContext).Build()
		resourcesBuilders = append(resourcesBuilders, singleResourceBuilder)

	}
	allowedResources, err := GetFilterObject(user, action, permitApi, resourcesBuilders...)
	if err != nil {
		log.Printf("Error While Checking the Permission: %v\n", err)
		return []string{}
	}
	if len(allowedResources) == 0 {
		panic("No roles ")
	}
	filterR := ResourceKeys{}
	for _, rs := range allowedResources {
		log.Println("Role: ", rs.Key)
		log.Println("Role Type: ", rs.GetType())
		filterR.GetResourceKey(rs.Key)
	}
	return filterR.Keys

}

func GetFilterObject(user enforcement.User, action enforcement.Action, permitApi *permit.Client, resources ...enforcement.Resource) ([]enforcement.Resource, error) {
	allowedResources := make([]enforcement.Resource, 0)
	for _, resource := range resources {
		decision, err := permitApi.Check(user, action, resource)
		if err != nil {
			return nil, err
		}
		if decision {
			allowedResources = append(allowedResources, resource)
		}
	}
	return allowedResources, nil

}
