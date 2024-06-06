package client

// Client is the interface for the client
type Client interface {
	// BelongToThisGroup checks if the client belongs to this group
	BelongToThisGroup(id int) bool
}
