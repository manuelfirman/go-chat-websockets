package message

import "github.com/manuelfirman/go-chat-websockets/internal/domain"

// Repository is the repository for the user
type Repository interface {
	// GetGroupMessages gets the messages between two users
	GetGroupMessages(groupId int) (messages []domain.Message, err error)
	// GetUserMessages gets the messages between two users
	GetUserMessages(userId int, friendId int) (messages []domain.Message, err error)
}
