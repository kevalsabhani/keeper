package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kevalsabhani/keeper/internal/configs"
	"github.com/kevalsabhani/keeper/internal/database"
	"github.com/kevalsabhani/keeper/internal/handlers"
	"github.com/kevalsabhani/keeper/internal/repositories"
	"github.com/kevalsabhani/keeper/internal/services"
)

func main() {

	config, err := configs.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

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

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)

		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err = server.Shutdown(ctx); err != nil {
			// TODO: log error and continue
		}
		db.Close()
	}()

	fmt.Println("Server Running on :3000")
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
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
