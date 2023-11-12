package kafka

import (
	"job-scheduler-worker/config"
	"job-scheduler-worker/internal/entity"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type messageBusRepository struct {
	producer *kafka.Producer
}

func NewMessageBusRepository(conf *config.Config) (entity.MessageBusRepository, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": conf.Kafka.BootstrapServers,
		"client.id":         conf.Kafka.ClientId,
		"acks":              conf.Kafka.ACKS,
	})
	if err != nil {
		return nil, err
	}

	return &messageBusRepository{
		producer: producer,
	}, nil
}

func (_self *messageBusRepository) PublishScheduledEvent(topic string, event entity.ScheduledEvent) error {
	// async write
	return _self.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(event.Message),
	}, nil)
}
