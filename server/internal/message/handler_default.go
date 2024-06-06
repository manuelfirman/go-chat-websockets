package message

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/manuelfirman/go-chat-websockets/pkg/crypto"
)

// HandlerDefault handles the messages
type HandlerDefault struct {
	rp Repository
}

// NewHandlerDefault creates a new handler
func NewHandlerDefault(rp Repository) *HandlerDefault {
	return &HandlerDefault{
		rp: rp,
	}
}

// HandleUserMessage handles the user messages
func (h *HandlerDefault) HandleUserMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetUserMessage(w, r)
		return
	}

	http.Error(w, "Method not allowed to /user-message", http.StatusMethodNotAllowed)
}

// HandleGroupMessage handles the group messages
func (h *HandlerDefault) HandleGroupMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetGroupMessage(w, r)
		return
	}

	http.Error(w, "Method not allowed to /group-message", http.StatusMethodNotAllowed)
}

// GetUserMessage gets the messages between two users
func (h *HandlerDefault) GetGroupMessage(w http.ResponseWriter, r *http.Request) bool {
	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		http.Error(w, "error JWT not found", http.StatusBadRequest)
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, "error JWT invalid", http.StatusBadRequest)
		return false
	}

	queryParams := r.URL.Query()
	groupIdStr := queryParams.Get("group_id")
	if groupIdStr == "" {
		http.Error(w, "error friend_id not found", http.StatusBadRequest)
		return false
	}
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		http.Error(w, "error friend_id not valid", http.StatusBadRequest)
		return false
	}

	messages, err := h.rp.GetGroupMessages(groupId)
	if err != nil {
		http.Error(w, "error getting messages", http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)

	return true
}

// GetUserMessage retrieves the messages between two users
func (h *HandlerDefault) GetUserMessage(w http.ResponseWriter, r *http.Request) bool {
	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		http.Error(w, "error JWT not found", http.StatusBadRequest)
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, "error JWT invalid", http.StatusBadRequest)
		return false
	}

	queryParams := r.URL.Query()
	friendIdStr := queryParams.Get("friend_id")
	if friendIdStr == "" {
		http.Error(w, "error friend_id not found", http.StatusBadRequest)
		return false
	}
	friendId, err := strconv.Atoi(friendIdStr)
	if err != nil {
		http.Error(w, "error friend_id not valid", http.StatusBadRequest)
		return false
	}

	messages, err := h.rp.GetUserMessages(*id, friendId)
	if err != nil {
		http.Error(w, "error getting messages", http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)

	return true
}
