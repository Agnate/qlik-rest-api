package user

import (
	"database/sql"
	"log"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) List() (Users, error) {
	users := make([]*User, 0)

	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		user := &User{}
		rows.Scan(&user.UUID, &user.Email, &user.APIKey, &user.LastAccess, &user.CreateDate, &user.Name)
		users = append(users, user)
	}

	return users, nil
}
