package services

import (
	"context"

	"github.com/kevalsabhani/keeper/internal/models"
	"github.com/kevalsabhani/keeper/internal/repositories"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, input *models.CreateUserInput) (*models.User, error) {
	user := &models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	// TODO: Validation

	// Delegate to repository
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, input *models.CreateUserInput, id int) error {
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return nil
}
