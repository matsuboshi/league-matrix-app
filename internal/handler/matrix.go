package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/matsuboshi/league-matrix-app/internal/domain"
	apperrors "github.com/matsuboshi/league-matrix-app/pkg/errors"
)

// MatrixHandlerInterface defines the contract for HTTP handlers that process matrix operations.
// It provides endpoints for listing available operations and processing matrices.
type MatrixHandlerInterface interface {
	// ListMatrixOperations handles requests to list all available matrix operations.
	// It responds with a text message showing available operations and a sample URL.
	ListMatrixOperations(w http.ResponseWriter, r *http.Request)

	// ProcessMatrix handles requests to perform specific matrix operations.
	// It extracts the operation from the URL path and the file path from query parameters,
	// then processes the matrix and returns the result.
	ProcessMatrix(w http.ResponseWriter, r *http.Request)

	// HealthCheck handles health check requests.
	// It returns HTTP 200 OK with "OK" message if the service is running and healthy.
	// This endpoint is intended for use with load balancers and container orchestration systems.
	HealthCheck(w http.ResponseWriter, r *http.Request)
}

type matrixHandler struct {
	matrixDomain domain.MatrixDomainInterface
}

// NewMatrixHandler creates a new instance of MatrixHandlerInterface with its dependencies.
// It initializes the handler with a matrix domain service for business logic processing.
func NewMatrixHandler() MatrixHandlerInterface {
	return &matrixHandler{
		matrixDomain: domain.NewMatrixDomain(),
	}
}

func (h *matrixHandler) ListMatrixOperations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := h.matrixDomain.ListMatrixOperations()
	if err != nil {
		statusCode := apperrors.GetHTTPStatusCode(err)
		slog.Error("failed to list operations",
			"error", err,
			"status_code", statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(result))
	if err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

func (h *matrixHandler) ProcessMatrix(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	operation := r.URL.Path[len("/matrix/"):]
	filePath := r.URL.Query().Get("file")

	result, err := h.matrixDomain.ProcessMatrix(r.Context(), operation, filePath)
	if err != nil {
		// Handle context errors specially
		if errors.Is(err, context.Canceled) {
			slog.Info("request cancelled by client",
				"operation", operation,
				"file_path", filePath)
			// Client already disconnected, no need to write response
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("request timeout",
				"operation", operation,
				"file_path", filePath)
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}

		// Handle other errors
		statusCode := apperrors.GetHTTPStatusCode(err)
		slog.Error("matrix operation failed",
			"operation", operation,
			"file_path", filePath,
			"error", err,
			"status_code", statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	slog.Info("matrix operation completed",
		"operation", operation,
		"file_path", filePath)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(result))
	if err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

func (h *matrixHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slog.Debug("health check request received")

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		slog.Error("failed to write health check response", "error", err)
	}
}
