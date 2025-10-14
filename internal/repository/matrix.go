package repository

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"

	apperrors "github.com/matsuboshi/league-matrix-app/pkg/errors"
)

const (
	// maxFileSizeBytes defines the maximum allowed file size in bytes (1KB).
	// This prevents denial of service attacks from extremely large files.
	// Maximum theoretical size for 10x10 matrix with 7-digit numbers is ~800 bytes.
	maxFileSizeBytes = 1024 // 1KB
)

// MatrixRepositoryInterface defines the contract for accessing matrix data from external sources.
type MatrixRepositoryInterface interface {
	// GetFileContent reads and parses a CSV file containing matrix data.
	// It returns the raw string content of the file organized as a 2D slice.
	GetFileContent(ctx context.Context, filePath string) (*MatrixFileContent, error)
}

// MatrixFileContent represents the raw content read from a matrix file.
// Content contains the parsed CSV data as rows and columns of strings.
type MatrixFileContent struct {
	Content [][]string
}

type matrixRepository struct{}

// NewMatrixRepository creates a new instance of MatrixRepositoryInterface.
// It returns a repository implementation that can read matrix data from CSV files.
func NewMatrixRepository() MatrixRepositoryInterface {
	return &matrixRepository{}
}

func (r *matrixRepository) GetFileContent(ctx context.Context, filePath string) (*MatrixFileContent, error) {
	// Check if context is already cancelled
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("failed to open file",
			"file_path", filePath,
			"error", err)
		return nil, fmt.Errorf("%w: failed to open file: %v", apperrors.ErrNotFound, err)
	}
	defer file.Close()

	// Get file info to check size
	fileInfo, err := file.Stat()
	if err != nil {
		slog.Error("failed to get file info",
			"file_path", filePath,
			"error", err)
		return nil, fmt.Errorf("%w: failed to get file info: %v", apperrors.ErrNotFound, err)
	}

	// Check file size BEFORE reading to prevent DoS attacks
	if fileInfo.Size() > maxFileSizeBytes {
		return nil, fmt.Errorf("%w: file too large: %d bytes (maximum: %d bytes)",
			apperrors.ErrPayloadTooLarge, fileInfo.Size(), maxFileSizeBytes)
	}

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		slog.Error("failed to parse CSV",
			"file_path", filePath,
			"error", err)
		return nil, fmt.Errorf("%w: failed to read CSV file: %v", apperrors.ErrUnprocessableEntity, err)
	}

	// Return the matrix file content
	return &MatrixFileContent{
		Content: records,
	}, nil
}
