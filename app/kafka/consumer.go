package kafka

import (
	"github.com/IBM/sarama"
	"log"
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
