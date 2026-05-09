package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
		respondError(w, http.StatusBadRequest, err)
		return
	}

	note, err := h.service.CreateNote(r.Context(), &input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusCreated, note)

}

func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	note, err := h.service.GetNoteByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusOK, note)

}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {
	notes, err := h.service.ListNotes(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusOK, notes)

}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	var input models.UpdateNoteInput

	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if err = h.service.UpdateNote(r.Context(), &input, id); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusOK, map[string]string{
		"message": "note updated.",
	})

}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if err = h.service.DeleteNote(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusOK, map[string]string{
		"message": "note deleted.",
	})
}
