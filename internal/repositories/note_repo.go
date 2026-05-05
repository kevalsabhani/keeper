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
	Update(context.Context, *models.Note) error
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
	return nil, nil
}

func (r *PostgresNoteRepository) List(ctx context.Context) ([]*models.Note, error) {
	return nil, nil
}

func (r *PostgresNoteRepository) Update(ctx context.Context, note *models.Note) error {
	return nil
}

func (r *PostgresNoteRepository) Delete(ctx context.Context, id int) error {
	return nil
}
