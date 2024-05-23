package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"madaurus/dev/material/app/interfaces"
	"time"
)

type KafkaInstance struct {
	Consumer sarama.Consumer
	Producer sarama.AsyncProducer
}

func KafkaInit(kafkaSettings interfaces.Kafka) *KafkaInstance {
	var err error
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version, err = sarama.ParseKafkaVersion(sarama.DefaultVersion.String())
	if err != nil {
		log.Fatal("Unable to parse Kafka version: ", err)
	}

	// Consumer Config
	kafkaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	kafkaConfig.Consumer.Return.Errors = true
	kafkaConfig.Consumer.Offsets.AutoCommit.Enable = true
	kafkaConfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
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
	return &KafkaInstance{
		Consumer: consumer,
		Producer: producer,
	}
}

type ConsumerFn func(ctx context.Context, message *sarama.ConsumerMessage) error

type HandlerMapConsumerGroupHandler struct {
	handlers map[string]ConsumerFn
}

func NewConsumerGroupHandler(handlers map[string]ConsumerFn) *HandlerMapConsumerGroupHandler {
	return &HandlerMapConsumerGroupHandler{
		handlers: handlers,
	}
}
func (h HandlerMapConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		handler, ok := h.handlers[msg.Topic]
		if !ok {
			return fmt.Errorf("missing handler for topic: %s", msg.Topic)
		}
		if err := handler(sess.Context(), msg); err != nil {
			return err
		}
		sess.MarkMessage(msg, "")
	}
	return nil
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
	return
}
