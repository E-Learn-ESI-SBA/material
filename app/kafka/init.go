package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/utils"
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

func (kafkaInstance *KafkaInstance) ProduceMessage(topic string, message string) error {
	kafkaInstance.Producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	select {
	case success := <-kafkaInstance.Producer.Successes():
		fmt.Printf("Message produced: offset=%d, timestamp=%v, partitions=%d\n", success.Offset, success.Timestamp, success.Partition)
	case err := <-kafkaInstance.Producer.Errors():
		fmt.Printf("Failed to produce message: %v\n", err)
		return err
	}
	return nil
}

func (kafkaInstance *KafkaInstance) EvaluationProducer(user *utils.UserDetails, resourceType string, evaluationPoint int32) error {

	evaluation := interfaces.EvaluationConsumer{
		UserId:          user.ID,
		EvaluationPoint: evaluationPoint,
		Date:            time.Now().Format("2006-01-02"),
	}

	// Marshal the evaluation object as string
	evaluationBytes, err := json.Marshal(evaluation)
	if err != nil {
		log.Printf("Error while marshalling the evaluation object: %v", err)
		return err
	}
	evaluationString := string(evaluationBytes)
	log.Println("Evaluation String: ", evaluationString)
	err = kafkaInstance.ProduceMessage(EVALUATION, evaluationString)

	if err != nil {
		return err
	}
	notification := NotificationEvent{
		Group: user.Group,
		Year:  user.Year,
		// Say congratulation  to username for completing the video and earning evaluation point
		Message:     fmt.Sprintf("Congratulations %s for completing the %s and earning %d evaluation points", user.Username, resourceType, evaluationPoint),
		UserId:      user.ID,
		Role:        user.Role,
		EnablePush:  true,
		PushTo:      USER_NOTIFICATION_TYPE,
		Title:       "Video Completion",
		EnableEmail: false,
	}
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error while marshalling the notification object: %v", err)
		return err
	}
	notificationString := string(notificationBytes)
	log.Println("Notification String: ", notificationString)
	err = kafkaInstance.ProduceMessage(NOTIFICATION_TOPIC, notificationString)
	return err
}

func (kafkaInstance *KafkaInstance) ResourceCreatingProducer(user *utils.UserDetails, resourceType string, resourceName string, notificationType string) error {
	notification := NotificationEvent{
		Group: user.Group,
		Year:  user.Year,
		// Say congratulation  to username for creating the resource
		Message:     fmt.Sprintf("New %s Created, Named  %s", resourceType, resourceName),
		UserId:      user.ID,
		Role:        user.Role,
		EnablePush:  true,
		PushTo:      notificationType,
		Title:       fmt.Sprintf("New %s Created", resourceType),
		EnableEmail: false,
	}
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error while marshalling the notification object: %v", err)
		return err
	}
	notificationString := string(notificationBytes)
	log.Println("Notification String: ", notificationString)
	err = kafkaInstance.ProduceMessage(NOTIFICATION_TOPIC, notificationString)
	return err
}
