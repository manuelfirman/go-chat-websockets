package domain

import "encoding/json"

// Friend is the domain model for a friend
type Friend struct {
	Id              int    `json:"id"`
	UserA           int    `json:"user_a"`
	UserB           int    `json:"user_b"`
	StartFriendship string `json:"start_friendship"`
}

// FriendToBody converts a friend to a byte array
func BodyToFriend(body []byte) *Friend {
	if len(body) == 0 {
		return nil
	}

	var friend Friend
	err := json.Unmarshal(body, &friend)
	if err != nil {
		return nil
	}

	return &friend
}
