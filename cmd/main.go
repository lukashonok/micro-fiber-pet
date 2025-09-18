package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lukashonok/micro-fiber-pet/api/routes"
	"github.com/lukashonok/micro-fiber-pet/pkg/book"
	migrate "github.com/xakep666/mongo-migrate"

	_ "github.com/lukashonok/micro-fiber-pet/cmd/migrations"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	db, cancel, err := databaseConnection()
	if err != nil {
		log.Fatal("Database Connection Error $s", err)
	}
	fmt.Println("Database connection success!")
	bookCollection := db.Collection("books")
	bookRepo := book.NewRepo(bookCollection)
	bookService := book.NewService(bookRepo)

	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Test Golang API Works!"))
	})
	api := app.Group("/api")
	routes.BookRouter(api, bookService)
	defer cancel()
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
