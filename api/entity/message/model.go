package message

import (
	"strconv"
)

type Input struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Message struct {
	Username   string `json:"username"`
	CreateDate string `json:"create_date"`
	Message    string `json:"message"`
	Palindrome string `json:"is_palindrome"`
}

func (input *Input) ToModel() *Message {
	sanitzedMsg := input.Message
	// createDate := time.Now()
	isPalindrome := false

	return &Message{
		Username: input.Username,
		Message:  sanitzedMsg,
		// CreateDate: createDate,
		Palindrome: strconv.FormatBool(isPalindrome),
	}
}
