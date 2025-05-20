package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"txsystem/internal/account/processor"
	"txsystem/internal/common/messaging"
	"txsystem/internal/common/types"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// LoadEnv loads env vars; already done in init()
func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func run() {
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}

	msgProcessor := processor.NewMessageProcessor(db)

	consumer := setupKafkaConsumer()
	if err != nil {
		log.Fatalf("kafka consumer setup failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer.StartConsumer(ctx, msgProcessor)
	log.Info("Kafka consumer started...")

	waitForShutdown(cancel)

	log.Info("Shutting down consumer...")
	consumer.Close()
	log.Info("Shutdown complete")
}

func main() {
	run()
}

func setupDatabase() (*gorm.DB, error) {
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	fmt.Println(host, port, user, password, dbname, "this is database")
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupKafkaConsumer() types.ConsumerConnection {
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC")

	if brokers == "" || topic == "" {
		log.Fatal("KAFKA_BROKERS or KAFKA_TOPIC env var not set")
	}

	consumer := messaging.NewKafkaConsumer(strings.Split(brokers, ","), topic)
	return consumer
}

func waitForShutdown(cancelFunc context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Info("Shutdown signal received")
	cancelFunc()
	os.Exit(0)
	time.Sleep(2 * time.Second)
}
