package errors

import (
	"errors"
	"net/http"
)

// Sentinel errors for error classification across the application.
// These errors are used to wrap domain-specific errors and map them to appropriate HTTP status codes.
var (
	// ErrInvalidInput maps to 400 Bad Request.
	ErrInvalidInput = errors.New("invalid input")

	// ErrNotFound maps to 404 Not Found.
	ErrNotFound = errors.New("not found")

	// ErrPayloadTooLarge maps to 413 Payload Too Large.
	ErrPayloadTooLarge = errors.New("payload too large")

	// ErrUnprocessableEntity maps to 422 Unprocessable Entity.
	ErrUnprocessableEntity = errors.New("unprocessable entity")
)

// GetHTTPStatusCode maps application errors to appropriate HTTP status codes.
// It uses errors.Is to check the error chain for sentinel errors and returns the corresponding status code.
// If no sentinel error is found, it defaults to 500 Internal Server Error.
func GetHTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch {
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest // 400
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound // 404
	case errors.Is(err, ErrPayloadTooLarge):
		return http.StatusRequestEntityTooLarge // 413
	case errors.Is(err, ErrUnprocessableEntity):
		return http.StatusUnprocessableEntity // 422
	default:
		return http.StatusInternalServerError // 500
	}
}
