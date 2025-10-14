package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	apperrors "github.com/matsuboshi/league-matrix-app/pkg/errors"
)

// mockMatrixDomain is a mock implementation of domain.MatrixDomainInterface for testing
type mockMatrixDomain struct {
	mock.Mock
}

func (m *mockMatrixDomain) ListMatrixOperations() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockMatrixDomain) ProcessMatrix(ctx context.Context, operation string, filePath string) (string, error) {
	args := m.Called(ctx, operation, filePath)
	return args.String(0), args.Error(1)
}

func TestMatrixHandler_ListMatrixOperations(t *testing.T) {
	tests := []struct {
		name             string
		method           string
		mockResponse     string
		mockError        error
		wantStatus       int
		wantBodyContains []string
		wantContentType  string
	}{
		{
			name:             "successfully list operations",
			method:           http.MethodGet,
			mockResponse:     "Are you lost?\nsum,multiply,echo,invert,flatten",
			mockError:        nil,
			wantStatus:       http.StatusOK,
			wantBodyContains: []string{"Are you lost?", "sum", "multiply"},
			wantContentType:  "text/plain",
		},
		{
			name:            "method not allowed - POST",
			method:          http.MethodPost,
			wantStatus:      http.StatusMethodNotAllowed,
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name:            "method not allowed - PUT",
			method:          http.MethodPut,
			wantStatus:      http.StatusMethodNotAllowed,
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name:            "method not allowed - DELETE",
			method:          http.MethodDelete,
			wantStatus:      http.StatusMethodNotAllowed,
			wantContentType: "text/plain; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock domain
			mockDomain := &mockMatrixDomain{}

			// Setup expectations only for GET requests
			if tt.method == http.MethodGet {
				mockDomain.On("ListMatrixOperations").Return(tt.mockResponse, tt.mockError)
			}

			// Create handler with mock
			handler := &matrixHandler{
				matrixDomain: mockDomain,
			}

			// Create request
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			// Execute
			handler.ListMatrixOperations(w, req)

			// Assert
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), tt.wantContentType)

			if tt.wantBodyContains != nil {
				for _, want := range tt.wantBodyContains {
					assert.Contains(t, w.Body.String(), want)
				}
			}
		})
	}
}

func TestMatrixHandler_ProcessMatrix(t *testing.T) {
	tests := []struct {
		name             string
		method           string
		path             string
		query            string
		mockResponse     string
		mockError        error
		wantStatus       int
		wantBodyContains string
		wantContentType  string
	}{
		{
			name:             "successfully process sum operation",
			method:           http.MethodGet,
			path:             "/matrix/sum",
			query:            "file=testdata/matrix1.csv",
			mockResponse:     "45",
			mockError:        nil,
			wantStatus:       http.StatusOK,
			wantBodyContains: "45",
			wantContentType:  "text/plain",
		},
		{
			name:             "successfully process multiply operation",
			method:           http.MethodGet,
			path:             "/matrix/multiply",
			query:            "file=testdata/matrix1.csv",
			mockResponse:     "362880",
			mockError:        nil,
			wantStatus:       http.StatusOK,
			wantBodyContains: "362880",
			wantContentType:  "text/plain",
		},
		{
			name:             "successfully process echo operation",
			method:           http.MethodGet,
			path:             "/matrix/echo",
			query:            "file=testdata/matrix1.csv",
			mockResponse:     "1,2,3\n4,5,6",
			mockError:        nil,
			wantStatus:       http.StatusOK,
			wantBodyContains: "1,2,3",
			wantContentType:  "text/plain",
		},
		{
			name:             "invalid operation",
			method:           http.MethodGet,
			path:             "/matrix/divide",
			query:            "file=testdata/matrix1.csv",
			mockError:        apperrors.ErrInvalidInput,
			wantStatus:       http.StatusBadRequest,
			wantBodyContains: "invalid input",
			wantContentType:  "text/plain; charset=utf-8",
		},
		{
			name:             "file not found",
			method:           http.MethodGet,
			path:             "/matrix/sum",
			query:            "file=testdata/notfound.csv",
			mockError:        apperrors.ErrNotFound,
			wantStatus:       http.StatusNotFound,
			wantBodyContains: "not found",
			wantContentType:  "text/plain; charset=utf-8",
		},
		{
			name:             "file too large",
			method:           http.MethodGet,
			path:             "/matrix/sum",
			query:            "file=testdata/large.csv",
			mockError:        apperrors.ErrPayloadTooLarge,
			wantStatus:       http.StatusRequestEntityTooLarge,
			wantBodyContains: "payload too large",
			wantContentType:  "text/plain; charset=utf-8",
		},
		{
			name:             "unprocessable entity",
			method:           http.MethodGet,
			path:             "/matrix/sum",
			query:            "file=testdata/matrix2.csv",
			mockError:        apperrors.ErrUnprocessableEntity,
			wantStatus:       http.StatusUnprocessableEntity,
			wantBodyContains: "unprocessable entity",
			wantContentType:  "text/plain; charset=utf-8",
		},
		{
			name:             "method not allowed - POST",
			method:           http.MethodPost,
			path:             "/matrix/sum",
			query:            "file=testdata/matrix1.csv",
			wantStatus:       http.StatusMethodNotAllowed,
			wantBodyContains: "method not allowed",
			wantContentType:  "text/plain; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock domain
			mockDomain := &mockMatrixDomain{}

			// Setup expectations only for GET requests
			if tt.method == http.MethodGet {
				operation := tt.path[len("/matrix/"):]
				filePath := ""
				if tt.query != "" {
					filePath = tt.query[len("file="):]
				}
				mockDomain.On("ProcessMatrix", mock.Anything, operation, filePath).
					Return(tt.mockResponse, tt.mockError)
			}

			// Create handler with mock
			handler := &matrixHandler{
				matrixDomain: mockDomain,
			}

			// Create request
			url := tt.path
			if tt.query != "" {
				url += "?" + tt.query
			}
			req := httptest.NewRequest(tt.method, url, nil)
			w := httptest.NewRecorder()

			// Execute
			handler.ProcessMatrix(w, req)

			// Assert
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), tt.wantContentType)
			assert.Contains(t, w.Body.String(), tt.wantBodyContains)
		})
	}
}

