package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevalsabhani/keeper/internal/models"
)

type NoteRepository interface {
	Create(context.Context, *models.Note) error
	GetByID(context.Context, int) (*models.Note, error)
	List(context.Context) ([]*models.Note, error)
	Update(context.Context, *models.Note, int) error
	Delete(context.Context, int) error
}

type PostgresNoteRepository struct {
	db *pgxpool.Pool
}

func NewPostgresNoteRepository(db *pgxpool.Pool) *PostgresNoteRepository {
	return &PostgresNoteRepository{
		db: db,
	}
}

func (r *PostgresNoteRepository) Create(ctx context.Context, note *models.Note) error {

	query := "INSERT INTO notes(title, user_id, content) VALUES ($1, $2, $3) RETURNING id"
	return r.db.QueryRow(ctx, query, note.Title, note.UserID, note.Content).Scan(&note.ID)
}

func (r *PostgresNoteRepository) GetByID(ctx context.Context, id int) (*models.Note, error) {
	var note models.Note

	query := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id=$1"

	row := r.db.QueryRow(ctx, query, id)

	if err := row.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *PostgresNoteRepository) List(ctx context.Context) ([]*models.Note, error) {
	var notes []*models.Note

	query := "SELECT id, user_id, title, content, created_at, updated_at FROM notes"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var note models.Note

		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}

		notes = append(notes, &note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *PostgresNoteRepository) Update(ctx context.Context, note *models.Note, id int) error {

	query := "UPDATE notes SET title=$1, content=$2, updated_at=$3 WHERE id=$4"

	if _, err := r.db.Exec(ctx, query, note.Title, note.Content, note.UpdatedAt, id); err != nil {
		return err
	}
	return nil
}

func (r *PostgresNoteRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM notes WHERE id=$1"
	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}
