package user

import (
	"database/sql"
)

type API struct {
	storage *UserStorage
}

func New(db *sql.DB) *API {
	return &API{
		storage: NewUserStorage(db),
	}
}
