package domain

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/matsuboshi/league-matrix-app/internal/entity"
	apperrors "github.com/matsuboshi/league-matrix-app/pkg/errors"
)

type Operation string

const (
	SumOperation      Operation = "sum"
	MultiplyOperation Operation = "multiply"
	EchoOperation     Operation = "echo"
	InvertOperation   Operation = "invert"
	FlattenOperation  Operation = "flatten"
)

var matrixOperations = map[Operation]bool{
	SumOperation:      true,
	MultiplyOperation: true,
	EchoOperation:     true,
	InvertOperation:   true,
	FlattenOperation:  true,
}

// MatrixOperationsDomainInterface defines the contract for performing operations on matrices.
// It provides methods to list, validate, and execute various matrix transformations and calculations.
type MatrixOperationsDomainInterface interface {
	// ListOperations returns a list of all supported matrix operation names.
	ListOperations() []string

	// IsValidOperation checks if the given operation name is supported.
	IsValidOperation(ctx context.Context, operation string) error

	// RunOperation executes the specified operation on the given matrix.
	// Returns the result as a formatted string or an error if the operation fails.
	RunOperation(ctx context.Context, matrix *entity.Matrix, operation string) (string, error)
}

type matrixOperationsDomain struct{}

// NewMatrixOperationsDomain creates a new instance of MatrixOperationsDomainInterface.
// It returns an operations service that can execute all supported matrix operations.
func NewMatrixOperationsDomain() MatrixOperationsDomainInterface {
	return &matrixOperationsDomain{}
}

func (d *matrixOperationsDomain) ListOperations() []string {
	operations := make([]string, 0, len(matrixOperations))
	for op := range matrixOperations {
		operations = append(operations, string(op))
	}
	return operations
}

func (d *matrixOperationsDomain) IsValidOperation(ctx context.Context, operation string) error {
	// Check if context is already cancelled
	if err := ctx.Err(); err != nil {
		return err
	}

	if !matrixOperations[Operation(operation)] {
		return fmt.Errorf("%w: invalid operation: %s", apperrors.ErrInvalidInput, operation)
	}
	return nil
}

func (d *matrixOperationsDomain) RunOperation(ctx context.Context, matrix *entity.Matrix, operation string) (string, error) {
	// Check if context is already cancelled
	if err := ctx.Err(); err != nil {
		return "", err
	}

	chosenOperation := Operation(operation)

	switch chosenOperation {
	case SumOperation:
		return d.sum(matrix)
	case MultiplyOperation:
		return d.multiply(matrix)
	case EchoOperation:
		return d.echo(matrix)
	case InvertOperation:
		return d.invert(matrix)
	case FlattenOperation:
		return d.flatten(matrix)
	default:
		return "", fmt.Errorf("%w: unsupported operation: %s", apperrors.ErrInvalidInput, operation)
	}
}

func (d *matrixOperationsDomain) sum(matrix *entity.Matrix) (string, error) {
	if matrix == nil || len(matrix.Data) == 0 {
		return "", fmt.Errorf("%w: empty matrix", apperrors.ErrInvalidInput)
	}

	// Use big.Int for arbitrary precision to avoid overflow
	sum := big.NewInt(0)
	for _, row := range matrix.Data {
		for _, val := range row {
			sum.Add(sum, big.NewInt(val))
		}
	}

	return sum.String(), nil
}

func (d *matrixOperationsDomain) multiply(matrix *entity.Matrix) (string, error) {
	if matrix == nil || len(matrix.Data) == 0 {
		return "", fmt.Errorf("%w: empty matrix", apperrors.ErrInvalidInput)
	}

	// Use big.Int for arbitrary precision to avoid overflow
	product := big.NewInt(1)
	for _, row := range matrix.Data {
		for _, val := range row {
			product.Mul(product, big.NewInt(val))
		}
	}

	return product.String(), nil
}

func (d *matrixOperationsDomain) echo(matrix *entity.Matrix) (string, error) {
	if matrix == nil || len(matrix.Data) == 0 {
		return "", fmt.Errorf("%w: empty matrix", apperrors.ErrInvalidInput)
	}

	var builder strings.Builder
	for i, row := range matrix.Data {
		for j, val := range row {
			if j > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(fmt.Sprintf("%d", val))
		}
		if i < len(matrix.Data)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String(), nil
}

func (d *matrixOperationsDomain) invert(matrix *entity.Matrix) (string, error) {
	if matrix == nil || len(matrix.Data) == 0 {
		return "", fmt.Errorf("%w: empty matrix", apperrors.ErrInvalidInput)
	}

	rows := len(matrix.Data)
	cols := len(matrix.Data[0])

	// Transpose the matrix
	inverted := make([][]int64, cols)
	for i := range inverted {
		inverted[i] = make([]int64, rows)
		for j := range inverted[i] {
			inverted[i][j] = matrix.Data[j][i]
		}
	}

	var builder strings.Builder
	for i, row := range inverted {
		for j, val := range row {
			if j > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(fmt.Sprintf("%d", val))
		}
		if i < len(inverted)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String(), nil
}

func (d *matrixOperationsDomain) flatten(matrix *entity.Matrix) (string, error) {
	if matrix == nil || len(matrix.Data) == 0 {
		return "", fmt.Errorf("%w: empty matrix", apperrors.ErrInvalidInput)
	}

	var builder strings.Builder
	first := true
	for _, row := range matrix.Data {
		for _, val := range row {
			if !first {
				builder.WriteString(",")
			}
			builder.WriteString(fmt.Sprintf("%d", val))
			first = false
		}
	}

	return builder.String(), nil
}
