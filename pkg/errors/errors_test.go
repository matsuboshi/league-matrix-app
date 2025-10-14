package errors

import (
	"errors"
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
			name:     "ErrInvalidInput returns 400",
			err:      ErrInvalidInput,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "wrapped ErrInvalidInput returns 400",
			err:      errors.New("wrapped: " + ErrInvalidInput.Error()),
			wantCode: http.StatusInternalServerError, // Not wrapped with %w, so won't match
		},
		{
			name:     "properly wrapped ErrInvalidInput returns 400",
			err:      errors.Join(ErrInvalidInput, errors.New("additional context")),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "ErrNotFound returns 404",
			err:      ErrNotFound,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "ErrPayloadTooLarge returns 413",
			err:      ErrPayloadTooLarge,
			wantCode: http.StatusRequestEntityTooLarge,
		},
		{
			name:     "ErrUnprocessableEntity returns 422",
			err:      ErrUnprocessableEntity,
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

func TestSentinelErrors(t *testing.T) {
	t.Run("sentinel errors are unique", func(t *testing.T) {
		assert.NotEqual(t, ErrInvalidInput, ErrNotFound)
		assert.NotEqual(t, ErrInvalidInput, ErrPayloadTooLarge)
		assert.NotEqual(t, ErrInvalidInput, ErrUnprocessableEntity)
		assert.NotEqual(t, ErrNotFound, ErrPayloadTooLarge)
		assert.NotEqual(t, ErrNotFound, ErrUnprocessableEntity)
		assert.NotEqual(t, ErrPayloadTooLarge, ErrUnprocessableEntity)
	})

	t.Run("sentinel errors have correct messages", func(t *testing.T) {
		assert.Equal(t, "invalid input", ErrInvalidInput.Error())
		assert.Equal(t, "not found", ErrNotFound.Error())
		assert.Equal(t, "payload too large", ErrPayloadTooLarge.Error())
		assert.Equal(t, "unprocessable entity", ErrUnprocessableEntity.Error())
	})
}
