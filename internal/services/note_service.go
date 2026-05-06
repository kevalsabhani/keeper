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

func (s *NoteService) GetNoteByID(ctx context.Context, id int) (*models.Note, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *NoteService) ListNotes(ctx context.Context) ([]*models.Note, error) {
	return s.repo.List(ctx)
}

func (s *NoteService) UpdateNote(ctx context.Context, input *models.UpdateNoteInput, id int) error {

	return s.repo.Update(ctx, input, id)
}

func (s *NoteService) DeleteNote(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
