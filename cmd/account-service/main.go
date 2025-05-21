package main

import (
	"fmt"
	"os"
	"strings"
	"txsystem/internal/account/handler"
	"txsystem/internal/account/models"
	"txsystem/pkg/common/messaging"
	"txsystem/pkg/common/types"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// LoadEnv loads env vars; already done in init()
func init() {
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
	if err := db.AutoMigrate(&models.Account{}); err != nil {
		return nil, err
	}

	return db, nil
}

func setupProducer() types.ProducerConnection {
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC_TRANSCATIONS")

	if brokers == "" || topic == "" {
		log.Fatal("KAFKA_BROKERS or KAFKA_TOPIC_TRANSCATIONS env var not set")
	}

	conn := messaging.GetProducerConnection(strings.Split(brokers, ","), topic)
	if conn == nil || !conn.IsConnected() {
		log.Fatal("Failed to connect Kafka producer")
	}

	return conn
}

func setupEchoServer(kafkaProducer types.ProducerConnection, db *gorm.DB) *echo.Echo {
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

func main() {
	run()
}
