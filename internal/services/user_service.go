package services

import (
	"context"

	"github.com/go-playground/validator/v10"
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

	// Input validation
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, err
	}

	// Delegate to repository
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, input *models.UpdateUserInput, id int) error {
	// Input validation
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return err
	}

	return s.repo.Update(ctx, input, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
