package application

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/manuelfirman/go-chat-websockets/internal/chatgroup"
	"github.com/manuelfirman/go-chat-websockets/internal/chatgrouprequest"
	"github.com/manuelfirman/go-chat-websockets/internal/client"
	"github.com/manuelfirman/go-chat-websockets/internal/friend"
	"github.com/manuelfirman/go-chat-websockets/internal/health"
	"github.com/manuelfirman/go-chat-websockets/internal/message"
	"github.com/manuelfirman/go-chat-websockets/internal/users"
	"github.com/manuelfirman/go-chat-websockets/internal/websocket"
)

// ApplicationDefault is the struct that implements the interface Application
type ApplicationDefault struct {
	Config
	// websocket
	WebSocket websocket.WebSocket
}

// Config is the configuration for the application
type Config struct {
	Addr     string
	MySQLDSN string
}

// NewApplicationDefault creates a new ApplicationDefault
func NewApplicationDefault(config Config) *ApplicationDefault {
	defaultCfg := Config{
		Addr:     ":8080",
		MySQLDSN: "",
	}
	if config.Addr == "" {
		config.Addr = defaultCfg.Addr
	}
	if config.MySQLDSN == "" {
		config.MySQLDSN = defaultCfg.MySQLDSN
	}

	return &ApplicationDefault{Config: config}
}

// SetUp sets up the application
func (a *ApplicationDefault) SetUp() {
	// Open connection to Database
	database, err := sql.Open("mysql", a.MySQLDSN)
	if err != nil {
		panic(err.Error())
	}
	// defer database.Close()

	// Check success connection to Database
	err = database.Ping()
	if err != nil {
		panic(err.Error())
	}
	println("Connected to the database")

	// Build dependencies
	http.HandleFunc("/", client.NewTemplateHandler().ServeHTTP)
	buildUserDependencies(database)
	buildChatGroupDependencies(database)
	buildChatGroupRequestDependencies(database)
	buildMessageDependencies(database)
	buildFriendDependencies(database)

	a.buildWebSocket(database)

	health := health.NewHealthChecker(database)
	http.HandleFunc("/health", health.Ping)

}

// Run runs the application
func (a *ApplicationDefault) Run() {
	go a.WebSocket.HandleMessages()

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)
	println("Server listening on port" + a.Addr + "...")
	if err := http.ListenAndServe(a.Addr, corsMiddleware(http.DefaultServeMux)); err != nil {
		panic(err)
	}
}

// buildWebSocket builds the websocket dependencies
func (a *ApplicationDefault) buildWebSocket(db *sql.DB) {
	rp := websocket.NewRepositoryMySQL(db)
	a.WebSocket = websocket.NewWebSocket(rp)

	http.HandleFunc("/ws", a.WebSocket.HandleWebSocket)
}

// build user dependencies
func buildUserDependencies(db *sql.DB) {
	rp := users.NewRepositoryMySQL(db)
	h := users.NewHandlerDefault(rp)

	http.HandleFunc("/user", h.HandleUser)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		h.Login(w, r)
	})

}

// build chat group dependencies
func buildChatGroupDependencies(db *sql.DB) {
	rp := chatgroup.NewRepositoryMySQL(db)
	h := chatgroup.NewHandlerDefault(rp)

	http.HandleFunc("/group", h.HandleGroups)
	http.HandleFunc("/group-users", h.HandleGroupUsers)
}

// build chat group request dependencies
func buildChatGroupRequestDependencies(db *sql.DB) {
	rp := chatgrouprequest.NewRepositoryMySQL(db)
	h := chatgrouprequest.NewHandlerDefault(rp)

	http.HandleFunc("/accept-group", h.HandleAcceptGroupRequest)
	http.HandleFunc("/group-request", h.HandleGroupRequest)
}

// build message dependencies
func buildMessageDependencies(db *sql.DB) {
	rp := message.NewRepositoryMySQL(db)
	h := message.NewHandlerDefault(rp)

	http.HandleFunc("/user-message", h.HandleUserMessage)
	http.HandleFunc("/group-message", h.HandleGroupMessage)
}

// build friend dependencies
func buildFriendDependencies(db *sql.DB) {
	rp := friend.NewRepositoryMySQL(db)
	h := friend.NewHandlerDefault(rp)

	http.HandleFunc("/friend-request", h.HandleFriendRequest)
	http.HandleFunc("/accept-friend", h.HandleAcceptFriend)
	http.HandleFunc("/friends", h.HandleFriends)
}
