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
	var user models.User
	query := "SELECT id, name, email from users WHERE id=$1"
	row := r.db.QueryRow(ctx, query, id)

	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) List(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	query := "SELECT id, name, email from users"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	return nil
}
