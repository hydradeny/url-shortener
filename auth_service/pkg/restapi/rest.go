package restapi

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Body  interface{} `json:"body,omitempty"`
	Error string      `json:"error,omitempty"`
}

func RespJSON(w http.ResponseWriter, body interface{}) {
	w.Header().Add("Content-Type", "application/json")
	respJSON, _ := json.Marshal(&ApiResponse{
		Body: body,
	})
	w.Write(respJSON)
}

func RespJSONError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	respJSON, _ := json.Marshal(&ApiResponse{
		Error: err.Error(),
	})
	w.Write(respJSON)
}
