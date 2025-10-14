# Download dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Run the application
.PHONY: run
run:
	go run cmd/main.go

# Run tests
.PHONY: test
test:
	go test -v -race ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html
	echo "Coverage report: coverage.html"

# Install mockery
.PHONY: mocks-install
mocks-install:
	go install github.com/vektra/mockery/v3@v3.5.0

# Initialize mockery
.PHONY: mocks-init
mocks-init: mocks-install
	echo "Initializing mockery..."
	mockery init $(shell go list -m)
	echo '      recursive: true' >> .mockery.yml

# Generate mocks
.PHONY: mocks-generate
mocks-generate: mocks-install
	echo "Generating mocks..."
	mockery

# Clean generated mocks
.PHONY: mocks-clean
mocks-clean:
	echo "Cleaning generated mocks..."
	rm -rf internal/mocks
