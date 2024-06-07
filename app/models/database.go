package models

import (
	"context"
	"log"
	"madaurus/dev/material/app/interfaces"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Setup(database *interfaces.Database) (*mongo.Client, *interfaces.Application, error) {
	log.Println("Connecting... to the database with : %v", database.Host)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(database.Host).SetServerAPIOptions(serverAPI)
	opts.SetConnectTimeout(10 * time.Second)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil, nil, err

	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil, nil, err
	}
	log.Println("Successfully Connected to the mongodb")
	app := interfaces.NewApp(client)
	return client, app, nil

}

func Close(client *mongo.Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to MongoDB closed.")
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
