package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kevalsabhani/keeper/internal/models"
	"github.com/kevalsabhani/keeper/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUserInput

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.CreateUser(r.Context(), &input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusCreated, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusOK, user)

}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusOK, users)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	var input models.UpdateUserInput

	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if err = h.service.UpdateUser(r.Context(), &input, id); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusOK, map[string]string{
		"message": "user updated.",
	})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if err = h.service.DeleteUser(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusOK, map[string]string{
		"message": "user deleted.",
	})
}
