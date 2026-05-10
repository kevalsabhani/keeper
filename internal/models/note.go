package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Note represents a note record stored in the database.
type Note struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id" validate:"required"`
	Title     string    `json:"title" db:"title" validate:"required,min=3,max=100"`
	Content   string    `json:"content" db:"content" validate:"required,min=3,max=100"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateNoteInput holds the fields required to create a new note.
type CreateNoteInput struct {
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdateNoteInput holds the optional fields for a partial note update.
// Pointer fields allow the caller to distinguish between "not provided" and zero value.
type UpdateNoteInput struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

// Validate runs struct-level validation rules defined via `validate` struct tags.
func (note *Note) Validate() error {
	return validator.New().Struct(note)
}
