package utils

import (
	"context"
	"errors"
	"github.com/permitio/permit-golang/pkg/enforcement"
	"github.com/permitio/permit-golang/pkg/models"
	"github.com/permitio/permit-golang/pkg/permit"
	"log"
	"madaurus/dev/material/app/logs"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/shared/iam"
	"time"
)

func CreateResourceInstance(permitApi *permit.Client, resourceType string, instanceKey string, leaderKey *string, leaderResourceType *string, relation *string) error {
	tenant := enforcement.DefaultTenant
	logs.Info("Creating the resource instance")
	resource := models.ResourceInstanceCreate{
		Resource: resourceType,
		Key:      instanceKey,
		Tenant:   &tenant,
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	resourceRead, err := permitApi.Api.ResourceInstances.Create(ctx, resource)
	if err != nil {
		return errors.New(shared.UNABLE_CREATE_INSTANCE)
	}
	// Create the relation
	if leaderKey != nil {
		relationTuple := models.RelationshipTupleCreate{
			Tenant:   &tenant,
			Relation: *relation,
			Object:   resourceType + ":" + resourceRead.Key,
			Subject:  *leaderResourceType + ":" + *leaderKey,
		}
		_, errT := permitApi.Api.RelationshipTuples.Create(ctx, relationTuple)
		if errT != nil {
			return errors.New(shared.UNABLE_CREATE_INSTANCE)
		}
	}
	return nil
}

type ResourceKeys struct {
	Keys []string
}

func (r *ResourceKeys) GetResourceKey(key string) {
	r.Keys = append(r.Keys, key)
}
func GetAllowedResources(actionName string, resourceType string, userKey string, permitApi *permit.Client) []string {
	logs.Info("Getting the allowed resources")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	requestContext := map[string]string{
		"source": "docs",
	}
	user := enforcement.UserBuilder(userKey).Build()
	action := enforcement.Action(actionName)
	// For loop to get all the resources
	var resourcesArr []models.ResourceInstanceRead
	// for loop to get to append the resources
	for i := 1; i < 10; i++ {
		resources, err := permitApi.Api.ResourceInstances.ListDetailed(ctx, i, 100, iam.TENANT, resourceType, "")
		if err != nil {
			log.Printf("Error While Checking the Permission: %v\n", err)
			break
		}
		// append the resources
		// Add resources array to resourcesArr
		resourcesArr = append(resourcesArr, *resources...)
	}

	var singleResourceBuilder enforcement.Resource
	var resourcesBuilders []enforcement.Resource
	for _, resource := range resourcesArr {
		singleResourceBuilder = enforcement.ResourceBuilder(resourceType).WithKey(resource.Key).WithTenant(enforcement.DefaultTenant).WithContext(requestContext).Build()
		resourcesBuilders = append(resourcesBuilders, singleResourceBuilder)

	}

	allowedResources, err := GetFilterObject(user, action, permitApi, requestContext, resourcesBuilders...)
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

func GetFilterObject(user enforcement.User, action enforcement.Action, permitApi *permit.Client, ctx map[string]string, resources ...enforcement.Resource) ([]enforcement.Resource, error) {
	allowedResources := make([]enforcement.Resource, 0)
	var checkRes []enforcement.CheckRequest
	for _, resource := range resources {
		checkR := enforcement.NewCheckRequest(user, action, resource, ctx)
		checkRes = append(checkRes, *checkR)
	}
	checks, err := permitApi.BulkCheck(checkRes...)
	if err != nil {
		return allowedResources, err
	}
	for i, check := range checks {
		if check {
			allowedResources = append(allowedResources, resources[i])

		}
	}
	return allowedResources, nil

}
