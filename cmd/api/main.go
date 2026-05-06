package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kevalsabhani/keeper/internal/configs"
	"github.com/kevalsabhani/keeper/internal/database"
	"github.com/kevalsabhani/keeper/internal/handlers"
	"github.com/kevalsabhani/keeper/internal/models"
	"github.com/kevalsabhani/keeper/internal/repositories"
	"github.com/kevalsabhani/keeper/internal/services"
)

// In-memory stores
var userStore = map[int]models.User{
	1: models.User{1, "Alice", "alice@ex.com"},
	2: models.User{2, "Bob", "bob@ex.com"},
}

var noteStore = map[int]models.Note{
	1: models.Note{1, 2, "keeper api key", "fjalkdfjldfjlsjf"},
	2: models.Note{2, 1, "aws secrets", "kadsfjdslfjlskdfjlksd"},
}

func main() {

	config := configs.Load()

	db, err := database.New(config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	userRepo := repositories.NewPostgresUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	noteRepo := repositories.NewPostgresNoteRepository(db)
	noteService := services.NewNoteService(noteRepo)
	noteHandler := handlers.NewNoteHandler(noteService)

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	// Routes
	r.Get("/", handleHome)
	r.Get("/health", handleHealth)

	// User routes
	userRoutes := func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Get("/", userHandler.List)
		r.Get("/{id}", userHandler.GetByID)
		r.Put("/{id}", handleUpdateUser)
		r.Delete("/{id}", handleDeleteUser)
	}

	// Note routes
	noteRoutes := func(r chi.Router) {
		r.Post("/", noteHandler.Create)
		r.Get("/", handleListNotes)
		r.Get("/{id}", handleGetNote)
		r.Put("/{id}", handleUpdateNote)
		r.Delete("/{id}", handleDeleteNote)
	}

	// v1 routes
	v1Routes := func(r chi.Router) {
		r.Route("/users", userRoutes)
		r.Route("/notes", noteRoutes)
	}

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", v1Routes)
	})

	// 404 handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "endpoint not found",
		})
	})

	fmt.Println("Server Running on :3000")
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), r)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "welcome to keeper APIs",
		"version": "0.0.1",
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// -----------------  User Handlers  -----------------

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid user ID",
		})
		return
	}

	var payload models.User
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request payload",
		})
		return
	}

	user, exists := userStore[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "user not found",
		})
		return
	}

	user.Name = payload.Name
	user.Email = payload.Email

	userStore[id] = user

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User updated.",
	})
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid user ID",
		})
		return
	}

	_, exists := userStore[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "user not found",
		})
		return
	}

	delete(userStore, id)
	w.WriteHeader(http.StatusNoContent)
}

// -----------------  Note Handlers  -----------------
func handleCreateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var note models.Note

	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request payload",
		})
		return
	}

	newID := len(noteStore) + 1
	note.ID = newID
	noteStore[newID] = note

	json.NewEncoder(w).Encode(map[string]any{
		"data": note,
	})
}

func handleListNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	noteList := make([]models.Note, 0, len(noteStore))

	for _, note := range noteStore {
		noteList = append(noteList, note)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"data": noteList,
	})
}

func handleGetNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid note ID",
		})
		return
	}

	note, exists := noteStore[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "note not found",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"data": note,
	})
}

func handleUpdateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid note ID",
		})
		return
	}

	var payload models.Note
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request payload",
		})
		return
	}

	note, exists := noteStore[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "user not found",
		})
		return
	}

	note.Title = payload.Title
	note.Content = payload.Content

	noteStore[id] = note

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Note updated.",
	})
}

func handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid note ID",
		})
		return
	}

	_, exists := noteStore[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "note not found",
		})
		return
	}

	delete(noteStore, id)
	w.WriteHeader(http.StatusNoContent)
}
