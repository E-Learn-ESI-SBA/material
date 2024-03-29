package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBHandler(uri string) *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	opts.SetConnectTimeout(10 * time.Second)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil
	}
	fmt.Println("Successfully Connected to the mongodb")
	return client
}

func ModuleCollection(client *mongo.Client, CollectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("materials").Collection(CollectionName)
	return collection

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

type Application struct {
	ContentCollection *mongo.Collection
	VideoCollection   *mongo.Collection
	LectureCollection *mongo.Collection
	SectionCollection *mongo.Collection
	CourseCollection  *mongo.Collection
	ModuleCollection  *mongo.Collection
}

func (app *Application) CreateApp(client *mongo.Client) {
	app = &Application{
		VideoCollection:   VideoCollection(client, "videos"),
		LectureCollection: LectureCollection(client, "lectures"),
		ContentCollection: ContentCollection(client, "contents"),
		SectionCollection: SectionCollection(client, "sections"),
		CourseCollection:  CourseCollection(client, "courses"),
		ModuleCollection:  ModuleCollection(client, "modules"),
	}

}
