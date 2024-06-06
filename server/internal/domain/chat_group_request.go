package domain

import "encoding/json"

// ChatGroupRequest struct to represent a chat group request
type ChatGroupRequest struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id"`
	GroupId int    `json:"group_id"`
	Sent    string `json:"sent"`
}

// BodyToChatGroupRequest converts a byte array to a ChatGroupRequest struct
func BodyToChatGroupRequest(body []byte) *ChatGroupRequest {
	if len(body) == 0 {
		return nil
	}

	var chatGroup ChatGroupRequest
	err := json.Unmarshal(body, &chatGroup)
	if err != nil {
		return nil
	}

	return &chatGroup
}
