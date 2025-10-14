package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/matsuboshi/league-matrix-app/internal/entity"
	"github.com/matsuboshi/league-matrix-app/internal/repository"
	apperrors "github.com/matsuboshi/league-matrix-app/pkg/errors"
)

// mockMatrixRepository is a mock implementation of repository.MatrixRepositoryInterface for testing
type mockMatrixRepository struct {
	mock.Mock
}

func (m *mockMatrixRepository) GetFileContent(ctx context.Context, filePath string) (*repository.MatrixFileContent, error) {
	args := m.Called(ctx, filePath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.MatrixFileContent), args.Error(1)
}

func TestMatrixDomain_ListMatrixOperations(t *testing.T) {
	tests := []struct {
		name           string
		mockOperations []string
		wantContains   []string
		wantErr        bool
	}{
		{
			name:           "successfully list operations",
			mockOperations: []string{"sum", "multiply", "echo", "invert", "flatten"},
			wantContains:   []string{"Are you lost?", "http://localhost:8080/matrix/sum?file=testdata/matrix1.csv", "sum", "multiply", "echo", "invert", "flatten"},
			wantErr:        false,
		},
		{
			name:           "list with single operation",
			mockOperations: []string{"sum"},
			wantContains:   []string{"Are you lost?", "sum"},
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks using Mockery v3
			mockOperations := NewMockMatrixOperationsDomainInterface(t)

			// Setup expectations using testify/mock syntax
			mockOperations.On("ListOperations").Return(tt.mockOperations)

			// Create domain with mocked dependencies
			domain := &matrixDomain{
				operationsDomain: mockOperations,
			}

			// Execute
			got, err := domain.ListMatrixOperations()

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				for _, want := range tt.wantContains {
					assert.Contains(t, got, want)
				}
			}
		})
	}
}

