package main

import (
	"encoding/json"
	"net/http"
)

type APIErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func respondAPIError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIErrorResponse{
		Code:    code,
		Message: message,
	})
}
