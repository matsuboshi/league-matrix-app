package domain

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/matsuboshi/league-matrix-app/internal/repository"
	apperrors "github.com/matsuboshi/league-matrix-app/pkg/errors"
)

// MatrixDomainInterface defines the main business logic contract for matrix processing.
// It coordinates between repository, validation, and operation layers to process matrix requests.
type MatrixDomainInterface interface {
	// ListMatrixOperations returns a formatted string listing all available matrix operations.
	// It includes a sample URL and all supported operation names.
	ListMatrixOperations() (string, error)

	// ProcessMatrix executes a specific matrix operation on a file.
	// It validates the operation, reads the file, validates the matrix data, and performs the operation.
	// Returns the result as a formatted string or an error if any step fails.
	ProcessMatrix(ctx context.Context, operation string, filePath string) (string, error)
}

type matrixDomain struct {
	matrixRepository repository.MatrixRepositoryInterface
	validatorDomain  MatrixValidatorDomainInterface
	operationsDomain MatrixOperationsDomainInterface
}

// NewMatrixDomain creates a new instance of MatrixDomainInterface with all required dependencies.
// It initializes the domain service with repository, validator, and operations components.
func NewMatrixDomain() MatrixDomainInterface {
	return &matrixDomain{
		matrixRepository: repository.NewMatrixRepository(),
		validatorDomain:  NewMatrixValidatorDomain(),
		operationsDomain: NewMatrixOperationsDomain(),
	}
}

func (d *matrixDomain) ListMatrixOperations() (string, error) {
	allOperations := d.operationsDomain.ListOperations()

	operationsStr := `
	Are you lost?
	Try using this sample URL: 
	http://localhost:8080/matrix/sum?file=testdata/matrix1.csv

	Other available operations: 
	`
	for i, op := range allOperations {
		if i > 0 {
			operationsStr += ","
		}
		operationsStr += op
	}

	return operationsStr, nil
}

func (d *matrixDomain) ProcessMatrix(ctx context.Context, operation string, filePath string) (string, error) {
	// Check if context is already cancelled
	if err := ctx.Err(); err != nil {
		return "", err
	}

	if operation == "" {
		return "", fmt.Errorf("%w: operation parameter is required", apperrors.ErrInvalidInput)
	}

	err := d.validatorDomain.ValidateFilePath(ctx, filePath)
	if err != nil {
		return "", err
	}

	err = d.operationsDomain.IsValidOperation(ctx, operation)
	if err != nil {
		return "", err
	}

	rawData, err := d.matrixRepository.GetFileContent(ctx, filePath)
	if err != nil {
		return "", err
	}

	validatedMatrix, err := d.validatorDomain.Validate(ctx, rawData)
	if err != nil {
		return "", err
	}

	result, err := d.operationsDomain.RunOperation(ctx, validatedMatrix, operation)
	if err != nil {
		slog.Error("operation execution failed",
			"operation", operation,
			"error", err)
		return "", err
	}

	return result, nil
}
