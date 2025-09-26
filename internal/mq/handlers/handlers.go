package handlers

import (
	"encoding/json"
	"log"

	"github.com/lukashonok/micro-fiber-pet/internal/mq"
	"github.com/lukashonok/micro-fiber-pet/internal/mq/messages"
	"github.com/lukashonok/micro-fiber-pet/pkg/entities"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func MqGeneratedPdf(services mq.DefaultServices) func(amqp.Delivery) {
	return func(msg amqp.Delivery) {
		var pdf messages.MqPdfGeneratedMsg

		if err := json.Unmarshal(msg.Body, &pdf); err != nil {
			log.Printf("failed to unmarshal pdf: %v", err)
			return
		}

		ID, err := bson.ObjectIDFromHex(pdf.BookID)
		if err != nil {
			log.Printf("failed to fetch pdf: %v", err)
		}
		services.BookService.UpdateBook(&entities.Book{
			ID:  ID,
			Url: pdf.URL,
		})

		log.Printf("Book pdf url was updated: %v", err)
	}
}
