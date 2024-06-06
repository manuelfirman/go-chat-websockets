package websocket

import "github.com/manuelfirman/go-chat-websockets/internal/domain"

type Repository interface {
	// InsertGroupMessage inserts a message in the group
	InsertGroupMessage(id int, message *domain.Message) (err error)
	// InsertUserMessage inserts a message in the user
	InsertUserMessage(id int, message *domain.Message) (err error)
	// GetGroupsOfUser gets the groups of a user
	GetGroupsOfUser(userId int) (groups []int, err error)
}
