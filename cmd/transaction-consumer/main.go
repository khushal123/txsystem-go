package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"txsystem/internal/transaction/messaging"
	"txsystem/internal/transaction/models"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func setupDatabase() (*gorm.DB, error) {
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Info("Migrating database...")
	if err := db.AutoMigrate(&models.Transaction{}); err != nil {
		return nil, err
	}

	return db, nil
}

func setupConsumer() *messaging.KafkaConsumer {
	brokers := os.Getenv("KAFKA_CONSUMER_BROKERS")
	topic := os.Getenv("KAFKA_CONSUMER_TOPIC")

	if brokers == "" || topic == "" {
		log.Fatal("Consumer Kafka config missing in env vars")
	}
	log.Info("Creating Kafka consumer...")

	consumer := messaging.NewKafkaConsumer(strings.Split(brokers, ","), topic)
	if consumer == nil || !consumer.IsConnected() {
		log.Fatal("Failed to connect Kafka consumer")
	}

	return consumer
}

func main() {
	loadEnv()

	db, err := setupDatabase()
	if err != nil {
		log.Fatal("Database setup failed:", err)
	}

	consumer := setupConsumer()
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer.StartConsumer(ctx)

	// Block main goroutine, e.g., wait for signal to exit
	select {}
}
