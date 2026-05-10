package services

import (
	"context"
	"time"

	errpkg "github.com/kevalsabhani/keeper/internal/errors"
	"github.com/kevalsabhani/keeper/internal/models"
	"github.com/kevalsabhani/keeper/internal/repositories"
	"github.com/kevalsabhani/keeper/internal/response"
)

// UserService contains the business logic for user operations.
type UserService struct {
	repo repositories.UserRepository
}

// NewUserService creates a UserService with the given repository dependency.
func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// CreateUser validates the input, then persists a new user to the database.
func (s *UserService) CreateUser(ctx context.Context, input *models.CreateUserInput) (*models.User, error) {
	user := &models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	// Input validation
	if err := user.Validate(); err != nil {
		return nil, errpkg.FromValidationError(err)
	}

	// Delegate to repository
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errpkg.FromDBError(err)
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID. Returns ErrNotFound if they do not exist.
func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errpkg.FromDBError(err)
	}
	return user, nil
}

// ListUsers returns a paginated list of users along with pagination metadata.
func (s *UserService) ListUsers(ctx context.Context, page, limit int) ([]*models.User, *response.Meta, error) {
	users, total, err := s.repo.List(ctx, page, limit)
	if err != nil {
		return nil, nil, errpkg.FromDBError(err)
	}

	return users, &response.Meta{
		CurrentPage: page,
		TotalPages:  (total + limit - 1) / limit,
		TotalCount:  total,
	}, nil
}

// UpdateUser applies partial changes to an existing user after re-validating the full record.
func (s *UserService) UpdateUser(ctx context.Context, input *models.UpdateUserInput, id int) error {

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errpkg.FromDBError(err)
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Email != nil {
		user.Email = *input.Email
	}

	// Input validation
	if err := user.Validate(); err != nil {
		return errpkg.FromValidationError(err)
	}

	user.UpdatedAt = time.Now()

	if err = s.repo.Update(ctx, user, id); err != nil {
		return errpkg.FromDBError(err)
	}

	return nil
}

// DeleteUser verifies the user exists, then removes them from the database.
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errpkg.FromDBError(err)
	}

	if err = s.repo.Delete(ctx, id); err != nil {
		return errpkg.FromDBError(err)
	}

	return nil
}
