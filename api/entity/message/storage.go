package message

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
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
	return s.scanMessages("SELECT * FROM messages WHERE logical_delete = $1", false)
}

func (s *MessageStorage) ListByUUID(uuid uuid.UUID) (Messages, error) {
	return s.scanMessages("SELECT * FROM messages WHERE uuid = $1 AND logical_delete = $2", uuid, false)
}

func (s *MessageStorage) Read(uuid uuid.UUID, createDate time.Time) (*Message, error) {
	msgs, err := s.scanMessages("SELECT * FROM messages WHERE uuid = $1 AND create_date = $2 AND logical_delete = $3", uuid, createDate, false)
	if err == nil && len(msgs) > 0 {
		return msgs[0], nil
	}
	return nil, err
}

func (s *MessageStorage) scanMessages(query string, queryParams ...any) (Messages, error) {
	msgs := make([]*Message, 0)

	rows, err := s.db.Query(query, queryParams...)
	if err != nil {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		log.Println(err)
		return msgs, err
	}

	for rows.Next() {
		msg := &Message{}
		// TODO: Stop making assumptions about how data will be returned, since this fails if the db schema changes.
		rows.Scan(&msg.UUID, &msg.CreateDate, &msg.Message, &msg.Palindrome, &msg.LastUpdated, &msg.LastUpdatedBy, &msg.Deleted)
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (s *MessageStorage) Create(msg *Message) (*Message, error) {
	// Create the new Message.
	_, err := s.db.Exec("INSERT INTO messages(uuid, message, is_palindrome, last_updated_by) VALUES($1, $2, $3, $4)",
		msg.UUID, msg.Message, msg.Palindrome, msg.LastUpdatedBy)
	if err != nil {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		log.Println(err)
		return nil, err
	}
	// Retrieve the message.
	newMsg, err := s.getLatest(msg.UUID)
	if err != nil {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		log.Println(err)
		return nil, err
	}
	return newMsg, nil
}

func (s *MessageStorage) getLatest(uuid uuid.UUID) (*Message, error) {
	msgs, err := s.scanMessages("SELECT * FROM messages WHERE uuid = $1 AND logical_delete = $2 ORDER BY create_date DESC LIMIT 1", uuid, false)
	if err == nil && len(msgs) > 0 {
		return msgs[0], nil
	}
	return nil, err
}

func (s *MessageStorage) get(uuid uuid.UUID, createDate time.Time) (*Message, error) {
	msgs, err := s.scanMessages("SELECT * FROM messages WHERE uuid = $1 AND create_date = $2", uuid, createDate)
	if err == nil && len(msgs) > 0 {
		return msgs[0], nil
	}
	return nil, err
}

func (s *MessageStorage) Update(msg *Message) (*Message, error) {
	// Update the Message.
	result, err := s.db.Exec("UPDATE messages SET message = $1, is_palindrome = $2, last_updated_by = $3, last_updated = $4 "+
		"WHERE uuid = $5 AND create_date = $6 AND last_updated = $7 AND logical_delete = $8",
		msg.Message, msg.Palindrome, msg.LastUpdatedBy, time.Now(), msg.UUID, msg.CreateDate, msg.LastUpdated, false)
	if err != nil {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		log.Println(err)
		return nil, err
	}

	// Check if any rows were updated.
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		// TODO: Prefer to use custom errors here so we can have better error handling to users.
		return nil, errors.New("no rows updated")
	}

	// Retrieve the message.
	// TODO: Make this transaction-based to avoid inconsistent read (if another update/delete)
	// has happened in between UPDATE and READ.
	updatedMsg, err := s.Read(msg.UUID, msg.CreateDate)
	if err != nil {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		log.Println(err)
		return nil, err
	}
	return updatedMsg, nil
}

func (s *MessageStorage) Delete(msg *Message) (*Message, error) {
	// Delete the Message.
	result, err := s.db.Exec("UPDATE messages SET logical_delete = $1, last_updated_by = $2, last_updated = $3 "+
		"WHERE uuid = $4 AND create_date = $5 AND last_updated = $6 AND logical_delete = $7",
		true, msg.LastUpdatedBy, time.Now(), msg.UUID, msg.CreateDate, msg.LastUpdated, false)
	if err != nil {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		log.Println(err)
		return nil, err
	}

	// Check if any rows were updated.
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		return nil, errors.New("no rows updated")
	}

	// Retrieve the message.
	// TODO: Make this transaction-based to avoid inconsistent read (if another update/delete)
	// has happened in between UPDATE and READ.
	deletedMsg, err := s.get(msg.UUID, msg.CreateDate)
	if err != nil {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		log.Println(err)
		return nil, err
	}
	return deletedMsg, nil
}
