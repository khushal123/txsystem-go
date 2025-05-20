package messaging

import (
	"context"
	"os"
	"time"
	"txsystem/internal/common/types"

	"github.com/labstack/gommon/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

type kafkaProducer struct {
	client     *kgo.Client
	topic      string
	connected  bool
	producerID []byte
}

func connectProducer(brokers []string, topic string) *kafkaProducer {
	log.Debug("kafka brokers", brokers)
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
		kgo.TransactionalID("transaction-producer-1"),
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

	return &kafkaProducer{
		client:     client,
		topic:      topic,
		producerID: []byte("transaction-producer-1"),
	}
}

func GetProducerConnection(brokers []string, topic string) types.ProducerConnection {
	log.Info("Connecting to Kafka brokers:", brokers)
	log.Info("Using topic:", topic)
	var instance *kafkaProducer = connectProducer(brokers, topic)
	if instance.client == nil {
		log.Error("Kafka client is nil")
		return nil
	}
	return instance
}

func (kp *kafkaProducer) Produce(message string) error {
	log.Info("Producing message to Kafka topic:", kp.topic)
	if err := kp.client.BeginTransaction(); err != nil {
		return err
	}
	record := &kgo.Record{
		Topic: kp.topic,
		Value: []byte(message),
		Key:   kp.producerID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := kp.client.ProduceSync(ctx, record); err != nil {
		log.Errorf("Failed to produce message: %v", err)
		if abortErr := kp.client.EndTransaction(ctx, kgo.TryCommit); abortErr != nil {
			log.Errorf("Additionally failed to abort transaction: %v", abortErr)
		}
		return nil
	}

	return nil
}

// IsConnected checks if the client is still connected
func (k *kafkaProducer) IsConnected() bool {
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

func (k *kafkaProducer) Close() {
	if k.client != nil {
		k.client.Close()
	}
}
