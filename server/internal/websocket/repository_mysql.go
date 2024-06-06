package websocket

import (
	"database/sql"

	"github.com/manuelfirman/go-chat-websockets/internal/domain"
)

// RepositoryMySQL is a struct to manage the MySQL database.
type RepositoryMySQL struct {
	db *sql.DB
}

// NewRepositoryMySQL creates a new RepositoryMySQL.
func NewRepositoryMySQL(db *sql.DB) *RepositoryMySQL {
	return &RepositoryMySQL{db}
}

func (r *RepositoryMySQL) InsertGroupMessage(id int, message *domain.Message) (err error) {
	// query to call the stored procedure
	query := "CALL InsertGroupMessage(?, ?, ?)"

	_, err = r.db.Exec(query, id, message.ToId, message.Message)
	if err != nil {
		return
	}

	return
}

func (r *RepositoryMySQL) InsertUserMessage(id int, message *domain.Message) (err error) {
	// query to call the stored procedure
	query := "CALL InsertUserMessage(?, ?, ?)"

	_, err = r.db.Exec(query, id, message.ToId, message.Message)
	if err != nil {
		return
	}

	return
}

// GetGroupsOfUser gets the groups of a user
func (r *RepositoryMySQL) GetGroupsOfUser(userId int) (groups []int, err error) {
	query := "SELECT group_id FROM UserChatGroup WHERE user_id = (?)"

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return
	}
	defer rows.Close()

	// Iterate Rows
	for rows.Next() {
		var chatId int
		err = rows.Scan(&chatId)
		if err != nil {
			return
		}
		groups = append(groups, chatId)
	}

	// Check Error on Rows
	if err = rows.Err(); err != nil {
		return
	}

	return
}
