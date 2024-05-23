package interfaces

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ModuleCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "courses._id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "courses.sections._id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "courses.sections.files._id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "courses.sections.videos._id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "courses.sections.lectures._id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexModels)
	if err != nil {
		log.Printf("Failed to create indexes: %v", err)
	}

	return collection
}

/*
	func ContentCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
		collection := client.Database("materials").Collection(CollectionName)
		collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys: bson.D{{"name", 1}},
		})
		return collection
	}
*/
func CommentCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)
	return collection
}
func UserCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)
	return collection
}

func QuizesCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)
	return collection
}

func QuestionCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)
	return collection
}

func GradesCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)
	return collection
}

type Application struct {
	ModuleCollection      *mongo.Collection
	CommentsCollection    *mongo.Collection
	UserCollection        *mongo.Collection
	QuizesCollection      *mongo.Collection
	SubmissionsCollection *mongo.Collection
}

func NewApp(client *mongo.Client) *Application {
	return &Application{
		ModuleCollection:      ModuleCollection(client, "modules"),
		CommentsCollection:    CommentCollection(client, "comments"),
		UserCollection:        UserCollection(client, "users"),
		QuizesCollection:      QuizesCollection(client, "quizes"),
		SubmissionsCollection: CommentCollection(client, "submissions"),
	}

}
