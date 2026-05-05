package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevalsabhani/keeper/internal/models"
)

type UserRepository interface {
	Create(context.Context, *models.User) error
	GetByID(context.Context, int) (*models.User, error)
	List(context.Context) ([]*models.User, error)
	Update(context.Context, *models.User) error
	Delete(context.Context, int) error
}

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users(name, email) VALUES ($1, $2) RETURNING id"
	return r.db.QueryRow(ctx, query, user.Name, user.Email).Scan(&user.ID)
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) List(ctx context.Context) ([]*models.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	return nil
}
