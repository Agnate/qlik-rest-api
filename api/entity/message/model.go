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

type MessageService interface {
	Message(id int) (*Message, error)
	Messages() ([]*Message, error)
	CreateMessage(m *Message) error
}

// func (input *Input) ToModel() *Message {
// 	sanitzedMsg := input.Message
// 	// now := time.Now()
// 	isPalindrome := false

// 	return &Message{
// 		UUID:    input.UUID,
// 		Message: sanitzedMsg,
// 		// CreateDate: now,
// 		Palindrome: strconv.FormatBool(isPalindrome),
// 		// LastUpdated: now
// 		LastUpdatedBy: input.UUID,
// 	}
// }
