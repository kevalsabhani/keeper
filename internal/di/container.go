package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevalsabhani/keeper/internal/handlers"
	"github.com/kevalsabhani/keeper/internal/repositories"
	"github.com/kevalsabhani/keeper/internal/services"
)

type Container struct {
	userHandler *handlers.UserHandler
	noteHandler *handlers.NoteHandler
}

func New(db *pgxpool.Pool) *Container {

	// Setup repositories
	userRepository := repositories.NewPostgresUserRepository(db)
	noteRepository := repositories.NewPostgresNoteRepository(db)

	// Setup services
	userService := services.NewUserService(userRepository)
	noteService := services.NewNoteService(noteRepository)

	// Setup handlers
	return &Container{
		userHandler: handlers.NewUserHandler(userService),
		noteHandler: handlers.NewNoteHandler(noteService),
	}
}
