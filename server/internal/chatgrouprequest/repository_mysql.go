package chatgrouprequest

import (
	"database/sql"

	"github.com/manuelfirman/go-chat-websockets/internal/domain"
)

// RepositoryMySQL is the repository for the chat group request
type RepositoryMySQL struct {
	db *sql.DB
}

// NewRepositoryMySQL creates a new RepositoryMySQL for the chat group request
func NewRepositoryMySQL(db *sql.DB) *RepositoryMySQL {
	return &RepositoryMySQL{db: db}
}

// InsertChatGroupRequest inserts a new chat group request
func (r *RepositoryMySQL) InsertChatGroupRequest(id int, groupId int) (err error) {
	query := "INSERT INTO GroupRequest (user_id, group_id) VALUES (?, ?)"

	_, err = r.db.Exec(query, id, groupId)
	if err != nil {
		return err
	}
	return
}

// GetChatGroupRequests gets the chat group requests
func (r *RepositoryMySQL) GetChatGroupRequests(id, groupId int, sented bool) (requests []domain.ChatGroupRequest, err error) {
	var query string
	var param int

	if sented {
		query = "SELECT * FROM GroupRequest WHERE user_id = (?)"
		param = id
	} else {
		query = "SELECT * FROM GroupRequest WHERE group_id = (?)"
		param = groupId
	}

	rows, err := r.db.Query(query, param)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var request domain.ChatGroupRequest
		var sentBytes []uint8
		err = rows.Scan(&request.Id, &request.UserId, &request.GroupId, &sentBytes)
		if err != nil {
			return
		}
		request.Sent = string(sentBytes)
		requests = append(requests, request)
	}

	// Check Error on Rows
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (r *RepositoryMySQL) GetGroupOwner(groupId int) (owner *int, err error) {
	query := "SELECT owner FROM ChatGroup WHERE id = (?)"

	err = r.db.QueryRow(query, groupId).Scan(&owner)
	if err != nil {
		return
	}

	return
}

func (r *RepositoryMySQL) AcceptGroupRequest(userId int, groupId int) (err error) {
	query := "CALL AcceptGroupRequest(?, ?)"
	_, err = r.db.Exec(query, userId, groupId)
	if err != nil {
		return
	}
	return
}
