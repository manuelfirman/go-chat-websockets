package friend

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/manuelfirman/go-chat-websockets/internal/domain"
	"github.com/manuelfirman/go-chat-websockets/pkg/crypto"
	"github.com/manuelfirman/go-chat-websockets/pkg/response"
)

// HandlerDefault handles the friend
type HandlerDefault struct {
	rp Repository
}

// NewHandlerDefault creates a new handler
func NewHandlerDefault(rp Repository) *HandlerDefault {
	return &HandlerDefault{
		rp: rp,
	}
}

// HandlerFriendRequest handles the friend request
func (h *HandlerDefault) HandleFriendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.SendFriendRequest(w, r)
		return
	}

	if r.Method == http.MethodGet {
		h.GetFriendRequests(w, r)
		return
	}

	http.Error(w, "Method not allowed to /friend-request", http.StatusMethodNotAllowed)
}

// HandlerFriends handles the friends
func (h *HandlerDefault) HandleFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetFriends(w, r)
		return
	}

	http.Error(w, "Method not allowed to /friends", http.StatusMethodNotAllowed)
}

// HandlerAcceptFriend handles the accept friend
func (h *HandlerDefault) HandleAcceptFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.AcceptFriendRequest(w, r)
		return
	}

	http.Error(w, "Method not allowed to /friends", http.StatusMethodNotAllowed)
}

// SendFriendRequest sends a friend request
func (h *HandlerDefault) SendFriendRequest(w http.ResponseWriter, r *http.Request) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading the body request", http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	friend := domain.BodyToFriend(body)
	if friend == nil {
		http.Error(w, "error unmarshalling friend", http.StatusBadRequest)
		return false
	}

	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		http.Error(w, "error getting the JWT token from the request", http.StatusBadRequest)
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, "error getting the id from the JWT token", http.StatusBadRequest)
		return false
	}
	err = h.rp.InsertFriendRequest(*id, friend.UserB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error inserting the friend request: " + err.Error()))
		return false
	}

	resData := response.ResponseData{
		Message: "SUCCESFULL POST REQUEST",
	}
	resJSON := response.GetResponseDataJSON(resData)

	if resJSON == nil {
		http.Error(w, "error converting the response data to JSON", http.StatusInternalServerError)
		return false
	}

	w.WriteHeader(http.StatusOK)
	w.Write(*resJSON)

	return true
}

// AcceptFriendRequest accepts a friend request
func (h *HandlerDefault) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading the body request", http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	friend := domain.BodyToFriend(body)
	if friend == nil {
		http.Error(w, "error unmarshalling friend", http.StatusBadRequest)
		return false
	}

	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		http.Error(w, "error getting the JWT token from the request", http.StatusBadRequest)
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, "error getting the id from the JWT token", http.StatusBadRequest)
		return false
	}

	// result, err := database.Exec(db.ACCEPT_FRINED_REQUEST, friend.UserA, *id)
	rowsAffected, err := h.rp.AcceptFriendRequest(*id, friend.UserA)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error accepting the friend request: " + err.Error()))
		return false
	}

	// Las filas afectadas dan siempre 1, sea cual sea el resultado.
	// Es decir, no podemos identificar cuando se acepta realmente un amigo
	// o cuando no se puede aceptar porque no existe la solicitud.
	println("Filas afectadas: ", rowsAffected)

	resData := response.ResponseData{
		Message: "SUCCESFULL POST REQUEST",
	}
	resJSON := response.GetResponseDataJSON(resData)

	if resJSON == nil {
		http.Error(w, "Error converting the response data to JSON. ", http.StatusInternalServerError)
		return false
	}

	w.WriteHeader(http.StatusOK)
	w.Write(*resJSON)

	return true
}

// GetFriendRequests gets the friend requests
func (h *HandlerDefault) GetFriendRequests(w http.ResponseWriter, r *http.Request) bool {
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

	sentedString := queryParams.Get("sented")
	sented, err := strconv.ParseBool(sentedString)
	if err != nil {
		http.Error(w, "error sented not valid", http.StatusBadRequest)
		return false
	}

	friends, err := h.rp.GetFriendRequests(*id, sented)
	if err != nil {
		http.Error(w, "error getting friend requests", http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(friends)

	return true
}

// GetFriends gets the friends
func (h *HandlerDefault) GetFriends(w http.ResponseWriter, r *http.Request) bool {
	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, "error JWT invalid", http.StatusBadRequest)
		return false
	}

	friends, err := h.rp.GetFriends(*id)
	if err != nil {
		println("error getting friends: ", err.Error())
		http.Error(w, "error getting friends", http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(friends)
	println("Friends: ", friends)

	return true
}
