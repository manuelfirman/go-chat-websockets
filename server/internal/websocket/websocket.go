package websocket

import (
	"net/http"

	"github.com/manuelfirman/go-chat-websockets/internal/client"
	"github.com/manuelfirman/go-chat-websockets/internal/domain"
)

type WebSocket interface {
	// HandleWebSocket handles the websocket connection
	HandleWebSocket(w http.ResponseWriter, r *http.Request)
	// HaveToReceiveThisMessage checks if the message is for the client
	HaveToReceiveThisMessage(message *domain.Message, client *client.ClientDefault) bool
	// HandleMessages handles the messages
	HandleMessages()
	// SendMessage sends a message to the client
	SendMessage(client *client.ClientDefault, message *domain.Message)
}
