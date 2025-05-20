package messaging

import (
	"context"
	"os"
	"time"
	"txsystem/internal/common/types"

	"github.com/labstack/gommon/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

type kafkaConsumer struct {
	client    *kgo.Client
	topic     string
	connected bool
}

func NewKafkaConsumer(brokers []string, topic string) types.ConsumerConnection {
	groupID := "transaction-consumer-group"

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
		kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, nil)),
	)
	if err != nil {
		log.Errorf("Failed to create Kafka consumer client: %v", err)
		return nil
	}

	return &kafkaConsumer{
		client: client,
		topic:  topic,
	}
}

func (kc *kafkaConsumer) consume(ctx context.Context, handler func(key, value []byte) error) {
	go func() {
		log.Infof("Starting Kafka consumer for topic: %s", kc.topic)

		for {
			select {
			case <-ctx.Done():
				log.Info("Kafka consumer shutting down")
				return
			default:
				fetches := kc.client.PollFetches(ctx)
				if errs := fetches.Errors(); len(errs) > 0 {
					for _, e := range errs {
						log.Errorf("Kafka fetch error: %v", e.Err)
					}
				}

				records := fetches.Records()
				if len(records) == 0 {
					time.Sleep(100 * time.Millisecond)
					continue
				}

				log.Infof("Received %d messages", len(records))

				for _, r := range records {
					// Add additional logging for debugging
					log.Debugf("Processing message: key=%s offset=%d partition=%d",
						string(r.Key), r.Offset, r.Partition)

					// Process the message
					err := handler(r.Key, r.Value)
					if err == nil {
						if commitErr := kc.client.CommitRecords(ctx, r); commitErr != nil {
							log.Errorf("Failed to commit offset: %v", commitErr)
						} else {
							log.Debugf("Successfully committed offset for key=%s", string(r.Key))
						}
					}
				}
			}
		}
	}()
}

func (kc *kafkaConsumer) StartConsumer(ctx context.Context, ms types.MessageProcessor) {
	handler := func(key, value []byte) error {
		log.Infof("Processing message: key=%s", string(key))
		return ms.ProcessMessage(string(value))
	}

	kc.consume(ctx, handler)
}

func (kc *kafkaConsumer) IsConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := kc.client.Ping(ctx)
	if err != nil {
		kc.connected = false
		return false
	}
	kc.connected = true
	return true
}

func (kc *kafkaConsumer) Close() {
	if kc.client != nil {
		kc.client.Close()
	}
}
