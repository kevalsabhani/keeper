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
	Update(context.Context, *models.User, int) error
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
	query := "SELECT id, name, email, created_at, updated_at from users WHERE id=$1"
	row := r.db.QueryRow(ctx, query, id)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) List(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	query := "SELECT id, name, email, created_at, updated_at from users"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User, id int) error {

	query := "UPDATE users SET name=$1, email=$2, updated_at=$3 WHERE id=$4"
	if _, err := r.db.Exec(ctx, query, user.Name, user.Email, user.UpdatedAt, id); err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id=$1"
	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}
