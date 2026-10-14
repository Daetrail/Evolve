package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func SendSuccessWithNoData(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response{
		Success: true,
	})
	if err != nil {
		log.Println("failed to encode response: ", err)
	}
}

func SendSuccess(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response{
		Success: true,
		Data:    data,
	})
	if err != nil {
		log.Println("failed to encode response: ", err)
	}
}

func SendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response{
		Success: false,
		Message: message,
	})
	if err != nil {
		log.Println("failed to encode response: ", err)
	}
}
