package chatgrouprequest

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/manuelfirman/go-chat-websockets/internal/domain"
	"github.com/manuelfirman/go-chat-websockets/pkg/crypto"
	"github.com/manuelfirman/go-chat-websockets/pkg/response"
)

type HandlerDefault struct {
	rp Repository
}

// NewHandlerDefault creates a new handler
func NewHandlerDefault(rp Repository) *HandlerDefault {
	return &HandlerDefault{
		rp: rp,
	}
}

// SendGroupRequest sends a group request given the group_id
func (c *HandlerDefault) SendGroupRequest(w http.ResponseWriter, r *http.Request) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading the body request", http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	request := domain.BodyToChatGroupRequest(body)
	if request == nil {
		http.Error(w, "error unmarshalling chat group request", http.StatusBadRequest)
		return false
	}

	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, "error JWT invalid", http.StatusBadRequest)
		return false
	}

	err = c.rp.InsertChatGroupRequest(*id, request.GroupId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error inserting the group request: " + err.Error()))
		return false
	}

	resData := response.ResponseData{
		Message: "SUCCESFULL POST REQUEST",
	}
	resJSON := response.GetResponseDataJSON(resData)

	if resJSON == nil {
		http.Error(w, "error converting the response data to json", http.StatusInternalServerError)
		return false
	}

	w.WriteHeader(http.StatusOK)
	w.Write(*resJSON)

	return true
}

// GetGroupOwner returns the owner of the group given the group_id
func (c *HandlerDefault) GetGroupsRequests(w http.ResponseWriter, r *http.Request) bool {
	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, "error JWT invalid", http.StatusBadRequest)
		return false
	}

	group_id, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	if err != nil {
		http.Error(w, "error getting the group_id", http.StatusBadRequest)
		return false
	}

	sented, err := strconv.ParseBool(r.URL.Query().Get("sented"))
	if err != nil {
		http.Error(w, "error getting the sented", http.StatusBadRequest)
		return false
	}

	requests, err := c.rp.GetChatGroupRequests(*id, group_id, sented)
	if err != nil {
		http.Error(w, "error getting the group requests", http.StatusBadRequest)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requests)

	return true
}

// AcceptGroupRequest accepts a group request given the group_id and the user_id
func (c *HandlerDefault) AcceptGroupRequest(w http.ResponseWriter, r *http.Request) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	request := domain.BodyToChatGroupRequest(body)
	if request == nil {
		http.Error(w, "error unmarshalling chat group request", http.StatusBadRequest)
		return false
	}

	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, "error JWT invalid", http.StatusBadRequest)
		return false
	}

	owner, err := c.rp.GetGroupOwner(request.GroupId)
	if err != nil {
		http.Error(w, "error getting the group owner", http.StatusBadRequest)
		return false
	}
	if *id != *owner {
		http.Error(w, "error user is not the owner of the group", http.StatusBadRequest)
		return false
	}

	err = c.rp.AcceptGroupRequest(request.UserId, request.GroupId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error inserting into database" + err.Error()))
		return false
	}

	resData := response.ResponseData{
		Message: "SUCCESFULL POST REQUEST",
	}
	resJSON := response.GetResponseDataJSON(resData)

	if resJSON == nil {
		http.Error(w, "error converting the response data to json", http.StatusInternalServerError)
		return false
	}

	w.WriteHeader(http.StatusOK)
	w.Write(*resJSON)

	return true
}

// HandleGroupRequest handles the group request
func (c *HandlerDefault) HandleGroupRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.SendGroupRequest(w, r)
		return
	}

	if r.Method == http.MethodGet {
		c.GetGroupsRequests(w, r)
		return
	}

	http.Error(w, "Method not allowed to /group-request", http.StatusMethodNotAllowed)
}

// HandleAcceptGroupRequest handles the accept group request
func (c *HandlerDefault) HandleAcceptGroupRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.AcceptGroupRequest(w, r)
		return
	}

	http.Error(w, "Method not allowed to /accept-request", http.StatusMethodNotAllowed)
}
