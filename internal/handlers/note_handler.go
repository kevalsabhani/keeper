package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	errpkg "github.com/kevalsabhani/keeper/internal/errors"
	"github.com/kevalsabhani/keeper/internal/models"
	"github.com/kevalsabhani/keeper/internal/response"
	"github.com/kevalsabhani/keeper/internal/services"
	"go.uber.org/zap"
)

// NoteHandler handles HTTP requests for note-related endpoints.
type NoteHandler struct {
	service *services.NoteService
	log     *zap.Logger
}

// NewNoteHandler creates a NoteHandler with the given service dependency.
func NewNoteHandler(service *services.NoteService, log *zap.Logger) *NoteHandler {
	return &NoteHandler{
		service: service,
		log:     log,
	}
}

// Create decodes a JSON body and creates a new note. Returns 201 on success.
func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateNoteInput

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("failed to decode create note request body", zap.Error(err))
		response.Error(w, errpkg.ErrInvalidInput)
		return
	}

	note, err := h.service.CreateNote(r.Context(), &input)
	if err != nil {
		h.log.Error("failed to create note", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusCreated, note, nil)
}

// GetByID fetches a single note by its URL path ID. Returns 404 if not found.
func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("invalid note id in path", zap.String("raw_id", chi.URLParam(r, "id")))
		response.Error(w, errpkg.ErrInvalidInput)
		return
	}

	note, err := h.service.GetNoteByID(r.Context(), id)
	if err != nil {
		h.log.Error("failed to get note", zap.Int("id", id), zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, note, &response.Meta{
		CurrentPage: 1,
		TotalPages:  1,
		TotalCount:  1,
	})
}

// List returns a paginated list of notes. Accepts optional `page` and `limit`
// query params; defaults to page=1, limit=20 (capped at 100).
func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {

	var (
		page    int
		limit   int
		listErr error
	)

	if r.URL.Query().Has("page") {
		page, listErr = strconv.Atoi(r.URL.Query().Get("page"))
		if listErr != nil {
			h.log.Warn("invalid page query parameter", zap.String("value", r.URL.Query().Get("page")))
			response.Error(w, errpkg.ErrInvalidInput)
			return
		}
	}

	if page < 1 {
		page = 1
	}

	if r.URL.Query().Has("limit") {
		limit, listErr = strconv.Atoi(r.URL.Query().Get("limit"))
		if listErr != nil {
			h.log.Warn("invalid limit query parameter", zap.String("value", r.URL.Query().Get("limit")))
			response.Error(w, errpkg.ErrInvalidInput)
			return
		}
	}

	if limit < 1 || limit > 100 {
		limit = 20
	}

	notes, meta, listErr := h.service.ListNotes(r.Context(), page, limit)
	if listErr != nil {
		h.log.Error("failed to list notes", zap.Int("page", page), zap.Int("limit", limit), zap.Error(listErr))
		response.Error(w, listErr)
		return
	}

	response.Success(w, http.StatusOK, notes, meta)
}

// Update applies partial changes to an existing note identified by URL path ID.
func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("invalid note id in path", zap.String("raw_id", chi.URLParam(r, "id")))
		response.Error(w, errpkg.ErrInvalidInput)
		return
	}

	var input models.UpdateNoteInput

	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("failed to decode update note request body", zap.Int("id", id), zap.Error(err))
		response.Error(w, errpkg.ErrInvalidInput)
		return
	}

	if err = h.service.UpdateNote(r.Context(), &input, id); err != nil {
		h.log.Error("failed to update note", zap.Int("id", id), zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, nil, nil)
}

// Delete removes a note by its URL path ID. Returns 404 if the note does not exist.
func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("invalid note id in path", zap.String("raw_id", chi.URLParam(r, "id")))
		response.Error(w, errpkg.ErrInvalidInput)
		return
	}

	if err = h.service.DeleteNote(r.Context(), id); err != nil {
		h.log.Error("failed to delete note", zap.Int("id", id), zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, nil, nil)
}
