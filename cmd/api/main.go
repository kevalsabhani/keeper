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
	"github.com/kevalsabhani/keeper/internal/di"
)

// main is the application entry point. It loads config, connects to the
// database, wires up the router, and starts the HTTP server with graceful shutdown.
func main() {

	config, err := configs.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := database.New(config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	container := di.New(db)

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	// Routes
	r.Get("/", handleHome)
	r.Get("/health", handleHealth)

	// User routes
	userRoutes := func(r chi.Router) {
		r.Post("/", container.UserHandler.Create)
		r.Get("/", container.UserHandler.List)
		r.Get("/{id}", container.UserHandler.GetByID)
		r.Patch("/{id}", container.UserHandler.Update)
		r.Delete("/{id}", container.UserHandler.Delete)
	}

	// Note routes
	noteRoutes := func(r chi.Router) {
		r.Post("/", container.NoteHandler.Create)
		r.Get("/", container.NoteHandler.List)
		r.Get("/{id}", container.NoteHandler.GetByID)
		r.Patch("/{id}", container.NoteHandler.Update)
		r.Delete("/{id}", container.NoteHandler.Delete)
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
		ReadTimeout:  time.Duration(config.ReadTimeOut) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeOut) * time.Second,
		IdleTimeout:  time.Duration(config.IdleTimeOut) * time.Second,
	}

	// Listen for OS signals and trigger graceful shutdown with a 10s drain window.
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

// handleHome returns a welcome message and the current API version.
func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "welcome to keeper APIs",
		"version": "0.0.1",
	})
}

// handleHealth returns a simple status check used by load balancers and monitoring tools.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
