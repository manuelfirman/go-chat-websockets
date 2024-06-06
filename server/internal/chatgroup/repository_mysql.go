package chatgroup

import (
	"database/sql"

	"github.com/manuelfirman/go-chat-websockets/internal/domain"
)

// RepositoryMySQL is the repository for the chat group
type RepositoryMySQL struct {
	db *sql.DB
}

// NewRepositoryMySQL creates a new repository
func NewRepositoryMySQL(db *sql.DB) *RepositoryMySQL {
	return &RepositoryMySQL{db: db}
}

// CreateGroup creates a new group
func (r *RepositoryMySQL) CreateGroup(id int, name, description string) (err error) {
	query := "CALL CreateChatGroup(?, ?, ?)"

	_, err = r.db.Exec(query, name, description, id)
	if err != nil {
		return err
	}
	return
}

// GetGroups gets the groups
func (r *RepositoryMySQL) GetGroups(id int) (groups []domain.ChatGroup, err error) {
	query := "SELECT * FROM ChatGroup CG WHERE CG.id IN ( SELECT UCG.group_id FROM UserChatGroup UCG WHERE UCG.user_id = (?))"

	rows, err := r.db.Query(query, id)
	if err != nil {
		return
	}
	defer rows.Close()

	var sentBytes []uint8
	for rows.Next() {
		var chatGroup domain.ChatGroup
		err = rows.Scan(&chatGroup.Id, &chatGroup.Name, &chatGroup.Description, &sentBytes, &chatGroup.Owner)
		if err != nil {
			return
		}

		chatGroup.Sent = string(sentBytes)
		groups = append(groups, chatGroup)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

// GetGroupUsers gets the users of a group
func (r *RepositoryMySQL) GetGroupUsers(groupId int) (users []domain.User, err error) {
	query := "SELECT U.id, U.name, U.email FROM User U INNER JOIN UserChatGroup UCG ON U.id = UCG.user_id WHERE group_id = (?)"

	rows, err := r.db.Query(query, groupId)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Iterate Rows
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

// GetGroupOwner gets the owner of a group
func (r *RepositoryMySQL) GetGroupOwner(id int) (owner *int, err error) {
	query := "SELECT owner FROM ChatGroup WHERE id = (?)"

	rows, err := r.db.Query(query, id)
	if err != nil {
		return
	}
	defer rows.Close()

	// Iterate Rows
	rows.Next()
	err = rows.Scan(&owner)
	if err != nil {
		return
	}

	// Check Error on Rows
	if err = rows.Err(); err != nil {
		return
	}

	return
}
