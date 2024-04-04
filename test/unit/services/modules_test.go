package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"testing"
	"time"
)

func TestCreateModule(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("Success", func(mt *mtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		module := models.Module{Name: "OOP", TeacherId: 125, Year: 2, Coefficient: 4, IsPublic: false, Semester: 2}
		err := services.CreateModule(ctx, mt.Coll, module)
		assert.Nil(t, err)
		defer cancel()

	})
	mt.Run("Error", func(mt *mtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{Code: 11000}))
		module := models.Module{Name: "OOP", TeacherId: 125, Year: 2, Coefficient: 4, IsPublic: false, Semester: 2}
		err := services.CreateModule(ctx, mt.Coll, module)
		assert.NotNil(t, err)
		defer cancel()
	})
}
func TestUpdateModule(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("Success", func(mt *mtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		module := models.Module{Name: "OOP", TeacherId: 125, Year: 2, Coefficient: 4, IsPublic: false, Semester: 2}
		err := services.UpdateModule(ctx, mt.Coll, module)
		assert.Nil(t, err)
		defer cancel()

	})
	mt.Run("Error", func(mt *mtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{Code: 11000}))
		module := models.Module{Name: "OOP", TeacherId: 125, Year: 2, Coefficient: 4, IsPublic: false, Semester: 2}
		_, err := mt.Coll.InsertOne(ctx, module)
		assert.NotNil(t, err)
		defer cancel()
	})
}
