package response

import (
	"encoding/json"
)

// ResponseData struct to represent a response message
type ResponseData struct {
	Message string `json:"message"`
}

// ResponseLogin struct to represent a response message
type ResponseLogin struct {
	Jwt      string      `json:"jwt"`
	RealUser interface{} `json:"user"`
}

// GetResponseDataJSON converts a ResponseData struct to a byte array
func GetResponseDataJSON(res ResponseData) *[]byte {
	resJSON, err := json.Marshal(res)
	if err != nil {
		return nil
	}

	return &resJSON
}

// GetResponseLoginJSON converts a ResponseLogin struct to a byte array
func GetResponseLoginJSON(res ResponseLogin) *[]byte {
	resJSON, err := json.Marshal(res)
	if err != nil {
		return nil
	}

	return &resJSON
}
