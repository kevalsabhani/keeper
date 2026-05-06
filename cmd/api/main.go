package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		r.Patch("/{id}", userHandler.Update)
		r.Delete("/{id}", userHandler.Delete)
	}

	// Note routes
	noteRoutes := func(r chi.Router) {
		r.Post("/", noteHandler.Create)
		r.Get("/", noteHandler.List)
		r.Get("/{id}", noteHandler.GetByID)
		r.Patch("/{id}", noteHandler.Update)
		r.Delete("/{id}", noteHandler.Delete)
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
