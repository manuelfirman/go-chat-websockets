package users

import (
	"database/sql"

	"github.com/manuelfirman/go-chat-websockets/internal/domain"
)

// RepositoryMySQL is the repository for the user
type RepositoryMySQL struct {
	db *sql.DB
}

// NewRepositoryMySQL creates a new RepositoryMySQL for the user
func NewRepositoryMySQL(db *sql.DB) *RepositoryMySQL {
	return &RepositoryMySQL{db: db}
}

// GetAllUsers gets all users from the database
func (r *RepositoryMySQL) GetAllUsers() (users []domain.User, err error) {
	query := "SELECT id, name, email FROM User"

	rows, err := r.db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// iterate rows
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return
		}
		users = append(users, user)
	}

	// Check Error on Rows
	if err = rows.Err(); err != nil {
		return
	}

	return
}

// CreateUser creates a new user in the database
func (r *RepositoryMySQL) CreateUser(user *domain.User) (id int, err error) {
	query := "INSERT INTO User (name, email, password) VALUES (?, ?, ?)"

	result, err := r.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return
	}

	newUserId, err := result.LastInsertId()
	if err != nil {
		return
	}

	id = int(newUserId)

	return
}

// GetUserFromId gets a user from the database given an id
func (r *RepositoryMySQL) GetUserFromId(id int) (user *domain.User, err error) {
	query := "SELECT id, name, email, created_at FROM User WHERE id = (?)"

	row := r.db.QueryRow(query, id)
	// Iterate Rows
	user = &domain.User{}
	err = row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return
	}

	// Check Error on Rows
	if err = row.Err(); err != nil {
		return
	}

	return
}

// GetUserFromEmail gets a user from the database given an email
func (r *RepositoryMySQL) GetUserFromEmail(email string) (user *domain.User, err error) {
	query := "SELECT id, name, email, created_at, password FROM `User` WHERE email = (?)"

	row := r.db.QueryRow(query, email)

	user = &domain.User{}
	err = row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.Password)
	if err != nil {
		return
	}

	// Check Error on Rows
	if err = row.Err(); err != nil {
		return
	}

	return
}
