package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"txsystem/internal/ledger/handlers"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "ledger_db"
	}

	return client.Database(dbName), nil
}

func setupEchoServer(db *mongo.Database) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	handlers.InitRoutes(e, db)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	return e
}

func waitForShutdown(e *echo.Echo, db *mongo.Client) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Disconnect(ctx); err != nil {
		log.Errorf("Error disconnecting from MongoDB: %v", err)
	}

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func run() {
	// Connect to MongoDB
	db, err := setupMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Info("Connected to MongoDB")

	// Setup Echo server
	e := setupEchoServer(db)

	// Start server
	port := os.Getenv("LEDGER_SERVICE_PORT")
	if port == "" {
		port = "9003"
	}

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
			log.Infof("Shutting down: %v", err)
		}
	}()

	log.Infof("Ledger service started on port %s", port)

	// Wait for shutdown signal
	waitForShutdown(e, db.Client())
	log.Info("Server gracefully stopped")
}

func main() {
	run()
}
