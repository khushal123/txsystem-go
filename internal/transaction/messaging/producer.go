package messaging

import (
	"context"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaProducer struct {
	client    *kgo.Client
	topic     string
	brokers   string
	connected bool
}

func connectProducer(brokers []string, topic string) *KafkaProducer {
	log.Debug("kafka brokers", brokers)
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
		kgo.WithLogger(kgo.BasicLogger(
			os.Stderr,        // log to stderr
			kgo.LogLevelInfo, // set log level (Debug, Info, Warn, Error)
			nil,              // use default formatter
		)),
	)
	if err != nil {
		log.Errorf("Failed to create Kafka client: %v", err)
		return nil
	}

	return &KafkaProducer{
		client: client,
		topic:  topic,
	}
}

func GetProducerConnection(brokers []string, topic string) ProducerConnection {
	log.Info("Connecting to Kafka brokers:", brokers)
	log.Info("Using topic:", topic)
	var instance *KafkaProducer = connectProducer(brokers, topic)
	if instance.client == nil {
		log.Error("Kafka client is nil")
		return nil
	}
	return instance
}

func (k *KafkaProducer) Produce(message string) error {
	log.Info("Producing message to Kafka topic:", k.topic)
	log.Info("Message:", k.IsConnected())
	record := &kgo.Record{
		Topic: k.topic,
		Value: []byte(message),
	}

	k.client.Produce(context.Background(), record, func(r *kgo.Record, err error) {
		if err != nil {
			log.Errorf("Failed to produce message: %v", err)
		}
	})

	return nil
}

// IsConnected checks if the client is still connected
func (k *KafkaProducer) IsConnected() bool {
	if k == nil || k.client == nil {
		return false
	}

	// Ping with a short timeout to check connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := k.client.Ping(ctx)
	if err != nil {
		k.connected = false
		return false
	}

	k.connected = true
	return true
}

func (k *KafkaProducer) Close() {
	if k.client != nil {
		k.client.Close()
	}
}
