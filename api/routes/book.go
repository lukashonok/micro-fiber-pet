package routes

import (
	"github.com/lukashonok/micro-fiber-pet/api/handlers"
	"github.com/lukashonok/micro-fiber-pet/internal/mq"

	"github.com/gofiber/fiber/v2"
)

// BookRouter is the Router for GoFiber App
func BookRouter(app fiber.Router, defaultServices mq.DefaultServices) {
	app.Get("/books", handlers.GetBooks(defaultServices.BookService))
	app.Post("/books", handlers.AddBook(defaultServices.BookService, defaultServices.BookPublisher))
	app.Put("/books", handlers.UpdateBook(defaultServices.BookService, defaultServices.BookPublisher))
	app.Delete("/books", handlers.RemoveBook(defaultServices.BookService, defaultServices.BookPublisher))
}
