package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/utils"
)

func ExampleConsumer(consumer sarama.Consumer) {
	consumerPartition, err := consumer.ConsumePartition("example", 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("Error While consuming ")

	}
	defer consumerPartition.Close()

	for msg := range consumerPartition.Messages() {
		log.Printf("Received message: %s\n", string(msg.Value))
	}
}

func UserMutationHandler(consumer sarama.Consumer, collection *mongo.Collection) {
	var user utils.LightUser
	consumerPartition, err := consumer.ConsumePartition(USER_MUTATION, 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("Error While consuming ")

	}
	ctx := context.Background()
	defer consumerPartition.Close()

	for msg := range consumerPartition.Messages() {
		// Parse the message as a user from json
		err = json.Unmarshal(msg.Value, &user)
		if err != nil {
			log.Printf("Error While Unmarshalling the User %v", err)
			log.Printf("Error While Unmarshalling the User: %v\n", err)
			continue
		}
		err = services.EditUser(ctx, user, collection)
		if err != nil {
			log.Printf("\t Event : Unable to edit the user from the ")
			log.Println("\t Event : Unable to edit the user from the ")
			return
		}
	}

}
