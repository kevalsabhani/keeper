package services

import (
	"context"
	"time"

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

	// Input alidation
	if err := note.Validate(); err != nil {
		return nil, err
	}

	// Delegate to repository
	if err := s.repo.Create(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteService) GetNoteByID(ctx context.Context, id int) (*models.Note, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *NoteService) ListNotes(ctx context.Context) ([]*models.Note, error) {
	return s.repo.List(ctx)
}

func (s *NoteService) UpdateNote(ctx context.Context, input *models.UpdateNoteInput, id int) error {

	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if input.Title != nil {
		note.Title = *input.Title
	}

	if input.Content != nil {
		note.Content = *input.Content
	}

	//Input validation
	if err := note.Validate(); err != nil {
		return err
	}

	note.UpdatedAt = time.Now()

	return s.repo.Update(ctx, note, id)
}

func (s *NoteService) DeleteNote(ctx context.Context, id int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
