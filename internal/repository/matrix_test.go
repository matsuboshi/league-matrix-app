package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	apperrors "github.com/matsuboshi/league-matrix-app/pkg/errors"
)

func TestMatrixRepository_GetFileContent(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		wantContent *MatrixFileContent
		wantErr     bool
		errType     error
	}{
		{
			name:     "successfully read matrix1.csv",
			filePath: "../../testdata/matrix1.csv",
			wantContent: &MatrixFileContent{
				Content: [][]string{
					{"1", "2", "3"},
					{"4", "5", "6"},
					{"7", "8", "9"},
					{"10", "11", "12"},
					{"13", "14", "15"},
					{"16", "17", "18"},
					{"19", "20", "21"},
					{"22", "23", "24"},
					{"25", "26", "27"},
				},
			},
			wantErr: false,
		},
		{
			name:     "successfully read matrix0.csv - large values",
			filePath: "../../testdata/matrix0.csv",
			wantContent: &MatrixFileContent{
				Content: [][]string{
					{"1000000", "1000000", "1000000"},
					{"1000000", "1000000", "1000000"},
					{"1000000", "1000000", "1000000"},
					{"1000000", "1000000", "1000000"},
					{"1000000", "1000000", "1000000"},
					{"1000000", "1000000", "1000000"},
					{"1000000", "1000000", "1000000"},
					{"1000000", "1000000", "1000000"},
					{"1000000", "1000000", "1000000"},
					{"1000000", "1000000", "1000000"},
				},
			},
			wantErr: false,
		},
		{
			name:     "file not found",
			filePath: "../../testdata/nonexistent.csv",
			wantErr:  true,
			errType:  apperrors.ErrNotFound,
		},
		{
			name:     "file too large",
			filePath: "../../testdata/gopher.jpg.csv",
			wantErr:  true,
			errType:  apperrors.ErrPayloadTooLarge,
		},
		{
			name:     "invalid CSV format in matrix2.csv",
			filePath: "../../testdata/matrix2.csv",
			wantContent: &MatrixFileContent{
				Content: [][]string{
					{"a", "2", "3"},
					{"4", "b", "6"},
					{"7", "8", "9"},
				},
			},
			wantErr: false, // CSV parsing succeeds, validation happens later
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMatrixRepository()

			got, err := repo.GetFileContent(context.Background(), tt.filePath)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantContent, got)
			}
		})
	}
}

func TestMatrixRepository_GetFileContent_ContextCancellation(t *testing.T) {
	t.Run("context cancelled before reading", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		repo := NewMatrixRepository()
		got, err := repo.GetFileContent(ctx, "../../testdata/matrix1.csv")

		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, context.Canceled)
	})
}

func TestMatrixRepository_GetFileContent_FileSize(t *testing.T) {
	// Create a temporary test file that exceeds the size limit
	t.Run("reject file larger than 1KB", func(t *testing.T) {
		// Create temp directory
		tmpDir := t.TempDir()
		largeFile := filepath.Join(tmpDir, "large.csv")

		// Create a file larger than 1KB
		content := make([]byte, 1025) // 1025 bytes > 1KB
		for i := range content {
			content[i] = 'a'
		}
		err := os.WriteFile(largeFile, content, 0o644)
		assert.NoError(t, err)

		repo := NewMatrixRepository()
		got, err := repo.GetFileContent(context.Background(), largeFile)

		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, apperrors.ErrPayloadTooLarge)
	})

	t.Run("accept file exactly 1KB", func(t *testing.T) {
		// Create temp directory
		tmpDir := t.TempDir()
		exactFile := filepath.Join(tmpDir, "exact.csv")

		// Create a file exactly 1KB with valid CSV content
		var content string
		for len(content) < 1024 {
			content += "1,2,3\n"
		}
		content = content[:1024] // Trim to exactly 1KB
		err := os.WriteFile(exactFile, []byte(content), 0o644)
		assert.NoError(t, err)

		repo := NewMatrixRepository()
		got, err := repo.GetFileContent(context.Background(), exactFile)

		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("accept file smaller than 1KB", func(t *testing.T) {
		// Create temp directory
		tmpDir := t.TempDir()
		smallFile := filepath.Join(tmpDir, "small.csv")

		// Create a small file
		content := "1,2,3\n4,5,6\n"
		err := os.WriteFile(smallFile, []byte(content), 0o644)
		assert.NoError(t, err)

		repo := NewMatrixRepository()
		got, err := repo.GetFileContent(context.Background(), smallFile)

		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, 2, len(got.Content)) // 2 rows
	})
}

func TestMatrixRepository_GetFileContent_EdgeCases(t *testing.T) {
	t.Run("empty CSV file", func(t *testing.T) {
		tmpDir := t.TempDir()
		emptyFile := filepath.Join(tmpDir, "empty.csv")
		err := os.WriteFile(emptyFile, []byte(""), 0o644)
		assert.NoError(t, err)

		repo := NewMatrixRepository()
		got, err := repo.GetFileContent(context.Background(), emptyFile)

		// Empty file should be parsed successfully (will fail validation later)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, 0, len(got.Content))
	})

	t.Run("CSV with single value", func(t *testing.T) {
		tmpDir := t.TempDir()
		singleFile := filepath.Join(tmpDir, "single.csv")
		err := os.WriteFile(singleFile, []byte("42"), 0o644)
		assert.NoError(t, err)

		repo := NewMatrixRepository()
		got, err := repo.GetFileContent(context.Background(), singleFile)

		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, 1, len(got.Content))
		assert.Equal(t, []string{"42"}, got.Content[0])
	})

	t.Run("CSV with trailing newline", func(t *testing.T) {
		tmpDir := t.TempDir()
		trailingFile := filepath.Join(tmpDir, "trailing.csv")
		err := os.WriteFile(trailingFile, []byte("1,2,3\n4,5,6\n"), 0o644)
		assert.NoError(t, err)

		repo := NewMatrixRepository()
		got, err := repo.GetFileContent(context.Background(), trailingFile)

		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, 2, len(got.Content))
	})
}
