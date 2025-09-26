package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Book Constructs your Book model under entities.
type Book struct {
	ID        bson.ObjectID `json:"id"  bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title,omitempty"`
	Author    string        `json:"author" bson:"author,omitempty"`
	Url       string        `json:"url" bson:"url"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// DeleteRequest struct is used to parse Delete Requests for Books
type DeleteRequest struct {
	ID string `json:"id"`
}
