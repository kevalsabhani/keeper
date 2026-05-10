package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevalsabhani/keeper/internal/handlers"
	"github.com/kevalsabhani/keeper/internal/repositories"
	"github.com/kevalsabhani/keeper/internal/services"
)

// Container holds all top-level HTTP handlers wired with their dependencies.
type Container struct {
	UserHandler *handlers.UserHandler
	NoteHandler *handlers.NoteHandler
}

// New builds the full dependency graph — repositories → services → handlers —
// and returns a Container ready to be used by the router.
func New(db *pgxpool.Pool) *Container {

	// Setup repositories
	userRepository := repositories.NewPostgresUserRepository(db)
	noteRepository := repositories.NewPostgresNoteRepository(db)

	// Setup services
	userService := services.NewUserService(userRepository)
	noteService := services.NewNoteService(noteRepository)

	// Setup handlers
	return &Container{
		UserHandler: handlers.NewUserHandler(userService),
		NoteHandler: handlers.NewNoteHandler(noteService),
	}
}
