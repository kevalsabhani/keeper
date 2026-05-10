package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	errpkg "github.com/kevalsabhani/keeper/internal/errors"
)

// Response is the standard JSON envelope returned by all API endpoints.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError carries a machine-readable error code, a human-readable message,
// and optional per-field validation details.
type APIError struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Fields  []errpkg.FieldError `json:"fields,omitempty"`
}

// Meta holds pagination information for list responses.
type Meta struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
}

// JSON sets Content-Type to application/json, writes the given status code,
// and encodes resp as JSON into the response body.
func JSON(w http.ResponseWriter, status int, resp *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// TODO: log error and continue
	}
}

// Success writes a successful JSON response. If data is a struct (or pointer to struct),
// it is automatically wrapped in a slice so the data field is always an array.
func Success(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	if data != nil {
		v := reflect.Indirect(reflect.ValueOf(data))
		if v.Kind() == reflect.Struct {
			data = []interface{}{data}
		}
	}

	JSON(w, status, &Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Error maps the given domain error to an HTTP status code and writes a JSON
// error response. Validation errors include a per-field breakdown in Fields.
func Error(w http.ResponseWriter, err error) {
	status, code := errpkg.MapErrorToStatusCode(err)

	apiErr := &APIError{
		Code:    code,
		Message: err.Error(),
	}

	var ve *errpkg.ValidationError
	if errors.As(err, &ve) {
		apiErr.Fields = ve.Fields
	}

	JSON(w, status, &Response{
		Success: false,
		Error:   apiErr,
	})
}
