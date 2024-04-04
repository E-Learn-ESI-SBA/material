package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBHandler(uri string) (*mongo.Client, error) {
	println("Connecting... to the database")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	opts.SetConnectTimeout(10 * time.Second)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil, err

	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil, err
	}
	log.Println("Successfully Connected to the mongodb")
	return client, nil
}

func ModuleCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	return client.Database("materials").Collection(CollectionName)

}

func CourseCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)
	return collection
}

func SectionCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)
	return collection
}
func LectureCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)
	return collection
}
func VideoCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)
	return collection
}
func ContentCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)
	return collection
}
func CommentCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	collection := client.Database("materials").Collection(CollectionName)
	return collection
}

type Application struct {
	ContentCollection  *mongo.Collection
	VideoCollection    *mongo.Collection
	LectureCollection  *mongo.Collection
	SectionCollection  *mongo.Collection
	CourseCollection   *mongo.Collection
	ModuleCollection   *mongo.Collection
	CommentsCollection *mongo.Collection
}

func NewApp(client *mongo.Client) *Application {
	return &Application{
		VideoCollection:    VideoCollection(client, "videos"),
		LectureCollection:  LectureCollection(client, "lectures"),
		ContentCollection:  ContentCollection(client, "contents"),
		SectionCollection:  SectionCollection(client, "sections"),
		CourseCollection:   CourseCollection(client, "courses"),
		ModuleCollection:   ModuleCollection(client, "modules"),
		CommentsCollection: CommentCollection(client, "comments"),
	}

}

/*

indexModel := mongo.IndexModel{
		Keys:    bson.D{{"_id", 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)

	}
*/
