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
