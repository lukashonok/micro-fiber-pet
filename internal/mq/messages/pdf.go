package messages

import (
	"time"
)

type MqPdfGeneratedMsg struct {
	Event     string    `json:"event"` // "pdf.generated"
	BookID    string    `json:"book_id"`
	ObjectKey string    `json:"object_key"`
	URL       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
}

type MqPdfDeletedMsg struct {
	Event     string    `json:"event"` // "pdf.deleted"
	BookID    string    `json:"book_id"`
	ObjectKey string    `json:"object_key"`
	Timestamp time.Time `json:"timestamp"`
}
