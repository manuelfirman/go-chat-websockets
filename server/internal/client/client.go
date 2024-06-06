package client

import (
	"github.com/gorilla/websocket"
)

// ClientDefault is the struct that implements the interface Client
type ClientDefault struct {
	Conn   *websocket.Conn
	Id     int
	Groups []int
}

// NewClientDefault creates a new ClientDefault
func NewClientDefault(conn *websocket.Conn, id int, groups []int) *ClientDefault {
	return &ClientDefault{Conn: conn, Id: id, Groups: groups}
}

// BelongToThisGroup checks if the client belongs to the group
func (c ClientDefault) BelongToThisGroup(id int) bool {
	for _, groupId := range c.Groups {
		if id == groupId {
			return true
		}
	}
	return false
}
