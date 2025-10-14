package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matsuboshi/league-matrix-app/internal/entity"
	"github.com/matsuboshi/league-matrix-app/internal/repository"
	apperrors "github.com/matsuboshi/league-matrix-app/pkg/errors"
)

func TestMatrixValidatorDomain_ValidateFilePath(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		wantErr  bool
		errType  error
	}{
		{
			name:     "valid file path",
			filePath: "testdata/matrix1.csv",
			wantErr:  false,
		},
		{
			name:     "valid file path with different name",
			filePath: "testdata/matrix2.csv",
			wantErr:  false,
		},
		{
			name:     "empty file path",
			filePath: "",
			wantErr:  true,
			errType:  apperrors.ErrInvalidInput,
		},
		{
			name:     "path traversal attempt - parent directory",
			filePath: "../secret.csv",
			wantErr:  true,
			errType:  apperrors.ErrInvalidInput,
		},
		{
			name:     "path traversal attempt - within testdata",
			filePath: "testdata/../secret.csv",
			wantErr:  true,
			errType:  apperrors.ErrInvalidInput,
		},
		{
			name:     "path traversal attempt - multiple levels",
			filePath: "testdata/../../etc/passwd",
			wantErr:  true,
			errType:  apperrors.ErrInvalidInput,
		},
		{
			name:     "file not in testdata directory",
			filePath: "data/matrix.csv",
			wantErr:  true,
			errType:  apperrors.ErrInvalidInput,
		},
		{
			name:     "file in root directory",
			filePath: "matrix.csv",
			wantErr:  true,
			errType:  apperrors.ErrInvalidInput,
		},
		{
			name:     "non-csv file extension",
			filePath: "testdata/matrix.txt",
			wantErr:  true,
			errType:  apperrors.ErrInvalidInput,
		},
		{
			name:     "non-csv file extension - json",
			filePath: "testdata/matrix.json",
			wantErr:  true,
			errType:  apperrors.ErrInvalidInput,
		},
		{
			name:     "file without extension",
			filePath: "testdata/matrix",
			wantErr:  true,
			errType:  apperrors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewMatrixValidatorDomain()

			err := validator.ValidateFilePath(context.Background(), tt.filePath)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMatrixValidatorDomain_Validate(t *testing.T) {
	tests := []struct {
		name       string
		rawData    *repository.MatrixFileContent
		wantMatrix *entity.Matrix
		wantErr    bool
		errType    error
	}{
		{
			name: "valid 2x2 matrix",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "2"},
					{"3", "4"},
				},
			},
			wantMatrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2},
					{3, 4},
				},
			},
			wantErr: false,
		},
		{
			name: "valid 3x3 matrix",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "2", "3"},
					{"4", "5", "6"},
					{"7", "8", "9"},
				},
			},
			wantMatrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			},
			wantErr: false,
		},
		{
			name: "valid matrix with negative numbers",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"-1", "-2"},
					{"-3", "-4"},
				},
			},
			wantMatrix: &entity.Matrix{
				Data: [][]int64{
					{-1, -2},
					{-3, -4},
				},
			},
			wantErr: false,
		},
		{
			name: "valid matrix with large numbers",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"1000000", "2000000"},
					{"3000000", "4000000"},
				},
			},
			wantMatrix: &entity.Matrix{
				Data: [][]int64{
					{1000000, 2000000},
					{3000000, 4000000},
				},
			},
			wantErr: false,
		},
		{
			name: "valid single element matrix",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{{"42"}},
			},
			wantMatrix: &entity.Matrix{
				Data: [][]int64{{42}},
			},
			wantErr: false,
		},
		{
			name: "valid 10x10 matrix - maximum size",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
				},
			},
			wantMatrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
					{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				},
			},
			wantErr: false,
		},
		{
			name:    "nil raw data",
			rawData: nil,
			wantErr: true,
			errType: apperrors.ErrUnprocessableEntity,
		},
		{
			name: "empty matrix data",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{},
			},
			wantErr: true,
			errType: apperrors.ErrUnprocessableEntity,
		},
		{
			name: "matrix with non-integer value - letter",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "a", "3"},
					{"4", "5", "6"},
				},
			},
			wantErr: true,
			errType: apperrors.ErrUnprocessableEntity,
		},
		{
			name: "matrix with inconsistent row lengths",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "2", "3"},
					{"4", "5"},
				},
			},
			wantErr: true,
			errType: apperrors.ErrUnprocessableEntity,
		},
		{
			name: "matrix with inconsistent row lengths - shorter first row",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "2"},
					{"3", "4", "5"},
				},
			},
			wantErr: true,
			errType: apperrors.ErrUnprocessableEntity,
		},
		{
			name: "matrix exceeds maximum rows (11 rows)",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "2"},
					{"1", "2"},
					{"1", "2"},
					{"1", "2"},
					{"1", "2"},
					{"1", "2"},
					{"1", "2"},
					{"1", "2"},
					{"1", "2"},
					{"1", "2"},
					{"1", "2"},
				},
			},
			wantErr: true,
			errType: apperrors.ErrUnprocessableEntity,
		},
		{
			name: "matrix exceeds maximum columns (11 columns)",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
					{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
				},
			},
			wantErr: true,
			errType: apperrors.ErrUnprocessableEntity,
		},
		{
			name: "matrix with empty string",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "", "3"},
					{"4", "5", "6"},
				},
			},
			wantErr: true,
			errType: apperrors.ErrUnprocessableEntity,
		},
		{
			name: "matrix with whitespace only",
			rawData: &repository.MatrixFileContent{
				Content: [][]string{
					{"1", "  ", "3"},
					{"4", "5", "6"},
				},
			},
			wantErr: true,
			errType: apperrors.ErrUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewMatrixValidatorDomain()

			gotMatrix, err := validator.Validate(context.Background(), tt.rawData)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				assert.Nil(t, gotMatrix)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotMatrix)
				assert.Equal(t, tt.wantMatrix, gotMatrix)
			}
		})
	}
}

func TestMatrixValidatorDomain_ContextCancellation(t *testing.T) {
	tests := []struct {
		name     string
		setupCtx func() context.Context
		wantErr  bool
	}{
		{
			name: "context cancelled before ValidateFilePath",
			setupCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
		{
			name: "context cancelled before Validate",
			setupCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewMatrixValidatorDomain()
			ctx := tt.setupCtx()

			if tt.name == "context cancelled before ValidateFilePath" {
				err := validator.ValidateFilePath(ctx, "testdata/matrix1.csv")
				if tt.wantErr {
					assert.Error(t, err)
					assert.ErrorIs(t, err, context.Canceled)
				}
			}

			if tt.name == "context cancelled before Validate" {
				rawData := &repository.MatrixFileContent{
					Content: [][]string{{"1", "2"}},
				}
				_, err := validator.Validate(ctx, rawData)
				if tt.wantErr {
					assert.Error(t, err)
					assert.ErrorIs(t, err, context.Canceled)
				}
			}
		})
	}
}
