package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID       uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	APIKey     string    `json:"api_key"`
	LastAccess time.Time `json:"last_access"`
	CreateDate time.Time `json:"create_date"`
	Name       string    `json:"name"`
}

type Users []*User

type UserService interface {
	User(id int) (*User, error)
	Users() ([]*User, error)
	CreateUser(u *User) error
	DeleteUser(id int) error
}
