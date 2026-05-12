package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	errpkg "github.com/kevalsabhani/keeper/internal/errors"
	"github.com/kevalsabhani/keeper/internal/models"
	"go.uber.org/zap"
)

// Response is the standard JSON envelope returned by all API endpoints.
type Response struct {
	Success bool               `json:"success"`
	Data    interface{}        `json:"data,omitempty"`
	Error   *APIError          `json:"error,omitempty"`
	Meta    *models.Pagination `json:"meta,omitempty"`
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

	respBytes, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error("failed to marshal response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(respBytes)
}

// Success writes a successful JSON response. If data is a struct (or pointer to struct),
// it is automatically wrapped in a slice so the data field is always an array.
func Success(w http.ResponseWriter, status int, data interface{}, meta *models.Pagination) {
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