func TestMatrixHandler_ProcessMatrix_ContextHandling(t *testing.T) {
	t.Run("context cancelled by client", func(t *testing.T) {
		mockDomain := &mockMatrixDomain{}
		mockDomain.On("ProcessMatrix", mock.Anything, "sum", "testdata/matrix1.csv").
			Return("", context.Canceled)

		handler := &matrixHandler{
			matrixDomain: mockDomain,
		}

		req := httptest.NewRequest(http.MethodGet, "/matrix/sum?file=testdata/matrix1.csv", nil)
		w := httptest.NewRecorder()

		handler.ProcessMatrix(w, req)

		// When context is cancelled, we don't write a response
		// The response should be empty or minimal
		assert.Equal(t, http.StatusOK, w.Code) // Default status when no write happens
	})

	t.Run("context deadline exceeded", func(t *testing.T) {
		mockDomain := &mockMatrixDomain{}
		mockDomain.On("ProcessMatrix", mock.Anything, "sum", "testdata/matrix1.csv").
			Return("", context.DeadlineExceeded)

		handler := &matrixHandler{
			matrixDomain: mockDomain,
		}

		req := httptest.NewRequest(http.MethodGet, "/matrix/sum?file=testdata/matrix1.csv", nil)
		w := httptest.NewRecorder()

		handler.ProcessMatrix(w, req)

		assert.Equal(t, http.StatusGatewayTimeout, w.Code)
		assert.Contains(t, w.Body.String(), "request timeout")
	})
}

func TestMatrixHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		wantStatus      int
		wantBody        string
		wantContentType string
	}{
		{
			name:            "successful health check",
			method:          http.MethodGet,
			wantStatus:      http.StatusOK,
			wantBody:        "OK",
			wantContentType: "text/plain",
		},
		{
			name:            "method not allowed - POST",
			method:          http.MethodPost,
			wantStatus:      http.StatusMethodNotAllowed,
			wantBody:        "method not allowed",
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name:            "method not allowed - PUT",
			method:          http.MethodPut,
			wantStatus:      http.StatusMethodNotAllowed,
			wantBody:        "method not allowed",
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name:            "method not allowed - DELETE",
			method:          http.MethodDelete,
			wantStatus:      http.StatusMethodNotAllowed,
			wantBody:        "method not allowed",
			wantContentType: "text/plain; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &matrixHandler{}

			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			handler.HealthCheck(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), tt.wantContentType)
			assert.Contains(t, w.Body.String(), tt.wantBody)
		})
	}
}

func TestMatrixHandler_ErrorHandling(t *testing.T) {
	t.Run("domain error is properly mapped to HTTP status", func(t *testing.T) {
		mockDomain := &mockMatrixDomain{}
		mockDomain.On("ProcessMatrix", mock.Anything, "sum", "invalid").
			Return("", errors.New("some domain error"))

		handler := &matrixHandler{
			matrixDomain: mockDomain,
		}

		req := httptest.NewRequest(http.MethodGet, "/matrix/sum?file=invalid", nil)
		w := httptest.NewRecorder()

		handler.ProcessMatrix(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("list operations error handling", func(t *testing.T) {
		mockDomain := &mockMatrixDomain{}
		mockDomain.On("ListMatrixOperations").
			Return("", errors.New("internal error"))

		handler := &matrixHandler{
			matrixDomain: mockDomain,
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler.ListMatrixOperations(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestNewMatrixHandler(t *testing.T) {
	t.Run("creates handler with dependencies", func(t *testing.T) {
		handler := NewMatrixHandler()

		assert.NotNil(t, handler)
		// Verify it implements the interface
		var _ MatrixHandlerInterface = handler
	})
}
