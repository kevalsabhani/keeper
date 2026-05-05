package services

import (
	"context"

	"github.com/kevalsabhani/keeper/internal/models"
	"github.com/kevalsabhani/keeper/internal/repositories"
)

type NoteService struct {
	repo repositories.NoteRepository
}

func NewNoteService(repo repositories.NoteRepository) *NoteService {
	return &NoteService{
		repo: repo,
	}
}

func (s *NoteService) CreateNote(ctx context.Context, input *models.CreateNoteInput) (*models.Note, error) {
	note := &models.Note{
		UserID:  input.UserID,
		Title:   input.Title,
		Content: input.Content,
	}

	// Validation

	// Delegate to repository
	if err := s.repo.Create(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}
