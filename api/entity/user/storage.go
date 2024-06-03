package user

import (
	"database/sql"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (us *UserStorage) List() (Users, error) {
	users := make([]*User, 0)
	return users, nil
}
