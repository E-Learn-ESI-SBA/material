package utils

import (
	"context"
	"errors"
	"github.com/permitio/permit-golang/pkg/enforcement"
	"github.com/permitio/permit-golang/pkg/models"
	"github.com/permitio/permit-golang/pkg/permit"
	"madaurus/dev/material/app/shared"
	"time"
)

func CreateResourceInstance(permitApi *permit.Client, resourceType string, instanceKey string, leaderKey *string, leaderResourceType *string, relation *string) error {
	tenant := enforcement.DefaultTenant
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
