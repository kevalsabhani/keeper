package services

import (
	"context"
	"time"

	errpkg "github.com/kevalsabhani/keeper/internal/errors"
	"github.com/kevalsabhani/keeper/internal/models"
	"github.com/kevalsabhani/keeper/internal/repositories"
	"github.com/kevalsabhani/keeper/internal/response"
)

// NoteService contains the business logic for note operations.
type NoteService struct {
	repo repositories.NoteRepository
}

// NewNoteService creates a NoteService with the given repository dependency.
func NewNoteService(repo repositories.NoteRepository) *NoteService {
	return &NoteService{
		repo: repo,
	}
}

// CreateNote validates the input, then persists a new note to the database.
func (s *NoteService) CreateNote(ctx context.Context, input *models.CreateNoteInput) (*models.Note, error) {
	note := &models.Note{
		UserID:  input.UserID,
		Title:   input.Title,
		Content: input.Content,
	}

	// Input validation
	if err := note.Validate(); err != nil {
		return nil, errpkg.FromValidationError(err)
	}

	// Delegate to repository
	if err := s.repo.Create(ctx, note); err != nil {
		return nil, errpkg.FromDBError(err)
	}

	return note, nil
}

// GetNoteByID retrieves a note by its ID. Returns ErrNotFound if it does not exist.
func (s *NoteService) GetNoteByID(ctx context.Context, id int) (*models.Note, error) {
	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errpkg.FromDBError(err)
	}

	return note, nil
}

// ListNotes returns a paginated list of notes along with pagination metadata.
func (s *NoteService) ListNotes(ctx context.Context, page, limit int) ([]*models.Note, *response.Meta, error) {
	notes, total, err := s.repo.List(ctx, page, limit)
	if err != nil {
		return nil, nil, errpkg.FromDBError(err)
	}

	return notes, &response.Meta{
		CurrentPage: page,
		TotalPages:  (total + limit - 1) / limit,
		TotalCount:  total,
	}, nil
}

// UpdateNote applies partial changes to an existing note after re-validating the full record.
func (s *NoteService) UpdateNote(ctx context.Context, input *models.UpdateNoteInput, id int) error {

	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errpkg.FromDBError(err)
	}

	if input.Title != nil {
		note.Title = *input.Title
	}

	if input.Content != nil {
		note.Content = *input.Content
	}

	// Input validation
	if err := note.Validate(); err != nil {
		return errpkg.FromValidationError(err)
	}

	note.UpdatedAt = time.Now()

	if err = s.repo.Update(ctx, note, id); err != nil {
		return errpkg.FromDBError(err)
	}

	return nil
}

// DeleteNote verifies the note exists, then removes it from the database.
func (s *NoteService) DeleteNote(ctx context.Context, id int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errpkg.FromDBError(err)
	}

	if err = s.repo.Delete(ctx, id); err != nil {
		return errpkg.FromDBError(err)
	}

	return nil
}
