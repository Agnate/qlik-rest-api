package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID       uuid.UUID `json:"user_id"`
	Email      string    `json:"-"`
	APIKey     string    `json:"-"`
	LastAccess time.Time `json:"last_access"`
	CreateDate time.Time `json:"-"`
	Name       string    `json:"name"`
}

type Users []*User
