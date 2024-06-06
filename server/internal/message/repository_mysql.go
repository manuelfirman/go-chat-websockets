package message

import (
	"database/sql"

	"github.com/manuelfirman/go-chat-websockets/internal/domain"
)

// Repository is the repository for the user
type RepositoryMySQL struct {
	db *sql.DB
}

// NewRepositoryMySQL creates a new RepositoryMySQL for the user
func NewRepositoryMySQL(db *sql.DB) *RepositoryMySQL {
	return &RepositoryMySQL{db: db}
}

// GetGroupMessages gets the messages between two users
func (r *RepositoryMySQL) GetGroupMessages(groupId int) (messages []domain.Message, err error) {
	query := "SELECT * FROM GroupMessage WHERE group_id = (?)"

	rows, err := r.db.Query(query, groupId)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var message domain.Message
		var sentBytes []uint8
		err = rows.Scan(&message.Id, &message.UserId, &message.ToId, &message.Message, &sentBytes)
		message.Sent = string(sentBytes)
		if err != nil {
			return
		}
		messages = append(messages, message)
	}

	// Check Error on Rows
	if err = rows.Err(); err != nil {
		return
	}

	return

}

// GetUserMessages gets the messages between two users
func (r *RepositoryMySQL) GetUserMessages(userId int, friendId int) (messages []domain.Message, err error) {
	query := "SELECT * FROM UserMessage WHERE (user_from = (?) OR user_to = (?)) AND (user_from = (?) OR user_to = (?))"

	rows, err := r.db.Query(query, userId, userId, friendId, friendId)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var message domain.Message
		var sentBytes []uint8
		err = rows.Scan(&message.Id, &message.UserId, &message.ToId, &message.Message, &sentBytes)
		message.Sent = string(sentBytes)
		if err != nil {
			return
		}
		messages = append(messages, message)
	}

	// Check Error on Rows
	if err = rows.Err(); err != nil {
		return
	}

	return
}
