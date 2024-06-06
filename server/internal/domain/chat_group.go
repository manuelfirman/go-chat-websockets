package domain

import "encoding/json"

// ChatGroup struct to represent a chat group
type ChatGroup struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sent        string `json:"sent"`
	Owner       int    `json:"owner"`
}

// BodyToChatGroup converts a byte array to a ChatGroup struct
func BodyToChatGroup(body []byte) *ChatGroup {
	if len(body) == 0 {
		return nil
	}

	var chatGroup ChatGroup
	err := json.Unmarshal(body, &chatGroup)
	if err != nil {
		return nil
	}

	return &chatGroup
}
