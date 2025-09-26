package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/lukashonok/micro-fiber-pet/api/routes"
	"github.com/lukashonok/micro-fiber-pet/internal/mq"
	mqhandlers "github.com/lukashonok/micro-fiber-pet/internal/mq/handlers"
	"github.com/lukashonok/micro-fiber-pet/pkg/book"
	"github.com/lukashonok/micro-fiber-pet/pkg/firebase"

	_ "github.com/lukashonok/micro-fiber-pet/cmd/migrations"
	amqp "github.com/rabbitmq/amqp091-go"
	migrate "github.com/xakep666/mongo-migrate"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	defaultServices, err := buildDefaultServices()
	if err != nil {
		log.Fatalf("services building failed: %v", err.Error())
	}
	setupInfrastructure(defaultServices)

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Golang API Works!"))
	})

	api := app.Group("/api")
	routes.BookRouter(api, defaultServices)
	routes.AuthRouter(api, defaultServices) // теперь принимает defaultServices с UserService

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}

func databaseConnection() (*mongo.Database, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	connectURI := os.Getenv("DB_DSN")
	client, err := mongo.Connect(options.Client().ApplyURI(connectURI).
		SetServerSelectionTimeout(5 * time.Second))
	if err != nil {
		cancel()
		return nil, nil, err
	}
	db := client.Database(os.Getenv("DB_NAME"))
	migrate.SetDatabase(db)
	if err := migrate.Up(ctx, migrate.AllAvailable); err != nil {
		return nil, cancel, err
	}
	return db, cancel, nil
}

func buildDefaultServices() (mq.DefaultServices, error) {
	services := mq.DefaultServices{}

	db, cancel, err := databaseConnection()
	if err != nil {
		log.Fatalf("Database Connection Error %s", err)
	}
	defer cancel()
	services.Db = db

	// RabbitMQ config
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("missing required environment variables (MINIO_*, RABBITMQ_URL, FONT_PATH)")
	}
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return services, fmt.Errorf("dial rabbitmq: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return services, fmt.Errorf("open channel: %w", err)
	}
	rabbitMQ := mq.RabbitMQ{
		Connection: conn,
		Channel:    ch,
	}
	services.RabbitMQ = rabbitMQ

	services.BookPublisher, err = mq.NewPublisher(rabbitMQ, "micro_exchange")
	if err != nil {
		return services, fmt.Errorf("error in creating publisher: %w", err)
	}

	// Book service
	bookCollection := services.Db.Collection("books")
	bookRepo := book.NewRepo(bookCollection)
	bookService := book.NewService(bookRepo)
	services.BookService = bookService

	authCredFile := os.Getenv("FIREBASE_CONFIG_FILE")
	services.FirebaseAuth = firebase.NewAuth(context.Background(), authCredFile)

	return services, nil
}

func setupInfrastructure(services mq.DefaultServices) {
	// Starting Consumers for actions
	pdfConsumer, err := mq.NewConsumer(services.RabbitMQ, "micro_exchange", "pdf_queue_fiber",
		services,
		map[string]mq.MqHandler{
			"pdf.generated": mqhandlers.MqGeneratedPdf,
		})
	if err != nil {
		log.Fatalf("rabbitmq new consumer error: %v", err)
	}
	pdfConsumer.Start()
}
