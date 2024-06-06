package chatgroup

import (
	"github.com/manuelfirman/go-chat-websockets/internal/domain"
)

type Repository interface {
	// CreateGroup creates a new group
	CreateGroup(id int, name, description string) (err error)
	// GetGroups gets the groups
	GetGroups(id int) (groups []domain.ChatGroup, err error)
	// GetGroupUsers gets the users of a group
	GetGroupUsers(groupId int) (users []domain.User, err error)
	// GetGroupOwner gets the owner of a group
	GetGroupOwner(groupId int) (owner *int, err error)
}
