package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetUser(ctx context.Context, userId primitive.ObjectID, collection *mongo.Collection) (utils.LightUser, error) {
	filter := bson.D{{"_id", userId}}
	var user utils.LightUser
	err := collection.FindOne(ctx, filter).Decode(&user)
	return user, err
}
