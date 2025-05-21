package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"txsystem/internal/ledger/processor"
	"txsystem/pkg/common/messaging"
	"txsystem/pkg/common/types"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Warn("Error loading .env file")
	}
}

func setupMongoDB() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "ledger_db"
	}

	return client.Database(dbName), nil
}

func setupKafkaConsumer() types.ConsumerConnection {
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC_TRANSCATIONS")

	if brokers == "" || topic == "" {
		log.Fatal("KAFKA_BROKERS or KAFKA_TOPIC_TRANSCATIONS env var not set")
	}

	consumer := messaging.NewKafkaConsumer(strings.Split(brokers, ","), topic)
	if consumer == nil || !consumer.IsConnected() {
		log.Fatal("Failed to connect to Kafka")
	}

	return consumer
}

func waitForShutdown(cancel context.CancelFunc, db *mongo.Database) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutdown signal received")
	cancel()

	ctx, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	if err := db.Client().Disconnect(ctx); err != nil {
		log.Errorf("Error disconnecting from MongoDB: %v", err)
	}

	log.Info("Gracefully shut down")
	os.Exit(0)
}

func run() {
	db, err := setupMongoDB()
	if err != nil {
		log.Fatalf("Failed to setup MongoDB: %v", err)
	}

	msgProcessor := processor.NewMessageProcessor(db)

	// Set up Kafka consumer
	consumer := setupKafkaConsumer()
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming messages
	consumer.StartConsumer(ctx, msgProcessor)
	log.Info("Ledger service started and consuming messages...")

	// Wait for shutdown signal
	waitForShutdown(cancel, db)
}

func main() {
	run()
}
