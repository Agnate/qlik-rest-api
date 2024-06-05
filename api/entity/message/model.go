package message

import (
	"time"

	"github.com/google/uuid"
)

type MessageInput struct {
	Message     string    `json:"message"`
	LastUpdated time.Time `json:"last_updated_date"` // Required for Updates to validate concurrency
}

type Message struct {
	UUID          uuid.UUID `json:"user_id"`
	CreateDate    time.Time `json:"create_date"`
	Message       string    `json:"message"`
	Palindrome    bool      `json:"is_palindrome"`
	LastUpdated   time.Time `json:"last_updated_date"`
	LastUpdatedBy uuid.UUID `json:"last_updated_by"`
	Deleted       bool      `json:"-"`
}

type Messages []*Message
