package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// User represents a user record stored in the database.
type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=3,max=100"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserInput holds the fields required to create a new user.
type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserInput holds the optional fields for a partial user update.
// Pointer fields allow the caller to distinguish between "not provided" and zero value.
type UpdateUserInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

// Validate runs struct-level validation rules defined via `validate` struct tags.
func (user *User) Validate() error {
	return validator.New().Struct(user)
}
