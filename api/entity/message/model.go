package message

// type Input struct {
// 	UUID    string `json:"user_id"`
// 	Message string `json:"message"`
// }

type Message struct {
	UUID          string `json:"user_id"`
	CreateDate    string `json:"create_date"`
	Message       string `json:"message"`
	Palindrome    string `json:"is_palindrome"`
	LastUpdated   string `json:"last_updated_date"`
	LastUpdatedBy string `json:"last_updated_by"`
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
