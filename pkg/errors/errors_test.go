package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHTTPStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "nil error returns 200",
			err:      nil,
			wantCode: http.StatusOK,
		},
		{
			name:     "400 ErrInvalidInput",
			err:      ErrInvalidInput,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "improper wrapping for 400 ErrInvalidInput",
			err:      errors.New("wrapped: " + ErrInvalidInput.Error()),
			wantCode: http.StatusInternalServerError, // Not wrapped with %w, so won't match
		},
		{
			name:     "fmt.Errorf with %w wrapping 400 ErrInvalidInput",
			err:      fmt.Errorf("%w: invalid operation: multiply", ErrInvalidInput),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "ErrNotFound returns 404",
			err:      ErrNotFound,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "fmt.Errorf with %w wrapping 404 ErrNotFound",
			err:      fmt.Errorf("%w: matrix not found with id: 123", ErrNotFound),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "ErrPayloadTooLarge returns 413",
			err:      ErrPayloadTooLarge,
			wantCode: http.StatusRequestEntityTooLarge,
		},
		{
			name:     "fmt.Errorf with %w wrapping 413 ErrPayloadTooLarge",
			err:      fmt.Errorf("%w: matrix size exceeds limit", ErrPayloadTooLarge),
			wantCode: http.StatusRequestEntityTooLarge,
		},
		{
			name:     "ErrUnprocessableEntity returns 422",
			err:      ErrUnprocessableEntity,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "fmt.Errorf with %w wrapping 422 ErrUnprocessableEntity",
			err:      fmt.Errorf("%w: unable to process matrix format", ErrUnprocessableEntity),
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "unknown error returns 500",
			err:      errors.New("unknown error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "custom error returns 500",
			err:      errors.New("custom application error"),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHTTPStatusCode(tt.err)
			assert.Equal(t, tt.wantCode, got)
		})
	}
}
