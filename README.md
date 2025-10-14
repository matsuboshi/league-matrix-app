# League Matrix App

A production-grade Go HTTP service for performing matrix operations on CSV files. Built with Clean Architecture, structured logging, and comprehensive error handling.

---
## 🎯 Features

The service provides 5 matrix operations on CSV files:

### 1. **Echo**
Returns the matrix in its original format.

**Input:**
```csv
1,2,3
4,5,6
7,8,9
```

**Output:**
```
1,2,3
4,5,6
7,8,9
```

### 2. **Invert (Transpose)**
Returns the matrix with rows and columns swapped.

**Output:**
```
1,4,7
2,5,8
3,6,9
```

### 3. **Flatten**
Returns all matrix values as a single comma-separated line.

**Output:**
```
1,2,3,4,5,6,7,8,9
```

### 4. **Sum**
Returns the sum of all integers in the matrix.

**Output:**
```
45
```

### 5. **Multiply**
Returns the product of all integers in the matrix (supports arbitrarily large numbers).

**Output:**
```
362880
```

---
## 📋 Requirements

- **Go 1.25** or higher
- Make (optional, for convenience commands)


---
## 🚀 Quick Start

### 1. Run the Server
```bash
make run
```

The server will start on `http://localhost:8080`

### 2. Test All Endpoints
```bash
sh test_all_endpoints.sh
```

This script tests all operations with various test cases including edge cases and error scenarios.

---
## 🔧 Usage

### API Endpoints

**Health Check:**
```bash
curl http://localhost:8080/health
```

**List Available Operations:**
```bash
curl http://localhost:8080/
```

**Perform Matrix Operations:**
```bash
# Sum operation
curl "http://localhost:8080/matrix/sum?file=testdata/matrix1.csv"

# Echo operation
curl "http://localhost:8080/matrix/echo?file=testdata/matrix1.csv"

# Invert operation
curl "http://localhost:8080/matrix/invert?file=testdata/matrix1.csv"

# Flatten operation
curl "http://localhost:8080/matrix/flatten?file=testdata/matrix1.csv"

# Multiply operation
curl "http://localhost:8080/matrix/multiply?file=testdata/matrix1.csv"
```

### URL Format

```
http://localhost:8080/matrix/{operation}?file={filepath}
```

- `{operation}`: sum, multiply, echo, invert, or flatten
- `{filepath}`: Path to CSV file (must be in `testdata/` directory)


---
## 📁 Project Structure

```
league-matrix-app/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── entity/                 # Domain entities
│   ├── handler/                # HTTP handlers
│   ├── domain/                 # Business logic
│   └── repository/             # Data access layer
└── pkg/
    └── errors/                 # Custom error types
```

---
## 🏗️ Architecture

This project follows **Clean Architecture** principles:

```
HTTP Layer (Handler)
        ↓
Business Logic (Domain)
        ↓
Data Access (Repository)
        ↓
File System
```

### Key Design Decisions

- **Interface-driven design**: All layers interact through interfaces
- **Dependency injection**: Clean, testable component initialization
- **Structured logging**: Uses Go's `log/slog` for production-grade logging
- **Context propagation**: Request cancellation and timeout support
- **Error handling**: Sentinel errors with proper HTTP status code mapping
- **Security**: Path traversal protection, file size limits, input validation

---
## 🔒 Security Features

- ✅ **Path traversal protection**: Blocks `../` in file paths
- ✅ **Directory sandboxing**: Only allows access to `testdata/` directory
- ✅ **File type validation**: Only `.csv` files accepted
- ✅ **File size limits**: Maximum 1KB to prevent DoS attacks
- ✅ **Matrix dimension limits**: Maximum 10x10 matrices
- ✅ **Input validation**: Multiple validation layers
- ✅ **Overflow protection**: Uses `big.Int` for large number operations

---
## 📊 Available Make Commands

```bash
# Run the application
make run

# Download dependencies
make deps

# Run tests
make test

# Run tests with coverage report
make test-coverage

# Generate mocks (using mockery v3)
make mocks-generate

# Clean generated mocks
make mocks-clean
```

---
## 🧪 Test Data

The `testdata/` directory contains various CSV files for testing:

