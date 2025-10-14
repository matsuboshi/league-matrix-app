package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matsuboshi/league-matrix-app/internal/entity"
	apperrors "github.com/matsuboshi/league-matrix-app/pkg/errors"
)

func TestMatrixOperationsDomain_ListOperations(t *testing.T) {
	domain := NewMatrixOperationsDomain()

	operations := domain.ListOperations()

	assert.NotEmpty(t, operations)
	assert.Contains(t, operations, "sum")
	assert.Contains(t, operations, "multiply")
	assert.Contains(t, operations, "echo")
	assert.Contains(t, operations, "invert")
	assert.Contains(t, operations, "flatten")
	assert.Len(t, operations, 5)
}

func TestMatrixOperationsDomain_IsValidOperation(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		wantErr   bool
		errType   error
	}{
		{
			name:      "valid operation - sum",
			operation: "sum",
			wantErr:   false,
		},
		{
			name:      "valid operation - multiply",
			operation: "multiply",
			wantErr:   false,
		},
		{
			name:      "valid operation - echo",
			operation: "echo",
			wantErr:   false,
		},
		{
			name:      "valid operation - invert",
			operation: "invert",
			wantErr:   false,
		},
		{
			name:      "valid operation - flatten",
			operation: "flatten",
			wantErr:   false,
		},
		{
			name:      "invalid operation - divide",
			operation: "divide",
			wantErr:   true,
			errType:   apperrors.ErrInvalidInput,
		},
		{
			name:      "invalid operation - subtract",
			operation: "subtract",
			wantErr:   true,
			errType:   apperrors.ErrInvalidInput,
		},
		{
			name:      "invalid operation - empty",
			operation: "",
			wantErr:   true,
			errType:   apperrors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := NewMatrixOperationsDomain()

			err := domain.IsValidOperation(context.Background(), tt.operation)

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

func TestMatrixOperationsDomain_Sum(t *testing.T) {
	tests := []struct {
		name    string
		matrix  *entity.Matrix
		want    string
		wantErr bool
		errType error
	}{
		{
			name: "sum of 2x2 matrix",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2},
					{3, 4},
				},
			},
			want:    "10",
			wantErr: false,
		},
		{
			name: "sum of 3x3 matrix",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			},
			want:    "45",
			wantErr: false,
		},
		{
			name: "sum with negative numbers",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{-1, -2},
					{-3, -4},
				},
			},
			want:    "-10",
			wantErr: false,
		},
		{
			name: "sum with large numbers",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1000000, 1000000},
					{1000000, 1000000},
				},
			},
			want:    "4000000",
			wantErr: false,
		},
		{
			name: "sum of single element",
			matrix: &entity.Matrix{
				Data: [][]int64{{42}},
			},
			want:    "42",
			wantErr: false,
		},
		{
			name:    "empty matrix",
			matrix:  &entity.Matrix{Data: [][]int64{}},
			want:    "",
			wantErr: true,
			errType: apperrors.ErrInvalidInput,
		},
		{
			name:    "nil matrix",
			matrix:  nil,
			want:    "",
			wantErr: true,
			errType: apperrors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := &matrixOperationsDomain{}

			got, err := domain.sum(tt.matrix)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMatrixOperationsDomain_Multiply(t *testing.T) {
	tests := []struct {
		name    string
		matrix  *entity.Matrix
		want    string
		wantErr bool
		errType error
	}{
		{
			name: "multiply 2x2 matrix",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{2, 3},
					{4, 5},
				},
			},
			want:    "120",
			wantErr: false,
		},
		{
			name: "multiply 3x3 matrix",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			},
			want:    "362880",
			wantErr: false,
		},
		{
			name: "multiply with zero",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 0, 3},
					{4, 5, 6},
				},
			},
			want:    "0",
			wantErr: false,
		},
		{
			name: "multiply with negative numbers",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{-2, 3},
					{4, -5},
				},
			},
			want:    "120",
			wantErr: false,
		},
		{
			name: "multiply single element",
			matrix: &entity.Matrix{
				Data: [][]int64{{7}},
			},
			want:    "7",
			wantErr: false,
		},
		{
			name:    "empty matrix",
			matrix:  &entity.Matrix{Data: [][]int64{}},
			want:    "",
			wantErr: true,
			errType: apperrors.ErrInvalidInput,
		},
		{
			name:    "nil matrix",
			matrix:  nil,
			want:    "",
			wantErr: true,
			errType: apperrors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := &matrixOperationsDomain{}

			got, err := domain.multiply(tt.matrix)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMatrixOperationsDomain_Echo(t *testing.T) {
	tests := []struct {
		name    string
		matrix  *entity.Matrix
		want    string
		wantErr bool
		errType error
	}{
		{
			name: "echo 2x2 matrix",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2},
					{3, 4},
				},
			},
			want:    "1,2\n3,4",
			wantErr: false,
		},
		{
			name: "echo 3x3 matrix",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			},
			want:    "1,2,3\n4,5,6\n7,8,9",
			wantErr: false,
		},
		{
			name: "echo single element",
			matrix: &entity.Matrix{
				Data: [][]int64{{42}},
			},
			want:    "42",
			wantErr: false,
		},
		{
			name: "echo with negative numbers",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{-1, -2},
					{-3, -4},
				},
			},
			want:    "-1,-2\n-3,-4",
			wantErr: false,
		},
		{
			name:    "empty matrix",
			matrix:  &entity.Matrix{Data: [][]int64{}},
			want:    "",
			wantErr: true,
			errType: apperrors.ErrInvalidInput,
		},
		{
			name:    "nil matrix",
			matrix:  nil,
			want:    "",
			wantErr: true,
			errType: apperrors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := &matrixOperationsDomain{}

			got, err := domain.echo(tt.matrix)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMatrixOperationsDomain_Invert(t *testing.T) {
	tests := []struct {
		name    string
		matrix  *entity.Matrix
		want    string
		wantErr bool
		errType error
	}{
		{
			name: "invert 2x2 matrix",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2},
					{3, 4},
				},
			},
			want:    "1,3\n2,4",
			wantErr: false,
		},
		{
			name: "invert 3x3 matrix",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			},
			want:    "1,4,7\n2,5,8\n3,6,9",
			wantErr: false,
		},
		{
			name: "invert rectangular matrix 2x3",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2, 3},
					{4, 5, 6},
				},
			},
			want:    "1,4\n2,5\n3,6",
			wantErr: false,
		},
		{
			name: "invert rectangular matrix 3x2",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2},
					{3, 4},
					{5, 6},
				},
			},
			want:    "1,3,5\n2,4,6",
			wantErr: false,
		},
		{
			name: "invert single element",
			matrix: &entity.Matrix{
				Data: [][]int64{{42}},
			},
			want:    "42",
			wantErr: false,
		},
		{
			name: "invert single row",
			matrix: &entity.Matrix{
				Data: [][]int64{{1, 2, 3, 4}},
			},
			want:    "1\n2\n3\n4",
			wantErr: false,
		},
		{
			name: "invert single column",
			matrix: &entity.Matrix{
				Data: [][]int64{{1}, {2}, {3}, {4}},
			},
			want:    "1,2,3,4",
			wantErr: false,
		},
		{
			name:    "empty matrix",
			matrix:  &entity.Matrix{Data: [][]int64{}},
			want:    "",
			wantErr: true,
			errType: apperrors.ErrInvalidInput,
		},
		{
			name:    "nil matrix",
			matrix:  nil,
			want:    "",
			wantErr: true,
			errType: apperrors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := &matrixOperationsDomain{}

			got, err := domain.invert(tt.matrix)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMatrixOperationsDomain_Flatten(t *testing.T) {
	tests := []struct {
		name    string
		matrix  *entity.Matrix
		want    string
		wantErr bool
		errType error
	}{
		{
			name: "flatten 2x2 matrix",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2},
					{3, 4},
				},
			},
			want:    "1,2,3,4",
			wantErr: false,
		},
		{
			name: "flatten 3x3 matrix",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			},
			want:    "1,2,3,4,5,6,7,8,9",
			wantErr: false,
		},
		{
			name: "flatten single element",
			matrix: &entity.Matrix{
				Data: [][]int64{{42}},
			},
			want:    "42",
			wantErr: false,
		},
		{
			name: "flatten single row",
			matrix: &entity.Matrix{
				Data: [][]int64{{1, 2, 3, 4, 5}},
			},
			want:    "1,2,3,4,5",
			wantErr: false,
		},
		{
			name: "flatten with negative numbers",
			matrix: &entity.Matrix{
				Data: [][]int64{
					{-1, -2},
					{-3, -4},
				},
			},
			want:    "-1,-2,-3,-4",
			wantErr: false,
		},
		{
			name:    "empty matrix",
			matrix:  &entity.Matrix{Data: [][]int64{}},
			want:    "",
			wantErr: true,
			errType: apperrors.ErrInvalidInput,
		},
		{
			name:    "nil matrix",
			matrix:  nil,
			want:    "",
			wantErr: true,
			errType: apperrors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := &matrixOperationsDomain{}

			got, err := domain.flatten(tt.matrix)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMatrixOperationsDomain_RunOperation(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		matrix    *entity.Matrix
		want      string
		wantErr   bool
		errType   error
	}{
		{
			name:      "run sum operation",
			operation: "sum",
			matrix: &entity.Matrix{
				Data: [][]int64{{1, 2}, {3, 4}},
			},
			want:    "10",
			wantErr: false,
		},
		{
			name:      "run multiply operation",
			operation: "multiply",
			matrix: &entity.Matrix{
				Data: [][]int64{{2, 3}, {4, 5}},
			},
			want:    "120",
			wantErr: false,
		},
		{
			name:      "run echo operation",
			operation: "echo",
			matrix: &entity.Matrix{
				Data: [][]int64{{1, 2}, {3, 4}},
			},
			want:    "1,2\n3,4",
			wantErr: false,
		},
		{
			name:      "run invert operation",
			operation: "invert",
			matrix: &entity.Matrix{
				Data: [][]int64{{1, 2}, {3, 4}},
			},
			want:    "1,3\n2,4",
			wantErr: false,
		},
		{
			name:      "run flatten operation",
			operation: "flatten",
			matrix: &entity.Matrix{
				Data: [][]int64{{1, 2}, {3, 4}},
			},
			want:    "1,2,3,4",
			wantErr: false,
		},
		{
			name:      "unsupported operation",
			operation: "unsupported",
			matrix: &entity.Matrix{
				Data: [][]int64{{1, 2}},
			},
			want:    "",
			wantErr: true,
			errType: apperrors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := NewMatrixOperationsDomain()

			got, err := domain.RunOperation(context.Background(), tt.matrix, tt.operation)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMatrixOperationsDomain_ContextCancellation(t *testing.T) {
	tests := []struct {
		name      string
		setupCtx  func() context.Context
		operation string
		wantErr   bool
	}{
		{
			name: "context cancelled before IsValidOperation",
			setupCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			operation: "sum",
			wantErr:   true,
		},
		{
			name: "context cancelled before RunOperation",
			setupCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			operation: "sum",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := NewMatrixOperationsDomain()
			ctx := tt.setupCtx()

			if tt.name == "context cancelled before IsValidOperation" {
				err := domain.IsValidOperation(ctx, tt.operation)
				if tt.wantErr {
					assert.Error(t, err)
					assert.ErrorIs(t, err, context.Canceled)
				}
			}

			if tt.name == "context cancelled before RunOperation" {
				matrix := &entity.Matrix{Data: [][]int64{{1, 2}}}
				_, err := domain.RunOperation(ctx, matrix, tt.operation)
				if tt.wantErr {
					assert.Error(t, err)
					assert.ErrorIs(t, err, context.Canceled)
				}
			}
		})
	}
}
