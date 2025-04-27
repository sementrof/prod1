package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

type ProducerWrapper struct {
	producer *kafka.Producer
}
type Producer interface {
	Publish(topic string, message []byte) error
	Close()
}

func NewProducer(config *kafka.ConfigMap) (Producer, error) {
	producer, err := kafka.NewProducer(config)
	if err != nil {
		logrus.Errorf("failed: %s", err)
		return nil, err
	}
	return &ProducerWrapper{producer}, nil
}

func (w *ProducerWrapper) Publish(topic string, message []byte) error {
	return w.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
}
func (w *ProducerWrapper) Close() {
	// Ожидаем отправку всех сообщений в течение 15 секунд (15000 миллисекунд)
	w.producer.Flush(1500)
	w.producer.Close()
}

func InitKafka() (Producer, error) {
	config := &kafka.ConfigMap{
		// "bootstrap.servers": os.Getenv("PORT_KAFKA1") + "," + os.Getenv("PORT_KAFKA2"),
		"bootstrap.servers": "localhost:9092,localhost:9093",
		"acks":              "all",
		"client.id":         "myProducer",
		"debug":             "msg",
	}
	producer, err := NewProducer(config)
	if err != nil {
		logrus.Errorf("failed: %s", err)
		return nil, err
	}
	return producer, nil

}
