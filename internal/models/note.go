package models

import "time"

type Note struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateNoteInput struct {
	UserID  int    `json:"user_id" validate:"required"`
	Title   string `json:"title" validate:"required,min=3,max=100"`
	Content string `json:"content" validate:"required,min=3,max=100"`
}

type UpdateNoteInput struct {
	Title   *string `json:"title" validate:"omitempty,min=3,max=100"`
	Content *string `json:"content" validate:"omitempty,min=3,max=100"`
}
