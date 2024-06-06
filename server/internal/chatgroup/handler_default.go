package chatgroup

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/manuelfirman/go-chat-websockets/internal/domain"
	"github.com/manuelfirman/go-chat-websockets/pkg/crypto"
	"github.com/manuelfirman/go-chat-websockets/pkg/response"
)

// HandlerDefault handles the chat group
type HandlerDefault struct {
	rp Repository
}

// NewHandlerDefault creates a new handler
func NewHandlerDefault(rp Repository) *HandlerDefault {
	return &HandlerDefault{
		rp: rp,
	}
}

// NewChatGroupHandler creates a new chat group handler
func (c *HandlerDefault) HandleGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c.GetGroups(w, r)
		return
	}

	if r.Method == http.MethodPost {
		c.CreateGroup(w, r)
		return
	}

	http.Error(w, "Method not allowed to /groups", http.StatusMethodNotAllowed)
}

// HandleGroupUsers handles the group users
func (c *HandlerDefault) HandleGroupUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c.GetGroupUsers(w, r)
		return
	}

	http.Error(w, "Method not allowed to /group-users", http.StatusMethodNotAllowed)
}

// CreateGroup creates a new group
func (c *HandlerDefault) CreateGroup(w http.ResponseWriter, r *http.Request) bool {
	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, "error JWT invalid", http.StatusBadRequest)
		return false
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading the body request", http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	chatGroup := domain.BodyToChatGroup(body)
	if chatGroup == nil {
		http.Error(w, "error unmarshalling chat group", http.StatusBadRequest)
		return false
	}

	err = c.rp.CreateGroup(*id, chatGroup.Name, chatGroup.Description)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error inserting the group into the database" + err.Error()))
		return false
	}

	resData := response.ResponseData{
		Message: "SUCCESSFUL POST REQUEST",
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

// GetGroups gets the groups
func (c *HandlerDefault) GetGroups(w http.ResponseWriter, r *http.Request) bool {
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

	chatGroups, err := c.rp.GetGroups(*id)
	if err != nil {
		http.Error(w, "error getting groups", http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chatGroups)

	return true
}

// GetGroupUsers gets the users of a group
func (c *HandlerDefault) GetGroupUsers(w http.ResponseWriter, r *http.Request) bool {
	// tokenString := crypto.GetJWTFromRequest(w, r)
	// if tokenString == nil {
	// 	http.Error(w, "error JWT not found", http.StatusBadRequest)
	// 	return false
	// }

	// id := crypto.GetIdFromJWT(*tokenString)
	// if id == nil {
	// 	http.Error(w, "error JWT invalid", http.StatusBadRequest)
	// 	return false
	// }

	queryParams := r.URL.Query()
	groupIdString := queryParams.Get("group_id")
	groupId, err := strconv.Atoi(groupIdString)
	if err != nil {
		http.Error(w, "error invalid id", http.StatusBadRequest)
		return false
	}

	users, err := c.rp.GetGroupUsers(groupId)
	if err != nil {
		http.Error(w, "error getting group users", http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

	return true
}

// GetGroupOwner gets the owner of a group
func (c *HandlerDefault) GetGroupOwner(id int) (owner *int) {
	owner, err := c.rp.GetGroupOwner(id)
	if err != nil {
		panic(err.Error())
	}

	return
}

// // GetGropusOfUser gets the groups of a user
// func (c *HandlerDefault) GetGroupsOfUser(id int) (groups []int) {
// 	groups, err := c.rp.GetGroupsOfUser(id)
// 	if err != nil {
// 		return
// 	}

// 	return
// }
