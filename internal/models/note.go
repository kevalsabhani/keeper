package models

type Note struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreateNoteInput struct {
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateNoteInput struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}
