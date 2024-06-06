package domain

import "encoding/json"

// Message is the domain model for a message
type Message struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id"`
	IsJWT   bool   `json:"is_jwt"`
	IsGroup bool   `json:"is_group"`
	ToId    int    `json:"to_id"`
	Message string `json:"message"`
	Sent    string `json:"sent"`
}

// MessageToBody converts a message to a byte array
func BodyToMessage(body []byte) *Message {
	if len(body) == 0 {
		return nil
	}

	var message Message
	err := json.Unmarshal(body, &message)
	if err != nil {
		println("Error: ", err.Error())
		return nil
	}

	return &message
}
