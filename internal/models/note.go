package models

type Note struct {
	ID      int
	UserID  int
	Title   string
	Content string
}

type CreateNoteInput struct {
	UserID  int
	Title   string
	Content string
}
