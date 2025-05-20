package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	_ "txsystem/internal/transaction/docs"
	"txsystem/internal/transaction/handler"
	"txsystem/internal/transaction/messaging"
	"txsystem/internal/transaction/models"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// LoadEnv loads env vars; already done in init()

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

	log.Info("Migrating database...")
	if err := db.AutoMigrate(&models.Transaction{}); err != nil {
		return nil, err
	}

	return db, nil
}

func setupProducer() messaging.ProducerConnection {
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC")

	if brokers == "" || topic == "" {
		log.Fatal("KAFKA_BROKERS or KAFKA_TOPIC env var not set")
	}

	conn := messaging.GetProducerConnection(strings.Split(brokers, ","), topic)
	if conn == nil || !conn.IsConnected() {
		log.Fatal("Failed to connect Kafka producer")
	}

	return conn
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

func setupEchoServer(kafkaProducer messaging.ProducerConnection, db *gorm.DB) *echo.Echo {
	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10,
		LogLevel:  log.ERROR,
	}))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	handler.InitRoutes(e, kafkaProducer, db)

	e.GET("/", func(c echo.Context) error {
		kafkaProducer.Produce("server started")
		return c.String(200, "OK")
	})

	return e
}

func run() {
	db, err := setupDatabase()
	if err != nil {
		log.Fatal("Database setup failed:", err)
	}

	producer := setupProducer()
	defer producer.Close()

	consumer := setupConsumer()
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer.StartConsumer(ctx)

	echoServer := setupEchoServer(producer, db)

	port := os.Getenv("ACCOUNT_SERVICE_PORT")
	if port == "" {
		port = "9001"
	}
	log.Infof("Starting server on port %s", port)
	if err := echoServer.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Echo server failed:", err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	run()
}
