package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevalsabhani/keeper/internal/models"
)

// UserRepository defines the data access contract for users.
type UserRepository interface {
	Create(context.Context, *models.User) error
	GetByID(context.Context, int) (*models.User, error)
	GetByEmail(context.Context, string) (*models.User, error)
	List(context.Context, int, int) ([]*models.User, int, error)
	Update(context.Context, *models.User, int) error
	Delete(context.Context, int) error
}

// PostgresUserRepository is the PostgreSQL implementation of UserRepository.
type PostgresUserRepository struct {
	db *pgxpool.Pool
}

// NewPostgresUserRepository creates a PostgresUserRepository backed by the given pool.
func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// Create inserts a new user and scans the generated id, created_at, and updated_at back into user.
func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users(name, email) VALUES ($1, $2) RETURNING id, created_at, updated_at"
	return r.db.QueryRow(ctx, query, user.Name, user.Email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetByID fetches a single user by their primary key. Returns an error if not found.
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, email, created_at, updated_at from users WHERE id=$1"
	row := r.db.QueryRow(ctx, query, id)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail fetches a single user by their email address. Returns an error if not found.
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, email, created_at, updated_at from users WHERE email=$1"
	row := r.db.QueryRow(ctx, query, email)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

// List returns a paginated slice of users ordered by id, along with the total row count.
func (r *PostgresUserRepository) List(ctx context.Context, page, limit int) ([]*models.User, int, error) {
	var total int
	users := make([]*models.User, 0)

	query := "SELECT COUNT(id) FROM users"
	row := r.db.QueryRow(ctx, query)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	query = "SELECT id, name, email, created_at, updated_at from users ORDER BY id ASC LIMIT $1 OFFSET $2"
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// Update persists changes to an existing user's name, email, and updated_at.
func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User, id int) error {

	query := "UPDATE users SET name=$1, email=$2, updated_at=$3 WHERE id=$4"
	if _, err := r.db.Exec(ctx, query, user.Name, user.Email, user.UpdatedAt, id); err != nil {
		return err
	}

	return nil
}

// Delete removes a user by their primary key.
func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id=$1"
	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}
