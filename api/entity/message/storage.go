package message

import (
	"database/sql"
	"log"
)

type MessageStorage struct {
	db *sql.DB
}

func NewMessageStorage(db *sql.DB) *MessageStorage {
	return &MessageStorage{
		db: db,
	}
}

func (s *MessageStorage) List() (Messages, error) {
	msgs := make([]*Message, 0)

	rows, err := s.db.Query("SELECT * FROM messages")
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		msg := &Message{}
		rows.Scan(&msg.UUID, &msg.CreateDate, &msg.Message, &msg.Palindrome, &msg.LastUpdated, &msg.LastUpdatedBy)
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

// func (m *MessageStorage) Read(id uuid.UUID, createDate uuid.Time) (*Message, error) {
// 	msg := &Message{}
// 	return msg, nil
// }
