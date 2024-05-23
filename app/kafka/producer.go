package kafka

import (
	"github.com/IBM/sarama"
	"log"
)

func ExampleProducer(producer sarama.AsyncProducer) {
	message := "Hello Material Service !"
	producer.Input() <- &sarama.ProducerMessage{
		Topic: "example",
		Value: sarama.StringEncoder(message),
	}
	select {
	case success := <-producer.Successes():
		log.Printf("Message produced: offset=%d, timestamp=%v, partitions=%d\n", success.Offset, success.Timestamp, success.Partition)
	case err := <-producer.Errors():
		log.Printf("Failed to produce message: %v\n", err)
	}
}
