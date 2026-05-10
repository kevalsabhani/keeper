package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ValidationError carries per-field validation failure messages.
type ValidationError struct {
	Fields []FieldError
}

// FieldError represents a single field's validation failure.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *ValidationError) Error() string {
	return ErrInvalidInput.Error()
}

// FromValidationError converts go-playground/validator errors into a ValidationError.
func FromValidationError(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make([]FieldError, 0, len(ve))
		for _, fe := range ve {
			fields = append(fields, FieldError{
				Field:   fe.Field(),
				Message: validationMessage(fe),
			})
		}
		return &ValidationError{Fields: fields}
	}
	return ErrInvalidInput
}

// validationMessage returns a human-readable message for a single field error.
func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}

// Domain-level sentinel errors
var (
	ErrNotFound            = errors.New("resource not found")
	ErrInvalidInput        = errors.New("invalid input")
	ErrConflict            = errors.New("resource already exists")
	ErrInternalServerError = errors.New("internal server error")
)

// Error Codes
var (
	ErrCodeInvalidInput        = "INVALID_INPUT"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrCodeConflict            = "CONFLICT"
)

// FromDBError maps DB errors to domain-level sentinel errors
func FromDBError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrConflict
		case "23503", "23502", "23514", "22001":
			return ErrInvalidInput
		}
	}

	return ErrInternalServerError
}

// MapErrorToStatusCode maps a domain-level sentinel error to an HTTP status code
// and an error code string suitable for use in API responses.
func MapErrorToStatusCode(err error) (int, string) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return http.StatusBadRequest, ErrCodeInvalidInput
	}
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, ErrCodeNotFound
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest, ErrCodeInvalidInput
	case errors.Is(err, ErrConflict):
		return http.StatusConflict, ErrCodeConflict
	default:
		return http.StatusInternalServerError, ErrCodeInternalServerError
	}
}
