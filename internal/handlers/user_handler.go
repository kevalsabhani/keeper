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

// UserHandler handles HTTP requests for user-related endpoints.
type UserHandler struct {
	service *services.UserService
	log     *zap.Logger
}

// NewUserHandler creates a UserHandler with the given service dependency.
func NewUserHandler(service *services.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		log:     log,
	}
}

// Create decodes a JSON body and creates a new user. Returns 201 on success.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUserInput

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("failed to decode create user request body", zap.Error(err))
		response.Error(w, errpkg.ErrInvalidInput)
		return
	}

	user, err := h.service.CreateUser(r.Context(), &input)
	if err != nil {
		h.log.Error("failed to create user", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusCreated, user, nil)
}

// GetByID fetches a single user by its URL path ID. Returns 404 if not found.
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("invalid user id in path", zap.String("raw_id", chi.URLParam(r, "id")))
		response.Error(w, errpkg.ErrInvalidInput)
		return
	}

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		h.log.Error("failed to get user", zap.Int("id", id), zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, user, &response.Meta{
		CurrentPage: 1,
		TotalPages:  1,
		TotalCount:  1,
	})
}

// List returns a paginated list of users. Accepts optional `page` and `limit`
// query params; defaults to page=1, limit=20 (capped at 100).
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {

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

	users, meta, listErr := h.service.ListUsers(r.Context(), page, limit)
	if listErr != nil {
		h.log.Error("failed to list users", zap.Int("page", page), zap.Int("limit", limit), zap.Error(listErr))
		response.Error(w, listErr)
		return
	}

	response.Success(w, http.StatusOK, users, meta)
}

// Update applies partial changes to an existing user identified by URL path ID.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("invalid user id in path", zap.String("raw_id", chi.URLParam(r, "id")))
		response.Error(w, errpkg.ErrInvalidInput)
		return
	}

	var input models.UpdateUserInput

	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Warn("failed to decode update user request body", zap.Int("id", id), zap.Error(err))
		response.Error(w, errpkg.ErrInvalidInput)
		return
	}

	if err = h.service.UpdateUser(r.Context(), &input, id); err != nil {
		h.log.Error("failed to update user", zap.Int("id", id), zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, nil, nil)
}

// Delete removes a user by its URL path ID. Returns 404 if the user does not exist.
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.log.Warn("invalid user id in path", zap.String("raw_id", chi.URLParam(r, "id")))
		response.Error(w, errpkg.ErrInvalidInput)
		return
	}

	if err = h.service.DeleteUser(r.Context(), id); err != nil {
		h.log.Error("failed to delete user", zap.Int("id", id), zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, nil, nil)
}
