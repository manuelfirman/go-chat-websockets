package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/manuelfirman/go-chat-websockets/internal/client"
	"github.com/manuelfirman/go-chat-websockets/internal/domain"
	"github.com/manuelfirman/go-chat-websockets/pkg/crypto"
)

// WebSocket is a struct that contains the configuration for the Web Socket.
type WebSocketDefault struct {
	// Upgrader is a struct that contains the configuration for the Web Socket.
	Upgrader websocket.Upgrader
	// Clients is a map that contains all the clients connected to the server.
	Clients map[*client.ClientDefault]bool
	// Broadcast is a channel that will be used to send messages to all the clients.
	Broadcast chan *domain.Message
	// rp is the repository that will be used to interact with the database.
	rp Repository
}

// NewWebSocket creates a new WebSocket struct.
func NewWebSocket(r Repository) *WebSocketDefault {
	return &WebSocketDefault{
		Clients:   make(map[*client.ClientDefault]bool),
		Broadcast: make(chan *domain.Message),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins. (if you want to restrict the origins, you can do it here.)
				return true
			},
			// Set the read buffer size
			ReadBufferSize: 1024,
			// Set the write buffer size
			WriteBufferSize: 1024,
		},
		rp: r,
	}
}

// HandleWebSocket handles the Web Socket connection.
func (wd *WebSocketDefault) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection
	conn, err := wd.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		println("Error al hacer el Upgrade.")
		// http.Error(w, errors.UPGRADE_FAILED+" "+err.Error(), http.StatusBadRequest)
		http.Error(w, "upgrade failed"+err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new client
	var groups []int
	client := client.NewClientDefault(conn, -1, groups)
	wd.Clients[client] = true

	println("Client [" + strconv.Itoa(client.Id) + "] Connected")

	// Close the connection and delete the client when the function ends
	defer func() {
		println("Closing the connection and deleting the client. " + "Client [" + strconv.Itoa(client.Id) + "]")
		conn.Close()
		delete(wd.Clients, client)
	}()

	// Handle the messages
	for {
		// Read the message from the client
		_, p, err := conn.ReadMessage()
		if err != nil {
			println(err)
			return
		}
		// Convert the bytes to a message
		message := domain.BodyToMessage(p)
		if message == nil {
			panic("Message is nil.")
		}

		// Check if the message is a JWT
		if message.IsJWT {
			id := crypto.GetIdFromJWT(message.Message)
			if id == nil {
				println("JWT INVALID")
				continue
			}
			client.Id = *id
			groups, err := wd.rp.GetGroupsOfUser(*id)
			if err != nil {
				println("Error getting the groups of the user.")
				continue
			}
			client.Groups = groups
			continue
		}
		// Check if the client has a JWT
		if client.Id == -1 {
			println("Client can't send messages, first have to provide the JWT.")
			continue
		}
		// Insert the message into the database
		if message.IsGroup {
			err = wd.rp.InsertGroupMessage(client.Id, message)
		} else {
			err = wd.rp.InsertUserMessage(client.Id, message)
		}

		if err != nil {
			println("Error sending the message.")
		}

		// Send the message to the Broadcast channel
		wd.Broadcast <- message
	}
}

// HandleMessages handles the messages that are sent to the clients. It will send the message to the correct client.
func (wd *WebSocketDefault) HandleMessages() {
	for {
		println("[HandleMessages] Looking for messages.")
		// Get the message from the Broadcast channel
		message := <-wd.Broadcast
		println("[HandleMessages] Message getted.")

		// Iterate over the clients
		for client := range wd.Clients {
			println("[HandleMessages] Looking for the correct client.")
			// Check if the client has to receive the message
			if wd.HaveToReceiveThisMessage(message, client) {
				println("[HandleMessages] Sending message to client.")
				// Send the message to the client
				wd.SendMessage(client, message)
				continue
			}
		}
	}
}

// HaveToReceiveThisMessage checks if the client has to receive the message. It will return true if the client has to receive the message.
func (wd *WebSocketDefault) HaveToReceiveThisMessage(message *domain.Message, client *client.ClientDefault) bool {
	// If the message is a group message, the client has to belong to the group and the client can't be the sender.
	if message.IsGroup {
		return client.BelongToThisGroup(message.ToId) && client.Id != message.UserId
	} else {
		return message.ToId == client.Id
	}
}

// SendMessage sends a message to the client. It will send the message in bytes.
func (wd *WebSocketDefault) SendMessage(client *client.ClientDefault, message *domain.Message) {
	// Convert the message to bytes
	b, err := json.Marshal(message)
	if err != nil {
		log.Println("Error wrapping the message to bytes. " + err.Error())
		client.Conn.Close()
		delete(wd.Clients, client)
	}
	// Send the message to the client
	err = client.Conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Println("Error writting the message into the Web Socket. ", err.Error())
		client.Conn.Close()
		delete(wd.Clients, client)
	}
}