- `matrix0.csv` - Large values (1,000,000) for overflow testing
- `matrix1.csv` - Standard 9x3 matrix
- `matrix2.csv` - Invalid matrix (non-integer values)
- `matrix3.csv` - Irregular matrix
- `matrix4.csv` - Matrix with more than 10 rows
- `matrix5.csv` - Matrix with more than 10 columns
- `matrix6.csv` - Empty matrix
- `gopher.jpg.csv` - Large file (>1KB) for size limit testing

---
## 🐛 Error Handling

The service returns appropriate HTTP status codes:

| Status Code | Error Type | Example |
|-------------|------------|---------|
| 400 | Bad Request | Invalid operation, missing parameters |
| 404 | Not Found | File doesn't exist |
| 413 | Payload Too Large | File exceeds 1KB limit |
| 422 | Unprocessable Entity | Invalid CSV format, matrix validation errors |
| 504 | Gateway Timeout | Request timeout |

---
## 📝 API Response Examples

**Health Check Response:**
```bash
$ curl http://localhost:8080/health
OK
```

**Success Response:**
```bash
$ curl "http://localhost:8080/matrix/sum?file=testdata/matrix1.csv"
351
```

**Error Response:**
```bash
$ curl "http://localhost:8080/matrix/sum?file=../secret.csv"
invalid input: path traversal not allowed
```

---
## 🔍 Logging

The application uses structured logging with `log/slog`:

```
2025-10-14T10:00:00.000Z INFO starting HTTP server port=8080 address=http://localhost:8080 read_timeout=7s write_timeout=30s
2025-10-14T10:00:01.000Z INFO matrix operation completed operation=sum file_path=testdata/matrix1.csv
2025-10-14T10:00:02.000Z ERROR matrix operation failed operation=divide file_path=testdata/matrix1.csv error="invalid input: invalid operation: divide" status_code=400
```

---
## 🛑 Graceful Shutdown

The server implements graceful shutdown to ensure in-flight requests complete before stopping:

```bash
# Press Ctrl+C or send SIGTERM to stop the server
$ make run
INFO starting HTTP server port=8080 address=http://localhost:8080
^C
INFO shutdown signal received signal=interrupt
INFO gracefully shutting down server timeout=30s
INFO server stopped gracefully
```

**How it works:**
- Listens for `SIGINT` (Ctrl+C) and `SIGTERM` signals
- Stops accepting new connections
- Waits up to 30 seconds for in-flight requests to complete
- Logs shutdown progress
- Exits cleanly with proper resource cleanup

**Use cases:**
- ✅ Safe deployments (zero downtime)
- ✅ Kubernetes pod termination
- ✅ Docker container stops
- ✅ Manual server restarts

---
## 🧪 Testing

### Test Infrastructure

The project uses **modern Go testing tools**:

- **Testing Framework**: Go's built-in `testing` package
- **Mock Generation**: [Mockery v3](https://github.com/vektra/mockery) with testify mocks
- **Assertions**: [testify/assert](https://github.com/stretchr/testify)
- **Coverage Reports**: Go's native coverage tools

### Running Tests

```bash
# Run all tests
make test

# Generate coverage report
make test-coverage
# Opens coverage.html showing line-by-line coverage
```

### Mock Generation

Mocks are automatically generated from interfaces:

```bash
# Generate mocks
make mocks-generate

# Clean all mocks
make mocks-clean
```

Mock configuration is defined in `.mockery.yml` and generates mock files in `internal/mocks/` directory.

---
## 📚 What Was Implemented

✅ Clean Architecture with proper layer separation  
✅ Interface-driven design for testability  
✅ Comprehensive error handling with sentinel errors  
✅ Structured logging with `log/slog`  
✅ Context propagation for request cancellation  
✅ Security measures (path validation, size limits)  
✅ Performance optimizations (`strings.Builder`, `big.Int`)  
✅ Production-grade code quality  
✅ GoDoc documentation  
✅ Modern testing infrastructure (Mockery v3, testify)  
✅ Test coverage reporting  
✅ Health check endpoint for monitoring  
✅ HTTP server timeouts configured  
✅ Graceful shutdown with signal handling  


---
## 🤝 Best Practices

This project follows Go best practices:

- Idiomatic Go code
- Clear separation of concerns
- Explicit error handling
- Comprehensive input validation
- Structured logging
- Context-aware operations
