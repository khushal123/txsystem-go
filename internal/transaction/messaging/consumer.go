package messaging

import (
	"context"
	"os"
	"time"
	"txsystem/internal/transaction/service"

	"github.com/labstack/gommon/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaConsumer struct {
	client    *kgo.Client
	topic     string
	groupID   string
	connected bool
}

func NewKafkaConsumer(brokers []string, topic string) *KafkaConsumer {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, nil)),
	)
	if err != nil {
		log.Errorf("Failed to create Kafka consumer client: %v", err)
		return nil
	}

	return &KafkaConsumer{
		client: client,
		topic:  topic,
	}
}
func (kc *KafkaConsumer) Consume(ctx context.Context, handler func(key, value []byte) error, ms *service.MessageProcessor) {
	go func() {
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

				for _, r := range records {
					if err := handler(r.Key, r.Value); err != nil {
						log.Errorf("Message handler error: %v", err)
					}
					kc.client.CommitRecords(ctx, r)
				}
			}
		}
	}()
}

func (kc *KafkaConsumer) StartConsumer(ctx context.Context, ms *service.MessageProcessor) {
	// Define your message handler here or accept it as param if you want flexibility
	handler := func(key, value []byte) error {
		log.Infof("Consumed message key=%s value=%s", string(key), string(value))
		// put your actual processing logic here, or call another internal function
		return nil
	}

	kc.Consume(ctx, handler, ms)
}

func (kc *KafkaConsumer) IsConnected() bool {
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

func (kc *KafkaConsumer) Close() {
	if kc.client != nil {
		kc.client.Close()
	}
}
