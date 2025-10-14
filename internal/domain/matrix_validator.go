package domain

import (
	"context"
	"fmt"
	"strings"

	"github.com/matsuboshi/league-matrix-app/internal/entity"
	"github.com/matsuboshi/league-matrix-app/internal/repository"
	apperrors "github.com/matsuboshi/league-matrix-app/pkg/errors"
)

const (
	maxInputMatrixRows = 10
	maxInputMatrixCols = 10
)

// MatrixValidatorDomainInterface defines the contract for validating and transforming raw matrix data.
// It ensures matrix data integrity and converts string data to typed entities.
type MatrixValidatorDomainInterface interface {
	ValidateFilePath(ctx context.Context, filePath string) error

	// Validate checks raw matrix file content for consistency and converts it to a typed Matrix entity.
	// It ensures all rows have equal length and all values are valid integers.
	// Returns a validated Matrix entity or an error if validation fails.
	Validate(ctx context.Context, matrix *repository.MatrixFileContent) (*entity.Matrix, error)
}

type matrixValidatorDomain struct{}

// NewMatrixValidatorDomain creates a new instance of MatrixValidatorDomainInterface.
// It returns a validator that can transform and validate raw matrix data.
func NewMatrixValidatorDomain() MatrixValidatorDomainInterface {
	return &matrixValidatorDomain{}
}

func (d *matrixValidatorDomain) ValidateFilePath(ctx context.Context, filePath string) error {
	// Check if context is already cancelled
	if err := ctx.Err(); err != nil {
		return err
	}

	if filePath == "" {
		return fmt.Errorf("%w: file parameter is required", apperrors.ErrInvalidInput)
	}
	if strings.Contains(filePath, "..") {
		return fmt.Errorf("%w: path traversal not allowed", apperrors.ErrInvalidInput)
	}
	if !strings.HasPrefix(filePath, "testdata/") {
		return fmt.Errorf("%w: only files in testdata/ are allowed", apperrors.ErrInvalidInput)
	}
	if !strings.HasSuffix(filePath, ".csv") {
		return fmt.Errorf("%w: only .csv files are supported", apperrors.ErrInvalidInput)
	}
	return nil
}

func (d *matrixValidatorDomain) Validate(ctx context.Context, rawData *repository.MatrixFileContent) (*entity.Matrix, error) {
	// Check if context is already cancelled
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if rawData == nil || len(rawData.Content) == 0 {
		return nil, fmt.Errorf("%w: empty matrix data", apperrors.ErrUnprocessableEntity)
	}

	rows := len(rawData.Content)
	cols := len(rawData.Content[0])

	// Validate maximum dimensions
	if rows > maxInputMatrixRows {
		return nil, fmt.Errorf("%w: matrix exceeds maximum row limit: got %d rows, maximum is %d",
			apperrors.ErrUnprocessableEntity, rows, maxInputMatrixRows)
	}

	if cols > maxInputMatrixCols {
		return nil, fmt.Errorf("%w: matrix exceeds maximum column limit: got %d columns, maximum is %d",
			apperrors.ErrUnprocessableEntity, cols, maxInputMatrixCols)
	}

	// Validate that all rows have the same number of columns
	for i, row := range rawData.Content {
		if len(row) != cols {
			return nil, fmt.Errorf("%w: inconsistent row length at row %d: expected %d columns, got %d",
				apperrors.ErrUnprocessableEntity, i, cols, len(row))
		}
	}

	// Convert string data to int64
	matrix := &entity.Matrix{
		Data: make([][]int64, rows),
	}

	for i, row := range rawData.Content {
		matrix.Data[i] = make([]int64, cols)
		for j, val := range row {
			var num int64
			_, err := fmt.Sscanf(val, "%d", &num)
			if err != nil {
				return nil, fmt.Errorf("%w: invalid integer value at row %d, column %d: %v",
					apperrors.ErrUnprocessableEntity, i, j, err)
			}
			matrix.Data[i][j] = num
		}
	}

	return matrix, nil
}