func TestMatrixDomain_ProcessMatrix(t *testing.T) {
	tests := []struct {
		name              string
		operation         string
		filePath          string
		mockFileContent   *repository.MatrixFileContent
		mockMatrix        *entity.Matrix
		mockResult        string
		mockValidateError error
		mockFileError     error
		mockOperationErr  error
		mockRunOpError    error
		want              string
		wantErr           bool
		expectedError     error
	}{
		{
			name:      "successfully process sum operation",
			operation: "sum",
			filePath:  "testdata/matrix1.csv",
			mockFileContent: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "2", "3"},
					{"4", "5", "6"},
				},
			},
			mockMatrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2, 3},
					{4, 5, 6},
				},
			},
			mockResult:        "21",
			mockValidateError: nil,
			mockFileError:     nil,
			mockOperationErr:  nil,
			mockRunOpError:    nil,
			want:              "21",
			wantErr:           false,
		},
		{
			name:      "successfully process multiply operation",
			operation: "multiply",
			filePath:  "testdata/matrix1.csv",
			mockFileContent: &repository.MatrixFileContent{
				Content: [][]string{
					{"2", "3"},
					{"4", "5"},
				},
			},
			mockMatrix: &entity.Matrix{
				Data: [][]int64{
					{2, 3},
					{4, 5},
				},
			},
			mockResult: "120",
			want:       "120",
			wantErr:    false,
		},
		{
			name:      "successfully process echo operation",
			operation: "echo",
			filePath:  "testdata/matrix1.csv",
			mockFileContent: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "2"},
					{"3", "4"},
				},
			},
			mockMatrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2},
					{3, 4},
				},
			},
			mockResult: "1,2\n3,4",
			want:       "1,2\n3,4",
			wantErr:    false,
		},
		{
			name:          "fail when operation is empty",
			operation:     "",
			filePath:      "testdata/matrix1.csv",
			want:          "",
			wantErr:       true,
			expectedError: apperrors.ErrInvalidInput,
		},
		{
			name:              "fail when file path is invalid - path traversal",
			operation:         "sum",
			filePath:          "../secret.csv",
			mockValidateError: apperrors.ErrInvalidInput,
			want:              "",
			wantErr:           true,
			expectedError:     apperrors.ErrInvalidInput,
		},
		{
			name:              "fail when file path is empty",
			operation:         "sum",
			filePath:          "",
			mockValidateError: apperrors.ErrInvalidInput,
			want:              "",
			wantErr:           true,
			expectedError:     apperrors.ErrInvalidInput,
		},
		{
			name:             "fail when operation is invalid",
			operation:        "divide",
			filePath:         "testdata/matrix1.csv",
			mockOperationErr: apperrors.ErrInvalidInput,
			want:             "",
			wantErr:          true,
			expectedError:    apperrors.ErrInvalidInput,
		},
		{
			name:          "fail when file not found",
			operation:     "sum",
			filePath:      "testdata/notfound.csv",
			mockFileError: apperrors.ErrNotFound,
			want:          "",
			wantErr:       true,
			expectedError: apperrors.ErrNotFound,
		},
		{
			name:      "fail when matrix validation fails",
			operation: "sum",
			filePath:  "testdata/matrix2.csv",
			mockFileContent: &repository.MatrixFileContent{
				Content: [][]string{
					{"a", "2", "3"},
					{"4", "b", "6"},
				},
			},
			mockMatrix:        nil,
			mockValidateError: nil,
			mockFileError:     nil,
			mockOperationErr:  nil,
			want:              "",
			wantErr:           true,
			expectedError:     apperrors.ErrUnprocessableEntity,
		},
		{
			name:      "fail when operation execution fails",
			operation: "sum",
			filePath:  "testdata/matrix6.csv",
			mockFileContent: &repository.MatrixFileContent{
				Content: [][]string{},
			},
			mockMatrix: &entity.Matrix{
				Data: [][]int64{},
			},
			mockValidateError: nil,
			mockFileError:     nil,
			mockOperationErr:  nil,
			mockRunOpError:    apperrors.ErrInvalidInput,
			want:              "",
			wantErr:           true,
			expectedError:     apperrors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockRepo := &mockMatrixRepository{}
			mockValidator := NewMockMatrixValidatorDomainInterface(t)
			mockOperations := NewMockMatrixOperationsDomainInterface(t)

			// Setup expectations based on test case
			if tt.operation != "" {
				mockValidator.On("ValidateFilePath", mock.Anything, tt.filePath).
					Return(tt.mockValidateError)
			}

			if tt.mockValidateError == nil && tt.operation != "" {
				mockOperations.On("IsValidOperation", mock.Anything, tt.operation).
					Return(tt.mockOperationErr)
			}

			if tt.mockOperationErr == nil && tt.mockValidateError == nil && tt.operation != "" {
				mockRepo.On("GetFileContent", mock.Anything, tt.filePath).
					Return(tt.mockFileContent, tt.mockFileError)
			}

			if tt.mockFileError == nil && tt.mockOperationErr == nil && tt.mockValidateError == nil && tt.operation != "" {
				var validateErr error
				if tt.mockMatrix == nil {
					validateErr = apperrors.ErrUnprocessableEntity
				}
				mockValidator.On("Validate", mock.Anything, tt.mockFileContent).
					Return(tt.mockMatrix, validateErr)

				if tt.mockMatrix != nil {
					mockOperations.On("RunOperation", mock.Anything, tt.mockMatrix, tt.operation).
						Return(tt.mockResult, tt.mockRunOpError)
				}
			}

			// Create domain with mocks
			domain := &matrixDomain{
				matrixRepository: mockRepo,
				validatorDomain:  mockValidator,
				operationsDomain: mockOperations,
			}

			// Execute
			got, err := domain.ProcessMatrix(context.Background(), tt.operation, tt.filePath)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMatrixDomain_ProcessMatrix_ContextCancellation(t *testing.T) {
	tests := []struct {
		name        string
		setupCtx    func() context.Context
		wantErr     bool
		expectedErr error
	}{
		{
			name: "context cancelled before processing",
			setupCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately
				return ctx
			},
			wantErr:     true,
			expectedErr: context.Canceled,
		},
		{
			name: "context with valid deadline",
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()

			if !tt.wantErr {
				// Setup mocks for successful case
				mockRepo := &mockMatrixRepository{}
				mockValidator := NewMockMatrixValidatorDomainInterface(t)
				mockOperations := NewMockMatrixOperationsDomainInterface(t)

				mockValidator.On("ValidateFilePath", mock.Anything, "testdata/matrix1.csv").Return(nil)
				mockOperations.On("IsValidOperation", mock.Anything, "sum").Return(nil)
				mockRepo.On("GetFileContent", mock.Anything, "testdata/matrix1.csv").Return(
					&repository.MatrixFileContent{Content: [][]string{{"1", "2"}}},
					nil,
				)
				mockValidator.On("Validate", mock.Anything, mock.Anything).Return(
					&entity.Matrix{Data: [][]int64{{1, 2}}},
					nil,
				)
				mockOperations.On("RunOperation", mock.Anything, mock.Anything, "sum").Return("3", nil)

				domain := &matrixDomain{
					matrixRepository: mockRepo,
					validatorDomain:  mockValidator,
					operationsDomain: mockOperations,
				}

				got, err := domain.ProcessMatrix(ctx, "sum", "testdata/matrix1.csv")
				assert.NoError(t, err)
				assert.Equal(t, "3", got)
			} else {
				domain := &matrixDomain{}
				got, err := domain.ProcessMatrix(ctx, "sum", "testdata/matrix1.csv")

				assert.Error(t, err)
				assert.Equal(t, "", got)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
			}
		})
	}
}

func TestMatrixDomain_ProcessMatrix_ErrorPropagation(t *testing.T) {
	t.Run("error from validator is properly wrapped", func(t *testing.T) {
		mockValidator := NewMockMatrixValidatorDomainInterface(t)
		mockValidator.On("ValidateFilePath", mock.Anything, "invalid/path").
			Return(errors.New("custom validation error"))

		domain := &matrixDomain{
			validatorDomain: mockValidator,
		}

		_, err := domain.ProcessMatrix(context.Background(), "sum", "invalid/path")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "custom validation error")
	})

	t.Run("error from repository is properly wrapped", func(t *testing.T) {
		mockValidator := NewMockMatrixValidatorDomainInterface(t)
		mockOperations := NewMockMatrixOperationsDomainInterface(t)
		mockRepo := &mockMatrixRepository{}

		mockValidator.On("ValidateFilePath", mock.Anything, "testdata/matrix1.csv").Return(nil)
		mockOperations.On("IsValidOperation", mock.Anything, "sum").Return(nil)
		mockRepo.On("GetFileContent", mock.Anything, "testdata/matrix1.csv").
			Return(nil, errors.New("file read error"))

		domain := &matrixDomain{
			matrixRepository: mockRepo,
			validatorDomain:  mockValidator,
			operationsDomain: mockOperations,
		}

		_, err := domain.ProcessMatrix(context.Background(), "sum", "testdata/matrix1.csv")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file read error")
	})
}
