package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
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
	if errors.Is(err, pgx.ErrNoRows) {
		responseJSON(w, http.StatusNotFound, errorReponse{Error: err.Error()})
		return
	}
	responseJSON(w, status, errorReponse{Error: err.Error()})
}

func respondSuccess(w http.ResponseWriter, status int, data any) {
	responseJSON(w, status, successResponse{Data: data})
}
