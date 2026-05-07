package handlers

import (
	"encoding/json"
	"net/http"
)

type errorReponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Data any `json:"data"`
}

func responseJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, status int, err error) {
	responseJSON(w, status, errorReponse{Error: err.Error()})
}

func respondSuccess(w http.ResponseWriter, status int, data any) {
	responseJSON(w, status, successResponse{Data: data})
}
