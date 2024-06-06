package friend

import (
	"database/sql"

	"github.com/manuelfirman/go-chat-websockets/internal/domain"
)

// RepositoryMySQL is the repository for the friend
type RepositoryMySQL struct {
	db *sql.DB
}

// NewRepositoryMySQL creates a new RepositoryMySQL for the friend
func NewRepositoryMySQL(db *sql.DB) *RepositoryMySQL {
	return &RepositoryMySQL{db: db}
}

// InsertFriendRequest inserts a new friend request in the database
func (r *RepositoryMySQL) InsertFriendRequest(id int, friendId int) (err error) {
	// query to call store procedure
	query := "CALL CreateFriendRequest(?, ?)"

	_, err = r.db.Exec(query, id, friendId)
	if err != nil {
		return
	}

	return

}

// AcceptFriendRequest accepts a friend request in the database
func (r *RepositoryMySQL) AcceptFriendRequest(id int, friendId int) (rowsAffected int, err error) {
	// query to call store procedure
	query := "CALL AcceptFriendRequest(?, ?)"

	result, err := r.db.Exec(query, id, friendId)
	if err != nil {
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return
	}

	rowsAffected = int(rows)

	return
}

// GetFriendRequests gets the friend requests
func (r *RepositoryMySQL) GetFriendRequests(id int, sented bool) (friends []domain.Friend, err error) {

	var query string
	if sented {
		query = "SELECT id, user_a, user_b, sent FROM FriendRequest WHERE user_a = (?)"
	} else {
		query = "SELECT id, user_a, user_b, sent FROM FriendRequest WHERE user_b = (?)"
	}

	rows, err := r.db.Query(query, id)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Iterate Rows
	for rows.Next() {
		var friend domain.Friend
		var sentBytes []uint8
		err = rows.Scan(&friend.Id, &friend.UserA, &friend.UserB, &sentBytes)
		friend.StartFriendship = string(sentBytes)
		if err != nil {
			return
		}
		friends = append(friends, friend)
	}

	// Check Error on Rows
	if err = rows.Err(); err != nil {
		return
	}

	return
}

// GetFriends gets the friends
func (r *RepositoryMySQL) GetFriends(id int) (friends []domain.Friend, err error) {
	query := "SELECT id, user_a, user_b, start_friendship FROM Friend WHERE user_a = (?) OR user_b = (?)"

	rows, err := r.db.Query(query, id, id)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// iterate rows
	for rows.Next() {
		var friend domain.Friend
		var sentBytes []uint8
		err = rows.Scan(&friend.Id, &friend.UserA, &friend.UserB, &sentBytes)
		friend.StartFriendship = string(sentBytes)
		if err != nil {
			return
		}
		friends = append(friends, friend)
	}

	// Check Error on Rows
	if err = rows.Err(); err != nil {
		return
	}

	return
}
