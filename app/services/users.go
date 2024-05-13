package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/utils"
)

func EditUser(ctx context.Context, user utils.LightUser, collection *mongo.Collection) error {
	filter := bson.D{{"_id", user.ID}}
	rs := collection.FindOneAndUpdate(ctx, filter, user)
	err := rs.Err()
	return err
}
func CreateUser(ctx context.Context, user utils.LightUser, collection *mongo.Collection) error {
	_, err := collection.InsertOne(ctx, user)
	return err
}

func DeleteUser(ctx context.Context, user utils.LightUser, collection *mongo.Collection) error {
	filter := bson.D{{"_id", user.ID}}
	rs := collection.FindOneAndDelete(ctx, filter)
	return rs.Err()
}

func GetUserById(ctx context.Context, userId string, collection *mongo.Collection) (models.User, error) {
	filter := bson.D{{"userId", userId}}
	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Printf("Error While Getting the User: %v\n", err.Error())
		return models.User{}, err
	}
	return user, nil
}
