package user

type User struct {
	UUID       string `json:"user_id"`
	Email      string `json:"email"`
	APIKey     string `json:"api_key"`
	LastAccess string `json:"last_access"`
	CreateDate string `json:"create_date"`
}

type Users []*User

type UserService interface {
	User(id int) (*User, error)
	Users() ([]*User, error)
	CreateUser(u *User) error
	DeleteUser(id int) error
}
