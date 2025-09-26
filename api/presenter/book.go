package presenter

import (
	"github.com/lukashonok/micro-fiber-pet/pkg/entities"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/gofiber/fiber/v2"
)

// Book is the presenter object which will be passed in the response by Handler
type Book struct {
	ID     bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title  string        `json:"title"`
	Author string        `json:"author"`
	Url    string        `json:"url"`
}

// BookSuccessResponse is the singular SuccessResponse that will be passed in the response by
// Handler
func BookSuccessResponse(data *entities.Book) *fiber.Map {
	book := Book{
		ID:     data.ID,
		Title:  data.Title,
		Author: data.Author,
		Url:    data.Url,
	}
	return &fiber.Map{
		"status": true,
		"data":   book,
		"error":  nil,
	}
}

// BooksSuccessResponse is the list SuccessResponse that will be passed in the response by Handler
func BooksSuccessResponse(data *[]Book) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data":   data,
		"error":  nil,
	}
}

// BookErrorResponse is the ErrorResponse that will be passed in the response by Handler
func BookErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"data":   "",
		"error":  err.Error(),
	}
}
