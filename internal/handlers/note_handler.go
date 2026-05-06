package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kevalsabhani/keeper/internal/models"
	"github.com/kevalsabhani/keeper/internal/services"
)

type NoteHandler struct {
	service *services.NoteService
}

func NewNoteHandler(service *services.NoteService) *NoteHandler {
	return &NoteHandler{
		service: service,
	}
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateNoteInput

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	note, err := h.service.CreateNote(r.Context(), &input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"data": note,
	})

}

func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {

}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {

}
