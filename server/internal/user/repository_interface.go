package users

import "github.com/manuelfirman/go-chat-websockets/internal/domain"

type Repository interface {
	// GetAllUsers gets all users from the database
	GetAllUsers() (users []domain.User, err error)
	// CreateUser creates a new user in the database
	CreateUser(user *domain.User) (id int, err error)
	// GetUserFromId gets a user from the database
	GetUserFromId(id int) (user *domain.User, err error)
	// GetUserFromEmail gets a user from the database
	GetUserFromEmail(email string) (user *domain.User, err error)
}
