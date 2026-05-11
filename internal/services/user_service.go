package services

import (
	"context"
	"time"

	errpkg "github.com/kevalsabhani/keeper/internal/errors"
	"github.com/kevalsabhani/keeper/internal/models"
	"github.com/kevalsabhani/keeper/internal/repositories"
	"github.com/kevalsabhani/keeper/internal/response"
	"go.uber.org/zap"
)

// UserService contains the business logic for user operations.
type UserService struct {
	repo repositories.UserRepository
	log  *zap.Logger
}

// NewUserService creates a UserService with the given repository dependency.
func NewUserService(repo repositories.UserRepository, log *zap.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
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
		s.log.Warn("create user validation failed", zap.Error(err))
		return nil, errpkg.FromValidationError(err)
	}

	// Delegate to repository
	if err := s.repo.Create(ctx, user); err != nil {
		s.log.Error("failed to insert user into db", zap.Error(err))
		return nil, errpkg.FromDBError(err)
	}

	s.log.Info("user created", zap.Int("id", user.ID), zap.String("email", user.Email))
	return user, nil
}

// GetUserByID retrieves a user by their ID. Returns ErrNotFound if they do not exist.
func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to fetch user from db", zap.Int("id", id), zap.Error(err))
		return nil, errpkg.FromDBError(err)
	}
	return user, nil
}

// ListUsers returns a paginated list of users along with pagination metadata.
func (s *UserService) ListUsers(ctx context.Context, page, limit int) ([]*models.User, *response.Meta, error) {
	users, total, err := s.repo.List(ctx, page, limit)
	if err != nil {
		s.log.Error("failed to fetch users from db", zap.Int("page", page), zap.Int("limit", limit), zap.Error(err))
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
		s.log.Error("failed to fetch user for update", zap.Int("id", id), zap.Error(err))
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
		s.log.Warn("update user validation failed", zap.Int("id", id), zap.Error(err))
		return errpkg.FromValidationError(err)
	}

	user.UpdatedAt = time.Now()

	if err = s.repo.Update(ctx, user, id); err != nil {
		s.log.Error("failed to update user in db", zap.Int("id", id), zap.Error(err))
		return errpkg.FromDBError(err)
	}

	s.log.Info("user updated", zap.Int("id", id))
	return nil
}

// DeleteUser verifies the user exists, then removes them from the database.
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to fetch user for deletion", zap.Int("id", id), zap.Error(err))
		return errpkg.FromDBError(err)
	}

	if err = s.repo.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete user from db", zap.Int("id", id), zap.Error(err))
		return errpkg.FromDBError(err)
	}

	s.log.Info("user deleted", zap.Int("id", id))
	return nil
}
