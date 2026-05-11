package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevalsabhani/keeper/internal/models"
)

// NoteRepository defines the data access contract for notes.
type NoteRepository interface {
	Create(context.Context, *models.Note) error
	GetByID(context.Context, int) (*models.Note, error)
	List(context.Context, int, int) ([]*models.Note, int, error)
	Update(context.Context, *models.Note, int) error
	Delete(context.Context, int) error
}

// PostgresNoteRepository is the PostgreSQL implementation of NoteRepository.
type PostgresNoteRepository struct {
	db *pgxpool.Pool
}

// NewPostgresNoteRepository creates a PostgresNoteRepository backed by the given pool.
func NewPostgresNoteRepository(db *pgxpool.Pool) NoteRepository {
	return &PostgresNoteRepository{
		db: db,
	}
}

// Create inserts a new note and scans the generated id, created_at, and updated_at back into note.
func (r *PostgresNoteRepository) Create(ctx context.Context, note *models.Note) error {

	query := "INSERT INTO notes(title, user_id, content) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at"
	return r.db.QueryRow(ctx, query, note.Title, note.UserID, note.Content).Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
}

// GetByID fetches a single note by its primary key. Returns an error if not found.
func (r *PostgresNoteRepository) GetByID(ctx context.Context, id int) (*models.Note, error) {
	var note models.Note

	query := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id=$1"

	row := r.db.QueryRow(ctx, query, id)

	if err := row.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
		return nil, err
	}
	return &note, nil
}

// List returns a paginated slice of notes ordered by id, along with the total row count.
func (r *PostgresNoteRepository) List(ctx context.Context, page, limit int) ([]*models.Note, int, error) {
	var total int
	notes := make([]*models.Note, 0)

	query := "SELECT COUNT(id) FROM notes"
	row := r.db.QueryRow(ctx, query)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	query = "SELECT id, user_id, title, content, created_at, updated_at FROM notes ORDER BY id ASC LIMIT $1 OFFSET $2"

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	for rows.Next() {
		var note models.Note

		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, 0, err
		}

		notes = append(notes, &note)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return notes, total, nil
}

// Update persists changes to an existing note's title, content, and updated_at.
func (r *PostgresNoteRepository) Update(ctx context.Context, note *models.Note, id int) error {

	query := "UPDATE notes SET title=$1, content=$2, updated_at=$3 WHERE id=$4"

	if _, err := r.db.Exec(ctx, query, note.Title, note.Content, note.UpdatedAt, id); err != nil {
		return err
	}
	return nil
}

// Delete removes a note by its primary key.
func (r *PostgresNoteRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM notes WHERE id=$1"
	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}
