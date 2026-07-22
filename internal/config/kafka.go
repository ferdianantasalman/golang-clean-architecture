package config

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func NewKafkaConsumer(log *logrus.Logger) *kafka.Consumer {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset": os.Getenv("KAFKA_AUTO_OFFSET_RESET"),
	}

	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	return consumer
}

func NewKafkaProducer(log *logrus.Logger) *kafka.Producer {
	enabled, _ := strconv.ParseBool(os.Getenv("KAFKA_PRODUCER_ENABLED"))
	if !enabled {
		log.Info("Kafka producer is disabled")
		return nil
	}

	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
	}

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	return producer
}
