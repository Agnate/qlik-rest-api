package user

import (
	"database/sql"
	"log"
)

type UserStorage struct {
	db *sql.DB
}

// Create a new User storage container/service.
func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

// Retrieve a list of Users.
func (s *UserStorage) List() (Users, error) {
	return s.scanUsers("SELECT * FROM users")
}

// Create a new User and retrieve them.
func (s *UserStorage) Create(user *User) (*User, error) {
	// Create the new User.
	_, err := s.db.Exec("INSERT INTO users(full_name, email, api_key) VALUES($1, $2, $3)",
		user.Name, user.Email, user.APIKey)
	if err != nil {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		log.Println(err)
		return nil, err
	}
	// Retrieve the new User.
	newUser, err := s.getNewUser(user)
	if err != nil {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		log.Println(err)
		return nil, err
	}
	return newUser, nil
}

// Get a User by their API key.
func (s *UserStorage) GetUserByAPIKey(apiKey string) (*User, error) {
	users, err := s.scanUsers("SELECT * FROM users WHERE api_key = $1 LIMIT 1", apiKey)
	if err == nil && len(users) > 0 {
		return users[0], nil
	}
	return nil, err
}

func (s *UserStorage) getNewUser(user *User) (*User, error) {
	users, err := s.scanUsersIncludeAPIKey("SELECT * FROM users WHERE full_name = $1 AND email = $2 AND api_key = $3 LIMIT 1",
		user.Name, user.Email, user.APIKey)
	if err == nil && len(users) > 0 {
		return users[0], nil
	}
	return nil, err
}

func (s *UserStorage) scanUsers(query string, queryParams ...any) (Users, error) {
	users, err := s.scanUsersIncludeAPIKey(query, queryParams...)
	if err != nil {
		return users, err
	}

	// Strip off API key hash.
	for _, user := range users {
		user.APIKey = ""
	}
	return users, nil
}

func (s *UserStorage) scanUsersIncludeAPIKey(query string, queryParams ...any) (Users, error) {
	users := make([]*User, 0)

	rows, err := s.db.Query(query, queryParams...)
	if err != nil {
		// TODO: Database errors should have better logging so they can be monitored and fixed.
		log.Println(err)
		return users, err
	}

	for rows.Next() {
		user := &User{}
		// TODO: Stop making assumptions about how data will be returned, since this fails if the db schema changes.
		rows.Scan(&user.UUID, &user.Email, &user.APIKey, &user.LastAccess, &user.CreateDate, &user.Name)
		users = append(users, user)
	}

	return users, nil
}
