package user

import (
	"time"

	"github.com/google/uuid"
)

type UserInput struct {
	Name  string `json:"full_name"`
	Email string `json:"email"`
}

type User struct {
	UUID       uuid.UUID `json:"user_id"`
	Email      string    `json:"-"`
	APIKey     string    `json:"api_key"`
	LastAccess time.Time `json:"last_access"`
	CreateDate time.Time `json:"-"`
	Name       string    `json:"name"`
	RawAPIKey  string    `json:"-"`
}

type Users []*User
