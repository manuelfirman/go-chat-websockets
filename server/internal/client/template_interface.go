package client

import "net/http"

// Template is the interface for the client
type Template interface {
	// ServeHTTP serves the http request
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
