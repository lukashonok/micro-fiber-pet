package handlers

import (
	"fmt"
	"net/http"

	"github.com/lukashonok/micro-fiber-pet/api/presenter"
	"github.com/lukashonok/micro-fiber-pet/internal/mq"
	"github.com/lukashonok/micro-fiber-pet/pkg/book"
	"github.com/lukashonok/micro-fiber-pet/pkg/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// AddBook is handler/controller which creates Books in the BookShop
func AddBook(service book.Service, bookEventPublisher *mq.Publisher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.Book
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.BookErrorResponse(err))
		}
		if requestBody.Author == "" || requestBody.Title == "" {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.BookErrorResponse(errors.New(
				"Please specify title and author")))
		}
		result, err := service.InsertBook(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.BookErrorResponse(err))
		}

		bookEventPublisher.Publish("book.created", requestBody)
		fmt.Println("bookEventPublisher.Publish('book.created', requestBody) sent")
		return c.JSON(presenter.BookSuccessResponse(result))
	}
}

// UpdateBook is handler/controller which updates data of Books in the BookShop
func UpdateBook(service book.Service, bookEventPublisher *mq.Publisher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.Book
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.BookErrorResponse(err))
		}
		result, err := service.UpdateBook(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.BookErrorResponse(err))
		}
		bookEventPublisher.Publish("book.updated", requestBody)
		return c.JSON(presenter.BookSuccessResponse(result))
	}
}

// RemoveBook is handler/controller which removes Books from the BookShop
func RemoveBook(service book.Service, bookEventPublisher *mq.Publisher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.DeleteRequest
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.BookErrorResponse(err))
		}
		bookID := requestBody.ID
		if bookID == "0" {
			books, _ := service.FetchBooks()
			for _, b := range *books {
				err := service.RemoveBook(string(b.ID.Hex()))
				if err != nil {
					fmt.Println(err.Error())
				}
			}

			return c.JSON(&fiber.Map{
				"status": true,
				"data":   "EVERYTHING DELETED",
				"err":    nil,
			})
		}
		err = service.RemoveBook(bookID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.BookErrorResponse(err))
		}
		bookEventPublisher.Publish("book.deleted", requestBody)
		return c.JSON(&fiber.Map{
			"status": true,
			"data":   "deleted successfully",
			"err":    nil,
		})
	}
}

// GetBooks is handler/controller which lists all Books from the BookShop
func GetBooks(service book.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fetched, err := service.FetchBooks()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.BookErrorResponse(err))
		}
		return c.JSON(presenter.BooksSuccessResponse(fetched))
	}
}
