package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"madaurus/dev/material/app/interfaces"
)

type KafkaInstance struct {
	Consumer sarama.Consumer
	Producer sarama.AsyncProducer
}

func KafkaInit(kafkaSettings interfaces.Kafka) KafkaInstance {
	var err error
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version, err = sarama.ParseKafkaVersion(sarama.DefaultVersion.String())
	if err != nil {
		log.Fatal("Unable to parse Kafka version: ", err)
	}

	// Consumer Config
	kafkaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	kafkaConfig.Consumer.Return.Errors = true
	kafkaConfig.Consumer.Offsets.AutoCommit.Enable = true
	kafkaConfig.Consumer.Offsets.AutoCommit.Interval = 1
	kafkaConfig.Consumer.Offsets.CommitInterval = 1
	kafkaConfig.Consumer.Offsets.Retry.Max = 3

	// Producer Config
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 3
	kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner

	// Create new consumer
	consumer, err := sarama.NewConsumer([]string{kafkaSettings.Host}, kafkaConfig)
	if err != nil {
		log.Fatalf("Error while creating consumer: %v", err)
	}

	// Create new producer
	producer, err := sarama.NewAsyncProducer([]string{kafkaSettings.Host}, kafkaConfig)
	if err != nil {
		log.Fatalf("Error while creating producer: %v", err)
	}

	return KafkaInstance{
		Consumer: consumer,
		Producer: producer,
	}
}

func ProduceMessage(producer sarama.AsyncProducer, topic string, message string) {
	producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	select {
	case success := <-producer.Successes():
		fmt.Printf("Message produced: offset=%d, timestamp=%v, partitions=%d\n", success.Offset, success.Timestamp, success.Partition)
	case err := <-producer.Errors():
		fmt.Printf("Failed to produce message: %v\n", err)
	}
}
