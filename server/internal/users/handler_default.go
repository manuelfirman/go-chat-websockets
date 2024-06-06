package users

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/manuelfirman/go-chat-websockets/internal/domain"
	"github.com/manuelfirman/go-chat-websockets/pkg/crypto"
	"github.com/manuelfirman/go-chat-websockets/pkg/response"
)

// HandlerDefault handles the user
type HandlerDefault struct {
	rp Repository
}

// NewHandlerDefault creates a new handler
func NewHandlerDefault(rp Repository) *HandlerDefault {
	return &HandlerDefault{
		rp: rp,
	}
}

// GetAllUsers gets all users from the database
func (h *HandlerDefault) GetAllUsers(w http.ResponseWriter) (users []domain.User) {
	users, err := h.rp.GetAllUsers()
	if err != nil {
		http.Error(w, "error getting all users", http.StatusInternalServerError)
		return
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

	return
}

// CreateUser creates a new user in the database
func (h *HandlerDefault) CreateUser(w http.ResponseWriter, r *http.Request) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	user := domain.BodyToUser(body)
	if user == nil {
		http.Error(w, "error unmarshalling user", http.StatusBadRequest)
		return false
	}

	newUserId, err := h.rp.CreateUser(user)
	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return false
	}

	tokenString := crypto.GenerateJWT(int(newUserId))

	resData := response.ResponseData{
		Message: *tokenString,
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

// GetUserFromId gets a user from the database given an id
func (h *HandlerDefault) GetUserFromId(id int, w http.ResponseWriter) (user *domain.User) {
	user, err := h.rp.GetUserFromId(id)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}

	return
}

// GetUserFromEmail gets a user from the database given an email
func (h *HandlerDefault) GetUserFromEmail(email string, w http.ResponseWriter) *domain.User {
	user, err := h.rp.GetUserFromEmail(email)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return nil
	}

	return user
}

// Login logs in a user
func (h *HandlerDefault) Login(w http.ResponseWriter, r *http.Request) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	user := domain.BodyToUser(body)
	if user == nil {
		http.Error(w, "error unmarshalling user", http.StatusBadRequest)
		return false
	}

	// Get the real user
	realUser, err := h.rp.GetUserFromEmail(user.Email)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return false
	}

	if realUser == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user not found"))
		return false
	}

	if realUser.Password != user.Password {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error password incorrect"))
		return false
	}

	realUser.Password = ""

	tokenString := crypto.GenerateJWT(realUser.Id)

	resData := response.ResponseLogin{
		Jwt:      *tokenString,
		RealUser: *realUser,
	}

	resJSON := response.GetResponseLoginJSON(resData)

	if resJSON == nil {
		http.Error(w, "error converting the response data to json", http.StatusInternalServerError)
		return false
	}

	w.WriteHeader(http.StatusOK)
	w.Write(*resJSON)

	return true
}

// GetUserMessages retrieves the messages between two users
func (h *HandlerDefault) GetParticularUser(idString string, w http.ResponseWriter, r *http.Request) bool {
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "error converting id to int", http.StatusBadRequest)
		return false
	}

	user, err := h.rp.GetUserFromId(id)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return false
	}
	if user == nil {
		http.Error(w, "error user not found", http.StatusBadRequest)
		return false
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*user)

	return true
}

// HandleUser handles the /user endpoint
func (h *HandlerDefault) HandleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.CreateUser(w, r)
		return
	}

	if r.Method == http.MethodGet {
		queryParams := r.URL.Query()
		id := queryParams.Get("id")
		if id == "" {
			h.GetAllUsers(w)
		} else {
			h.GetParticularUser(id, w, r)
		}

		return
	}

	http.Error(w, "Method not allowed to /user", http.StatusMethodNotAllowed)
}
