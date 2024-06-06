package chatgrouprequest

import "github.com/manuelfirman/go-chat-websockets/internal/domain"

type Repository interface {
	// InsertChatGroupRequest inserts a new chat group request
	InsertChatGroupRequest(id int, groupId int) (err error)
	// GetChatGroupRequests gets the chat group requests
	GetChatGroupRequests(id, groupId int, sented bool) (requests []domain.ChatGroupRequest, err error)
	// GetGroupOwner gets the owner of a group
	GetGroupOwner(groupId int) (owner *int, err error)
	// AcceptGroupRequest accepts a group request
	AcceptGroupRequest(userId int, groupId int) (err error)
}
