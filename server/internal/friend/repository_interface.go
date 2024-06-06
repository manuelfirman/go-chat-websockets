package friend

import "github.com/manuelfirman/go-chat-websockets/internal/domain"

type Repository interface {
	// InsertFriendRequest inserts a new friend request in the database
	InsertFriendRequest(id int, friendId int) (err error)
	// AcceptFriendRequest accepts a friend request in the database
	AcceptFriendRequest(id int, friendId int) (rowsAffected int, err error)
	// GetFriendRequests gets the friend requests
	GetFriendRequests(id int, sented bool) (friends []domain.Friend, err error)
	// GetFriends gets the friends
	GetFriends(id int) (friends []domain.Friend, err error)
}
