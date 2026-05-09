package errors

import "errors"

// Domain-level sentinel errors
var (
	ErrNotFound = errors.New("resource not found")
)
